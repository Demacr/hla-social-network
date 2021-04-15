package storages

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/Demacr/otus-hl-socialnetwork/internal/config"
	"github.com/Demacr/otus-hl-socialnetwork/internal/models"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// DB struct contains sql.DB pointer of MySQL database
//
type DB struct {
	db *sql.DB
}

// NewDB creates new DB struct
//
func NewDB(cfg *config.Config) *DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?autocommit=true",
		cfg.MySQL.Login,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Database,
	))
	if err != nil {
		panic(err)
	}
	return &DB{db: db}
}

// WriteProfile writes to DB registration profile
//
func (db *DB) WriteProfile(profile *models.Profile) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(profile.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	result, err := db.db.Exec("INSERT INTO users(name, surname, age, sex, interests, city, email, password) VALUES(?, ?, ?, ?, ?, ?, ?, ?);",
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

// CheckCredentials checks valid credentials or not
func (db *DB) CheckCredentials(credentials *models.Credentials) (bool, error) {
	var hashedPassword string
	err := db.db.QueryRow("SELECT password FROM users WHERE email = ?", credentials.Email).Scan(&hashedPassword)
	if err != nil {
		return false, errors.Wrap(err, "CheckCredentials: No user found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
	if err != nil {
		return false, errors.Wrap(err, "CheckCredentials: Login or password mismatched")
	}

	return true, nil
}
