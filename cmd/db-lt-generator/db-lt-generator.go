package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	"github.com/Demacr/otus-hl-socialnetwork/internal/domain"

	"github.com/bxcodec/faker/v3"
)

var (
	flags             = flag.NewFlagSet("db-lt-generator", flag.ExitOnError)
	WORKERS           = *flags.Int("workers", 4, "Number of workers")
	RECORD_PER_WORKER = *flags.Int("records", 250000, "Number of records per worker")
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

	fmt.Println("Start generating")

	mx := sync.Mutex{}
	waitCh := make(chan interface{})
	wg := sync.WaitGroup{}
	var count int64 = 0
	go func() {
		for i := 0; i < WORKERS; i++ {
			wg.Add(1)
			go worker(&count, &wg, db, &mx)
		}
		wg.Wait()
		close(waitCh)
	}()

	for {
		select {
		case <-time.After(time.Second * 10):
			fmt.Println("Added: ", count)
		case <-waitCh:
			break
		}
	}
}

func worker(count *int64, wg *sync.WaitGroup, db *sql.DB, mx *sync.Mutex) {
	defer wg.Done()
	for i := 0; i < RECORD_PER_WORKER; i++ {
		p := &domain.Profile{}
		mx.Lock()
		if err := faker.FakeData(p); err != nil {
			log.Println(err)
		}
		mx.Unlock()
		if err := AddRecordToDB(db, p); err != nil {
			log.Println(err)
			continue
		}
		atomic.AddInt64(count, 1)
	}
}

func AddRecordToDB(db *sql.DB, profile *domain.Profile) error {
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
		// // Check duplicate email error
		// if driverErr, ok := err.(*mysql.MySQLError); ok {
		// 	if driverErr.Number == mysqlerr.ER_DUP_ENTRY {
		// 		return errors.Wrap(driverErr, "email exists")
		// 	}
		// }
		log.Println(err)
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
