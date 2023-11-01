package storages

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync/atomic"

	"github.com/Demacr/otus-hl-socialnetwork/internal/config"
	"github.com/Demacr/otus-hl-socialnetwork/internal/domain"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// DB struct contains sql.DB pointer of Postgres database.
type postgresSocialNetworkRepository struct {
	Conn   *sql.DB
	slaves []*sql.DB
	count  uint64
}

// NewDB creates new DB struct.
func NewPostgresSocialNetworkRepository(cfg *config.PostgreSQLConfig) SocialNetworkRepository {
	DSNMaster := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		5432,
		cfg.Login,
		cfg.Password,
		cfg.Database,
	)

	db, err := sql.Open("postgres", DSNMaster)
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

			dbSlave, err := sql.Open("postgres", DSNSlave)
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

	return &postgresSocialNetworkRepository{Conn: db, slaves: slaves}
}

func (m *postgresSocialNetworkRepository) Slave() *sql.DB {
	if len(m.slaves) != 0 {
		return m.slaves[atomic.AddUint64(&m.count, 1)%uint64(len(m.slaves))]
	}
	return m.Conn
}

// WriteProfile writes to DB registration profile.
func (m *postgresSocialNetworkRepository) WriteProfile(profile *domain.Profile) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(profile.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	result, err := m.Conn.Exec("INSERT INTO users(name, surname, age, sex, interests, city, email, password) VALUES($1, $2, $3, $4, $5, $6, $7, $8);",
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
		return err
	}
	if rowsaffected, err := result.RowsAffected(); err != nil {
		return err
	} else if rowsaffected != 1 {
		return errors.New("email exists")
	}

	return nil
}

func (m *postgresSocialNetworkRepository) GetProfileByEmail(email string) (*domain.Profile, error) {
	var profile domain.Profile
	err := m.Slave().QueryRow("SELECT id, name, surname, age, sex, city, interests, email, password FROM users WHERE email = $1", email).Scan(
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

func (m *postgresSocialNetworkRepository) GetRelatedProfileById(id, related_id int) (*domain.RelatedProfile, error) {
	var profile domain.RelatedProfile
	err := m.Slave().QueryRow("SELECT id, name, surname, age, sex, city, interests, COALESCE((SELECT true from friendship where id1=$1 and id2=$2), false), COALESCE((SELECT true from friendship_requests where id_from=$1 and id_to=$2), false) FROM users WHERE id = $2", related_id, id).Scan(
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

func (m *postgresSocialNetworkRepository) GetLastProfileId() (int, error) {
	var lastId int
	err := m.Conn.QueryRow("SELECT MAX(id) FROM users").Scan(&lastId)
	if err != nil {
		return lastId, errors.Wrap(err, "PostgresRepository.GetLastProfileId.QueryRow")
	}

	return lastId, nil
}

func (m *postgresSocialNetworkRepository) CreateFriendRequest(id, friend_id int) (bool, error) {
	// TODO: check cross-request
	result, err := m.Conn.Exec("INSERT INTO friendship_requests(id_from, id_to) VALUES($1, $2)", id, friend_id)
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

func (m *postgresSocialNetworkRepository) GetRandomProfiles(exclude_id int) ([]domain.Profile, error) {
	result := make([]domain.Profile, 0, 10)

	// SELECT * from (SELECT id, name, surname, age, sex, city, interests FROM users ORDER BY rand() LIMIT 10) u left join friendship on u.id = friendship.id1 left join (select * from friendship_requests where id_from = 3) fr on u.id=fr.id_to;
	profiles, err := m.Slave().Query("SELECT id, name, surname, age, sex, city, interests FROM users WHERE id != $1 ORDER BY RANDOM() LIMIT 10", exclude_id)
	if err != nil || profiles.Err() != nil {
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

func (m *postgresSocialNetworkRepository) GetProfilesBySearchPrefixes(first_name string, last_name string) ([]domain.Profile, error) {
	result := []domain.Profile{}

	profiles, err := m.Slave().Query("SELECT id, name, surname, age, sex, city, interests FROM users WHERE name LIKE $1 AND surname LIKE $2 ORDER BY id ASC",
		first_name+"%",
		last_name+"%",
	)
	if err != nil || profiles.Err() != nil {
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

func (m *postgresSocialNetworkRepository) GetFriendRequests(id int) ([]domain.FriendRequest, error) {
	result := []domain.FriendRequest{}

	fr, err := m.Conn.Query("SELECT id_from FROM friendship_requests WHERE id_to = $1", id)
	if err != nil || fr.Err() != nil {
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

func (m *postgresSocialNetworkRepository) AcceptFriendship(id, friend_id int) (bool, error) {
	tx, err := m.Conn.Begin()
	if err != nil {
		return false, errors.Wrap(err, "creating transaction in accepting friendship")
	}
	defer func() {
		if err != nil {
			err = errors.Wrap(tx.Rollback(), "error during rollback")
		}
	}()

	delete_result, err := tx.Exec("DELETE FROM friendship_requests WHERE id_from = $1 AND id_to = $2", friend_id, id)
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

	_, err = tx.Exec("INSERT INTO friendship(id1, id2) VALUES($1, $2), ($2, $1)", friend_id, id)
	if err != nil {
		return false, errors.Wrap(err, "adding friend in accepting friendship")
	}

	if err = tx.Commit(); err != nil {
		return false, errors.Wrap(err, "committing in accepting friendship")
	}

	return true, nil
}

func (m *postgresSocialNetworkRepository) DeclineFriendship(id, friend_id int) (bool, error) {
	delete_result, err := m.Conn.Exec("DELETE FROM friendship_requests WHERE id_from = $1 AND id_to = $2", friend_id, id)
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

func (m *postgresSocialNetworkRepository) GetFriends(id int) ([]int, error) {
	rows, err := m.Conn.Query("SELECT id2 FROM friendship WHERE id1 = $1", id)
	if err != nil || rows.Err() != nil {
		return nil, errors.Wrap(err, "PostgresRepository.GetFriends.Query")
	}
	defer rows.Close()

	result := make([]int, 0, 100)
	var resInt int

	for rows.Next() {
		if err = rows.Scan(&resInt); err != nil {
			return nil, errors.Wrap(err, "PostgresRepository.GetFriends.Scan")
		}

		result = append(result, resInt)
	}

	return result, nil
}

func (m *postgresSocialNetworkRepository) CreatePost(profile_id int, post *domain.Post) (post_id int, err error) {
	err = m.Conn.QueryRow("INSERT INTO posts(profile_id, title, text) VALUES($1, $2, $3) RETURNING id", profile_id, post.Title, post.Text).Scan(&post_id)
	if err != nil {
		return 0, errors.Wrapf(err, "error during creating post for user %d", profile_id)
	}

	return post_id, nil
}

func (m *postgresSocialNetworkRepository) UpdatePost(profile_id int, post *domain.Post) error {
	tx, err := m.Conn.Begin()
	if err != nil {
		return errors.Wrap(err, "error during updating post")
	}
	defer func() {
		if err != nil {
			err = errors.Wrap(tx.Rollback(), "error during rollback")
		}
	}()

	var old_post domain.Post
	err = tx.QueryRow("SELECT id, title, text FROM posts WHERE id = $1 and profile_id = $2", post.ID, profile_id).Scan(
		&old_post.ID,
		&old_post.Title,
		&old_post.Text,
	)
	if err == sql.ErrNoRows {
		return errors.New("wrong permissions to update post")
	} else if err != nil {
		return errors.Wrap(err, "error during updating post")
	}

	result, err := tx.Exec("UPDATE posts SET title = $1, text = $2 WHERE id = $3 and profile_id = $4", post.Title, post.Text, post.ID, profile_id)
	if err != nil {
		return errors.Wrapf(err, "error during updating post of %did with \"%s\" title and \"%s\" text", profile_id, post.Title, post.Text)
	}

	if affected, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "error during updating post rowsaffected")
	} else if affected != 1 {
		return errors.New(fmt.Sprintf("updating post affected not 1 rows but %d", affected))
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "committing in updating post")
	}

	return nil
}

func (m *postgresSocialNetworkRepository) DeletePost(profile_id int, post *domain.Post) error {
	tx, err := m.Conn.Begin()
	if err != nil {
		return errors.Wrap(err, "error during deleting post")
	}
	defer func() {
		if err != nil {
			err = errors.Wrap(tx.Rollback(), "error during rollback")
		}
	}()

	var existing_post domain.Post
	err = tx.QueryRow("SELECT id FROM posts WHERE id = $1 and profile_id = $2", post.ID, profile_id).Scan(
		&existing_post.ID,
	)
	if err == sql.ErrNoRows {
		return errors.New("wrong permissions to delete post")
	} else if err != nil {
		return errors.Wrapf(err, "error during deleting post %d", post.ID)
	}

	result, err := tx.Exec("DELETE FROM posts WHERE id = $1 and profile_id = $2", post.ID, profile_id)
	if err != nil {
		return errors.Wrapf(err, "error during deleting post of %did", post.ID)
	}

	if affected, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "error during deleting post rowsaffected")
	} else if affected != 1 {
		return errors.New(fmt.Sprintf("deleting post affected not 1 rows but %d", affected))
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "committing in deleting post")
	}

	return nil
}

func (m *postgresSocialNetworkRepository) GetPost(post_id int) (*domain.Post, error) {
	var post domain.Post
	err := m.Conn.QueryRow("SELECT id, profile_id, title, text FROM posts WHERE id = $1", post_id).Scan(
		&post.ID,
		&post.ProfileID,
		&post.Title,
		&post.Text,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error during getting post %d id", post_id)
	}

	return &post, nil
}

func (m *postgresSocialNetworkRepository) GetFeedLastN(profileId int, n int) (result []int, err error) {
	rows, err := m.Conn.Query("SELECT posts.id FROM friendship JOIN posts on friendship.id2 = posts.profile_id WHERE id1 = $1 ORDER BY posts.id LIMIT $2", profileId, n)
	if err != nil || rows.Err() != nil {
		return nil, errors.Wrap(err, "PostgresRepository.GetFeedLastN.Query")
	}
	defer rows.Close()

	result = make([]int, 0, n)

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Println(errors.Wrap(err, "PostgresRepository.GetFeedLastN.Scan"))
		}
		result = append(result, id)
	}

	return result, nil
}

// TODO: move out getting dialog_id and possibly cache it locally.
func (m *postgresSocialNetworkRepository) CreateMessage(message *domain.Message) error {
	id1, id2 := message.From, message.To
	if id2 < id1 {
		id1, id2 = id2, id1
	}

	tx, err := m.Conn.Begin()
	if err != nil {
		return errors.Wrap(err, "PostgresRepository.CreateMessage.Begin")
	}
	defer func() {
		if err != nil {
			err = errors.Wrap(tx.Rollback(), "error during rollback")
		}
	}()

	var dialogID int
	err = tx.QueryRow("SELECT id FROM dialogs WHERE id1 = $1 AND id2 = $2", id1, id2).Scan(&dialogID)
	if errors.Is(err, sql.ErrNoRows) {
		if err := tx.QueryRow("INSERT INTO dialogs(id1, id2) VALUES($1, $2) RETURNING id", id1, id2).Scan(&dialogID); err != nil {
			return errors.Wrap(err, "PostgresRepository.CreateMessage.QueryRow.INSERTINTODialogs")
		}
	} else if err != nil {
		return errors.Wrap(err, "PostgresRepository.CreateMessage.QueryRow")
	}

	result, err := tx.Exec("INSERT INTO messages(dialog_id, id_from, id_to, seq, ts, text) VALUES($1, $2, $3, (SELECT COALESCE(MAX(seq), 0) FROM messages as m WHERE m.dialog_id = $1) + 1, $4, $5);",
		dialogID,
		message.From,
		message.To,
		message.Timestamp,
		message.Text,
	)
	if err != nil {
		return errors.Wrap(err, "PostgresRepository.CreateMessage.Exec")
	}

	if rows, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "PostgresRepository.CreateMessage.RowsAffected")
	} else if rows != 1 {
		return errors.New("PostgresRepository.CreateMessage.RowsAffected: affected rows doesn't equal to 1")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "PostgresRepository.CreateMessage.Commit")
	}

	return nil
}

func (m *postgresSocialNetworkRepository) GetDialog(id1 int, id2 int) ([]*domain.Message, error) {
	if id1 > id2 {
		id1, id2 = id2, id1
	}

	rows, err := m.Conn.Query("SELECT id_from, id_to, ts, text FROM messages WHERE dialog_id = (SELECT id from dialogs WHERE id1 = $1 AND id2 = $2) ORDER BY seq", id1, id2)
	if err != nil || rows.Err() != nil {
		return nil, errors.Wrap(err, "PostgresRepository.GetDialog.Query")
	}
	defer rows.Close()

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
			return nil, errors.Wrap(err, "PostgresRepository.GetDialog.Scan")
		}

		result = append(result, message)
	}

	return result, nil
}

func (m *postgresSocialNetworkRepository) GetDialogList(id int) ([]*domain.DialogPreview, error) {
	rows, err := m.Conn.Query(`select messages.dialog_id, messages.id_from, messages.id_to, messages.text from messages 
															join
																(select dialog_id, max(seq) as ms from messages where id_from = $1 or id_to = $1 group by dialog_id) as mm 
															on messages.dialog_id = mm.dialog_id and messages.seq = mm.ms 
															order by ts desc`, id)
	if err != nil || rows.Err() != nil {
		return nil, errors.Wrap(err, "PostgresRepository.GetDialogList.Query")
	}
	defer rows.Close()

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
			return nil, errors.Wrap(err, "PostgresRepository.GetDialogList.Scan")
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
