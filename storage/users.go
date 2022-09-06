package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type User struct {
	UID      int       `json:"uid"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
}

func GetUser(db *sql.DB, username string) (*User, error) {
	queryString := fmt.Sprintf("SELECT * FROM users WHERE username='%s'", username)
	rows, err := db.Query(queryString)
	if err != nil {
		return nil, err
	}

	var users []User

	for rows.Next() {
		var user User
		err = rows.Scan(&user.UID, &user.Username, &user.Password, &user.Created)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	rows.Close()

	return &users[0], err
}

func GetUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	var users []User

	for rows.Next() {
		var user User
		err = rows.Scan(&user.UID, &user.Username, &user.Password, &user.Created)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	rows.Close()

	return users, err
}

func CreateUser(db *sql.DB, user *User) error {
	now := time.Now()
	queryString := fmt.Sprintf("INSERT INTO users(username, password, created) values('%s', '%s', '%s')", user.Username, user.Password, now)
	res, err := db.Exec(queryString)
	if err != nil {
		return err
	}

	log.Printf("CreateUser result: %v", res)

	return err
}
