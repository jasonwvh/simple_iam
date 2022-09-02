package storage

import (
	"database/sql"
	"log"
	"time"
)

type Storage struct {
	DB *sql.DB
}

type User struct {
	UID      int       `json:"uid"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
}

func (storage *Storage) GetUsers() ([]User, error) {
	log.Printf("GetUsers()")
	log.Printf("DB: %v", storage)
	rows, err := storage.DB.Query("SELECT * FROM users")
	if err != nil {
		log.Fatalf("error query %s", err)
	}

	var users []User
	var uid int
	var username string
	var password string
	var created time.Time

	for rows.Next() {
		err = rows.Scan(&uid, &username, &password, &created)
		if err != nil {
			log.Fatalf("error scan %s", err)
		}

		user := User{
			UID:      uid,
			Username: username,
			Password: password,
			Created:  created,
		}
		users = append(users, user)
	}

	rows.Close()

	return users, err
}
