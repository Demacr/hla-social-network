package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Demacr/otus-hl-socialnetwork/internal/domain"
	"github.com/VividCortex/mysqlerr"
	"github.com/bxcodec/faker/v3"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	flags          = flag.NewFlagSet("db-lt-generator", flag.ExitOnError)
	Workers        = *flags.Int("workers", 4, "Number of workers")
	UsersPerWorker = *flags.Int("records", 2500, "Number of records per worker")
)

const (
	Timer = 10 * time.Second
)

func main() {
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()

		return
	}

	dbstring := args[0]

	db, err := sql.Open("mysql", dbstring)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(150)
	db.SetMaxIdleConns(150)
	db.SetConnMaxLifetime(time.Minute * 5)

	fmt.Println("Start generating")

	mx := sync.Mutex{}
	waitCh := make(chan interface{})
	wg := sync.WaitGroup{}
	var count int64 = 0

	go func() {
		for i := 0; i < Workers; i++ {
			wg.Add(1)

			go worker(&count, &wg, db, &mx)
		}
		wg.Wait()
		close(waitCh)
	}()

L:
	for {
		select {
		case <-time.After(Timer):
			fmt.Println("Added: ", count)
		case <-waitCh:
			fmt.Println("Added: ", count)

			break L
		}
	}

	// Friendship workers.
	count = 0
	waitCh = make(chan interface{})
	go func() {
		for i := 0; i < Workers; i++ {
			wg.Add(1)
			i := i
			go workerFriendship(i, &count, &wg, db)
		}

		wg.Wait()
		close(waitCh)
	}()
	<-waitCh
}

func worker(count *int64, wg *sync.WaitGroup, db *sql.DB, mx *sync.Mutex) {
	defer wg.Done()

	for i := 0; i < UsersPerWorker; i++ {
		p := &domain.Profile{}

		mx.Lock()

		if err := faker.FakeData(p); err != nil {
			log.Println(err)
		}

		mx.Unlock()

		profileId, err := AddRecordToDB(db, p)
		if err != nil {
			log.Println(err)

			continue
		}

		postCount := int(rand.NormFloat64()*3 + 10)
		if postCount > 0 {
			posts := make([]domain.Post, postCount)
			for _, post := range posts {
				mx.Lock()
				if err := faker.FakeData(&post); err != nil {
					log.Println(err)
				}
				post.ProfileId = int(profileId)
				mx.Unlock()

				AddPostToDB(db, &post)
			}
		}
		atomic.AddInt64(count, 1)
	}
}

func workerFriendship(index int, count *int64, wg *sync.WaitGroup, db *sql.DB) {
	defer wg.Done()

	for i := 0; i < UsersPerWorker; i++ {
		friendshipsCount := int(rand.NormFloat64()*5 + 150)
		for j := 0; j < friendshipsCount; j++ {
			AddFriendshipToDB(db, index*UsersPerWorker+i, rand.Intn(Workers*UsersPerWorker)+1)
		}

	}
}

func AddRecordToDB(db *sql.DB, profile *domain.Profile) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(profile.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	result, err := db.Exec("INSERT INTO users(name, surname, age, sex, interests, city, email, password) VALUES(?, ?, ?, ?, ?, ?, ?, ?);",
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
				return 0, errors.Wrap(driverErr, "email exists")
			}
		}

		return 0, errors.Wrap(err, "error during adding new record")
	}

	profileId, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "error during getting affected rows")
	}

	return profileId, nil
}

func AddPostToDB(db *sql.DB, post *domain.Post) error {
	result, err := db.Exec("INSERT INTO posts(profile_id, title, text) values(?, ?, ?)", post.ProfileId, post.Title, post.Text)
	if err != nil {
		return errors.Wrap(err, "AddPostToDB.Exec")
	}

	_, err = result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "error during getting affected rows")
	}

	return nil
}

func AddFriendshipToDB(db *sql.DB, id1 int, id2 int) error {
	_, err := db.Exec("INSERT INTO friendship VALUES (?, ?)", id1, id2)
	if err != nil {
		return errors.Wrap(err, "AddFriendshipToDB.Exec")
	}

	_, err = db.Exec("INSERT INTO friendship VALUES (?, ?)", id2, id1)
	if err != nil {
		return errors.Wrap(err, "AddFriendshipToDB.Exec")
	}

	return nil
}
