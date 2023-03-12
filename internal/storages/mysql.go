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
	DSNMaster := fmt.Sprintf("%s:%s@%s/%s?autocommit=true&interpolateParams=true",
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
			DSNSlave := fmt.Sprintf("%s:%s@%s/%s?autocommit=true&interpolateParams=true",
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
	err := m.Slave().QueryRow("SELECT id, name, surname, age, sex, city, interests FROM users WHERE email = ?", email).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Surname,
		&profile.Age,
		&profile.Sex,
		&profile.City,
		&profile.Interests,
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

// CheckCredentials checks valid credentials or not
func (m *mysqlSocialNetworkRepository) CheckCredentials(credentials *domain.Credentials) (bool, error) {
	var hashedPassword string
	err := m.Conn.QueryRow("SELECT password FROM users WHERE email = ?", credentials.Email).Scan(&hashedPassword)
	if err != nil {
		return false, errors.Wrap(err, "CheckCredentials: No user found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
	if err != nil {
		return false, errors.Wrap(err, "CheckCredentials: Login or password mismatched")
	}

	return true, nil
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
