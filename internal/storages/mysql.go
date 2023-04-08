package storages

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Demacr/otus-hl-socialnetwork/internal/config"
	"github.com/Demacr/otus-hl-socialnetwork/internal/domain"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	MAX_OPEN_CONNECTIONS int           = 150
	MAX_IDLE_CONNECTIONS int           = 150
	CONNECTION_IDLE_TIME time.Duration = time.Minute * 5
)

// DB struct contains sql.DB pointer of MySQL database.
type mysqlSocialNetworkRepository struct {
	Conn   *sql.DB
	slaves []*sql.DB
	count  uint64
}

// NewDB creates new DB struct.
func NewMysqlSocialNetworkRepository(cfg *config.MySQLConfig) SocialNetworkRepository {
	DSNMaster := fmt.Sprintf("%s:%s@%s/%s?autocommit=true&interpolateParams=true&parseTime=true",
		cfg.Login,
		cfg.Password,
		cfg.Host,
		cfg.Database,
	)

	db, err := sql.Open("mysql", DSNMaster)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	// Configure pool
	db.SetMaxOpenConns(MAX_OPEN_CONNECTIONS)
	db.SetMaxIdleConns(MAX_IDLE_CONNECTIONS)
	db.SetConnMaxIdleTime(CONNECTION_IDLE_TIME)

	dbSlave := db
	slaves := []*sql.DB{}
	if cfg.SlaveHosts != "" {
		slaveHosts := strings.Split(cfg.SlaveHosts, ";")
		for _, slaveHost := range slaveHosts {
			DSNSlave := fmt.Sprintf("%s:%s@%s/%s?autocommit=true&interpolateParams=true&parseTime=true",
				cfg.Login,
				cfg.Password,
				slaveHost,
				cfg.Database,
			)

			dbSlave, err = sql.Open("mysql", DSNSlave)
			if err != nil {
				panic(err)
			}
			if err = dbSlave.Ping(); err != nil {
				panic(err)
			}

			// Configure pool
			dbSlave.SetMaxOpenConns(MAX_OPEN_CONNECTIONS)
			dbSlave.SetMaxIdleConns(MAX_IDLE_CONNECTIONS)
			dbSlave.SetConnMaxIdleTime(CONNECTION_IDLE_TIME)

			slaves = append(slaves, dbSlave)
		}
	}

	return &mysqlSocialNetworkRepository{Conn: db, slaves: slaves}
}

func (m *mysqlSocialNetworkRepository) Slave() *sql.DB {
	if len(m.slaves) != 0 {
		return m.slaves[atomic.AddUint64(&m.count, 1)%uint64(len(m.slaves))]
	}
	return m.Conn
}

// WriteProfile writes to DB registration profile.
func (m *mysqlSocialNetworkRepository) WriteProfile(profile *domain.Profile) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(profile.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	result, err := m.Conn.Exec("INSERT INTO users(name, surname, age, sex, interests, city, email, password) VALUES(?, ?, ?, ?, ?, ?, ?, ?);",
		profile.Name,
		profile.Surname,
		profile.Age,
		profile.Sex,
		profile.Interests,
		profile.City,
		profile.Email,
		string(hash),
	)
	if err != nil {
		// Check duplicate email error
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqlerr.ER_DUP_ENTRY {
				return errors.Wrap(driverErr, "email exists")
			}
		}
		log.Println(err)
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (m *mysqlSocialNetworkRepository) GetProfileByEmail(email string) (*domain.Profile, error) {
	var profile domain.Profile
	err := m.Slave().QueryRow("SELECT id, name, surname, age, sex, city, interests, email, password FROM users WHERE email = ?", email).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Surname,
		&profile.Age,
		&profile.Sex,
		&profile.City,
		&profile.Interests,
		&profile.Email,
		&profile.Password,
	)
	if err != nil {
		return nil, errors.Wrap(err, "GetProfileByEmail")
	}

	return &profile, err
}

func (m *mysqlSocialNetworkRepository) GetRelatedProfileById(id, related_id int) (*domain.RelatedProfile, error) {
	var profile domain.RelatedProfile
	err := m.Slave().QueryRow("SELECT id, name, surname, age, sex, city, interests, IFNULL((SELECT true from friendship where id1=? and id2=?), false), IFNULL((SELECT true from friendship_requests where id_from=? and id_to=?), false) FROM users WHERE id = ?", related_id, id, related_id, id, id).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Surname,
		&profile.Age,
		&profile.Sex,
		&profile.City,
		&profile.Interests,
		&profile.IsFriend,
		&profile.IsRequestSent,
	)
	if err != nil {
		return nil, errors.Wrap(err, "GetRelatedProfileById")
	}

	return &profile, nil
}

func (m *mysqlSocialNetworkRepository) GetLastProfileId() (int, error) {
	var lastId int
	err := m.Conn.QueryRow("SELECT MAX(id) FROM users").Scan(&lastId)
	if err != nil {
		return lastId, errors.Wrap(err, "MySQLRepository.GetLastProfileId.QueryRow")
	}

	return lastId, nil
}

func (m *mysqlSocialNetworkRepository) CreateFriendRequest(id, friend_id int) (bool, error) {
	// TODO: check cross-request
	result, err := m.Conn.Exec("INSERT INTO friendship_requests(id_from, id_to) VALUES(?, ?)", id, friend_id)
	if err != nil {
		// Check duplicate email error
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqlerr.ER_DUP_ENTRY {
				return false, nil
			}
		}
		wrapped_error := errors.Wrap(err, "request exists")
		log.Println(wrapped_error)
		return false, wrapped_error
	}

	if _, err := result.RowsAffected(); err != nil {
		return false, err
	}

	return true, nil
}

func (m *mysqlSocialNetworkRepository) GetRandomProfiles(exclude_id int) ([]domain.Profile, error) {
	result := make([]domain.Profile, 0, 10)

	// SELECT * from (SELECT id, name, surname, age, sex, city, interests FROM users ORDER BY rand() LIMIT 10) u left join friendship on u.id = friendship.id1 left join (select * from friendship_requests where id_from = 3) fr on u.id=fr.id_to;
	profiles, err := m.Slave().Query("SELECT id, name, surname, age, sex, city, interests FROM users WHERE id != ? ORDER BY rand() LIMIT 10", exclude_id)
	if err != nil {
		return nil, errors.Wrap(err, "error during select random profiles")
	}
	defer profiles.Close()

	for profiles.Next() {
		profile := domain.Profile{}
		if err = profiles.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Surname,
			&profile.Age,
			&profile.Sex,
			&profile.City,
			&profile.Interests,
		); err != nil {
			return nil, errors.Wrap(err, "error during scan random profile")
		}

		result = append(result, profile)
	}

	return result, nil
}

func (m *mysqlSocialNetworkRepository) GetProfilesBySearchPrefixes(first_name string, last_name string) ([]domain.Profile, error) {
	result := []domain.Profile{}

	profiles, err := m.Slave().Query("SELECT id, name, surname, age, sex, city, interests FROM users WHERE name LIKE ? AND surname LIKE ? ORDER BY id ASC",
		first_name+"%",
		last_name+"%",
	)
	if err != nil {
		return nil, errors.Wrap(err, "error during searching profiles")
	}

	defer profiles.Close()
	for profiles.Next() {
		profile := domain.Profile{}
		if err = profiles.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Surname,
			&profile.Age,
			&profile.Sex,
			&profile.City,
			&profile.Interests,
		); err != nil {
			return nil, errors.Wrap(err, "error during scan searched profiles")
		}

		result = append(result, profile)
	}

	return result, nil
}

func (m *mysqlSocialNetworkRepository) GetFriendRequests(id int) ([]domain.FriendRequest, error) {
	result := []domain.FriendRequest{}

	fr, err := m.Conn.Query("SELECT id_from FROM friendship_requests WHERE id_to = ?", id)
	if err != nil {
		return nil, errors.Wrap(err, "error during select friend requests")
	}
	defer fr.Close()

	for fr.Next() {
		req := domain.FriendRequest{}
		if err = fr.Scan(
			&req.FriendID,
		); err != nil {
			return nil, errors.Wrap(err, "error during scan friend requests")
		}

		result = append(result, req)
	}

	return result, nil
}

func (m *mysqlSocialNetworkRepository) AcceptFriendship(id, friend_id int) (bool, error) {
	tx, err := m.Conn.Begin()
	if err != nil {
		return false, errors.Wrap(err, "creating transaction in accepting friendship")
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Fatalln(err)
		}
	}()

	delete_result, err := tx.Exec("DELETE FROM friendship_requests WHERE id_from = ? AND id_to = ?", friend_id, id)
	if err != nil {
		return false, errors.Wrap(err, "deleting request in accepting friendship")
	}

	rows, err := delete_result.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "getting affected rows in accepting friendship")
	}
	if rows != 1 {
		return false, nil
	}

	_, err = tx.Exec("INSERT INTO friendship(id1, id2) VALUES(?, ?)", friend_id, id)
	if err != nil {
		return false, errors.Wrap(err, "adding friend in accepting friendship")
	}

	_, err = tx.Exec("INSERT INTO friendship(id2, id1) VALUES(?, ?)", friend_id, id)
	if err != nil {
		return false, errors.Wrap(err, "adding friend in accepting friendship")
	}

	if err = tx.Commit(); err != nil {
		return false, errors.Wrap(err, "commiting in accepting friendship")
	}

	return true, nil
}

func (m *mysqlSocialNetworkRepository) DeclineFriendship(id, friend_id int) (bool, error) {
	delete_result, err := m.Conn.Exec("DELETE FROM friendship_requests WHERE id_from = ? AND id_to = ?", friend_id, id)
	if err != nil {
		return false, errors.Wrap(err, "deleting request in declining friendship")
	}

	rows, err := delete_result.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "getting affected rows in accepting friendship")
	}
	if rows != 1 {
		return false, nil
	}

	return true, nil
}

func (m *mysqlSocialNetworkRepository) GetFriends(id int) ([]int, error) {
	rows, err := m.Conn.Query("SELECT id2 FROM friendship WHERE id1 = ?", id)
	if err != nil {
		return nil, errors.Wrap(err, "MySQLRepository.GetFriends.Query")
	}

	result := make([]int, 0, 100)
	var resInt int

	for rows.Next() {
		if err = rows.Scan(&resInt); err != nil {
			return nil, errors.Wrap(err, "MySQLRepository.GetFriends.Scan")
		}

		result = append(result, resInt)
	}

	return result, nil
}

func (m *mysqlSocialNetworkRepository) CreatePost(profile_id int, post *domain.Post) (post_id int, err error) {
	result, err := m.Conn.Exec("INSERT INTO posts(profile_id, title, text) VALUES(?, ?, ?)", profile_id, post.Title, post.Text)
	if err != nil {
		return 0, errors.Wrapf(err, "error during creating post for user %d", profile_id)
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "MySQLRepository.CreatePost.LastInsertId")
	}

	if _, err := result.RowsAffected(); err != nil {
		return 0, errors.Wrapf(err, "error during creating post: rowsaffected")
	}

	return int(insertId), nil
}

func (m *mysqlSocialNetworkRepository) UpdatePost(profile_id int, post *domain.Post) error {
	tx, err := m.Conn.Begin()
	if err != nil {
		return errors.Wrap(err, "error during updating post")
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Fatalln(err)
		}
	}()

	var old_post domain.Post
	err = tx.QueryRow("SELECT id, title, text FROM posts WHERE id = ? and profile_id = ?", post.Id, profile_id).Scan(
		&old_post.Id,
		&old_post.Title,
		&old_post.Text,
	)
	if err == sql.ErrNoRows {
		return errors.New("wrong permissions to update post")
	} else if err != nil {
		return errors.Wrap(err, "error during updating post")
	}

	result, err := tx.Exec("UPDATE posts SET title = ?, text = ? WHERE id = ? and profile_id = ?", post.Title, post.Text, post.Id, profile_id)
	if err != nil {
		return errors.Wrapf(err, "error during updating post of %did with \"%s\" title and \"%s\" text", profile_id, post.Title, post.Text)
	}

	if affected, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "error during updating post rowsaffected")
	} else if affected != 1 {
		return errors.New(fmt.Sprintf("updating post affected not 1 rows but %d", affected))
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commiting in updating post")
	}

	return nil
}

func (m *mysqlSocialNetworkRepository) DeletePost(profile_id int, post *domain.Post) error {
	tx, err := m.Conn.Begin()
	if err != nil {
		return errors.Wrap(err, "error during deleting post")
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Fatalln(err)
		}
	}()

	var existing_post domain.Post
	err = tx.QueryRow("SELECT id FROM posts WHERE id = ? and profile_id = ?", post.Id, profile_id).Scan(
		&existing_post.Id,
	)
	if err == sql.ErrNoRows {
		return errors.New("wrong permissions to delete post")
	} else if err != nil {
		return errors.Wrapf(err, "error during deleting post %d", post.Id)
	}

	result, err := tx.Exec("DELETE FROM posts WHERE id = ? and profile_id = ?", post.Id, profile_id)
	if err != nil {
		return errors.Wrapf(err, "error during deleting post of %did", post.Id)
	}

	if affected, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "error during deleting post rowsaffected")
	} else if affected != 1 {
		return errors.New(fmt.Sprintf("deleting post affected not 1 rows but %d", affected))
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commiting in deleting post")
	}

	return nil
}

func (m *mysqlSocialNetworkRepository) GetPost(post_id int) (*domain.Post, error) {
	var post domain.Post
	err := m.Conn.QueryRow("SELECT id, profile_id, title, text FROM posts WHERE id = ?", post_id).Scan(
		&post.Id,
		&post.ProfileId,
		&post.Title,
		&post.Text,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error during getting post %d id", post_id)
	}

	return &post, nil
}

func (m *mysqlSocialNetworkRepository) GetFeedLastN(profileId int, N int) (result []int, err error) {
	rows, err := m.Conn.Query("SELECT posts.id FROM friendship JOIN posts on friendship.id2 = posts.profile_id WHERE id1 = ? ORDER BY posts.id LIMIT ?", profileId, N)
	if err != nil {
		return nil, errors.Wrap(err, "MySQLRepository.GetFeedLastN.Query")
	}

	result = make([]int, 0, N)

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Println(errors.Wrap(err, "MySQLRepository.GetFeedLastN.Scan"))
		}
		result = append(result, id)
	}

	return result, nil
}

// TODO: move out getting dialog_id and possibly cache it locally
func (m *mysqlSocialNetworkRepository) CreateMessage(message *domain.Message) error {
	id1, id2 := message.From, message.To
	if id2 < id1 {
		id1, id2 = id2, id1
	}

	tx, err := m.Conn.Begin()
	if err != nil {
		return errors.Wrap(err, "MySQLRepository.CreateMessage.Begin")
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Fatalln(err)
		}
	}()

	var dialogId int
	err = tx.QueryRow("SELECT id FROM dialogs WHERE id1 = ? AND id2 = ?", id1, id2).Scan(&dialogId)
	if errors.Is(err, sql.ErrNoRows) {
		res, err := tx.Exec("INSERT INTO dialogs(id1, id2) VALUES(?, ?)", id1, id2)
		if err != nil {
			return errors.Wrap(err, "MySQLRepository.CreateMessage.Exec.INSERTINTODialogs")
		}
		dialogId64, err := res.LastInsertId()
		if err != nil {
			return errors.Wrap(err, "MySQLRepository.CreateMessage.LastInsertId")
		}
		dialogId = int(dialogId64)
	} else if err != nil {
		return errors.Wrap(err, "MySQLRepository.CreateMessage.QueryRow")
	}

	result, err := tx.Exec("INSERT INTO messages(dialog_id, id_from, id_to, seq, ts, text) VALUES(?, ?, ?, (SELECT COALESCE(MAX(seq), 0) FROM messages as m WHERE m.dialog_id = ?) + 1, ?, ?);",
		dialogId,
		message.From,
		message.To,
		dialogId,
		message.Timestamp,
		message.Text,
	)
	if err != nil {
		return errors.Wrap(err, "MySQLRepository.CreateMessage.Exec")
	}

	if rows, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "MySQLRepository.CreateMessage.RowsAffected")
	} else if rows != 1 {
		return errors.New("MySQLRepository.CreateMessage.RowsAffected: affected rows doesn't equal to 1")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "MySQLRepository.CreateMessage.Commit")
	}

	return nil
}

func (m *mysqlSocialNetworkRepository) GetDialog(id1 int, id2 int) ([]*domain.Message, error) {
	if id1 > id2 {
		id1, id2 = id2, id1
	}

	rows, err := m.Conn.Query("SELECT id_from, id_to, ts, text FROM messages WHERE dialog_id = (SELECT id from dialogs WHERE id1 = ? AND id2 = ?) ORDER BY seq", id1, id2)
	if err != nil {
		return nil, errors.Wrap(err, "MySQLRepository.GetDialog.Query")
	}

	result := make([]*domain.Message, 0)
	for rows.Next() {
		message := &domain.Message{}
		err = rows.Scan(
			&message.From,
			&message.To,
			&message.Timestamp,
			&message.Text,
		)
		if err != nil {
			return nil, errors.Wrap(err, "MySQLRepository.GetDialog.Scan")
		}

		result = append(result, message)
	}

	return result, nil
}

func (m *mysqlSocialNetworkRepository) GetDialogList(id int) ([]*domain.DialogPreview, error) {
	rows, err := m.Conn.Query(`select messages.dialog_id, messages.id_from, messages.id_to, messages.text from messages 
															join
																(select dialog_id, max(seq) as ms from messages where id_from = ? or id_to = ? group by dialog_id) as mm 
															on messages.dialog_id = mm.dialog_id and messages.seq = mm.ms 
															order by ts desc`, id, id)
	if err != nil {
		return nil, errors.Wrap(err, "MySQLRepository.GetDialogList.Query")
	}

	result := make([]*domain.DialogPreview, 0)
	tempId1, tempId2 := 0, 0
	for rows.Next() {
		dialogPreview := &domain.DialogPreview{}
		err = rows.Scan(
			&dialogPreview.DialogID,
			&tempId1,
			&tempId2,
			&dialogPreview.LastMessage,
		)
		if err != nil {
			return nil, errors.Wrap(err, "MySQLRepository.GetDialogList.Scan")
		}

		if tempId1 != id {
			dialogPreview.FriendID = tempId1
		} else {
			dialogPreview.FriendID = tempId2
		}

		result = append(result, dialogPreview)
	}

	return result, nil
}
