package main

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	dbHost     = "db"
	dbPort     = "5432"
	dbUser     = "postgres-dev"
	dbPassword = "s3cr3tp4ssw0rd"
	dbName     = "dev"
)

func awaitDb() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	var err error
	var db *sql.DB

	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal("Failed to open connection to db")
	}

	err = db.Ping()

	for err != nil {
		log.Error(err)
		time.Sleep(time.Second)
		err = db.Ping()
	}
	log.Println("db connection established")
	return db
}

func insertUser(register *registerPayload, hash []byte, db *sql.DB) (int, error) {
	var userid int
	row := db.QueryRow(`INSERT INTO users(username, email_address, password) VALUES($1, $2, $3) RETURNING uid`, register.Username, register.Email, hash)

	err := row.Scan(&userid)
	if err != nil {
		return -1, err
	}
	return userid, nil
}

func getUser(login *loginPayload, db *sql.DB) (*user, error) {
	row := db.QueryRow(`SELECT username, email_address, password from users WHERE username = $1`, login.Username)

	user := new(user)
	if err := row.Scan(&user.Username, &user.Email, &user.Password); err != nil {
		return nil, err
	}
	return user, nil
}

func getUsers(db *sql.DB) ([]user, error) {
	rows, err := db.Query(`SELECT username, email_address, password from users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user

	for rows.Next() {
		var user user
		err := rows.Scan(&user.Username, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}
