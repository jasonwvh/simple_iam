package storage

import (
	"database/sql"
	"time"
)

type User struct {
	UID      int       `json:"uid"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
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
