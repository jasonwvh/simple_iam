package storage

import (
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

type User struct {
	UID      int       `json:"uid"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Salt     string    `json:"salt"`
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
		err = rows.Scan(&user.UID, &user.Username, &user.Password, &user.Salt, &user.Created)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	rows.Close()

	return &users[0], nil
}

func GetUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	var users []User

	for rows.Next() {
		var user User
		err = rows.Scan(&user.UID, &user.Username, &user.Password, &user.Salt, &user.Created)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	rows.Close()

	return users, nil
}

func CreateUser(db *sql.DB, user *User) error {
	salt := generateSalt()
	saltPass := user.Password + salt

	hasher := sha512.New512_256()
	hasher.Write([]byte(saltPass))
	hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	now := time.Now()
	queryString := fmt.Sprintf("INSERT INTO users(username, password, salt, created) values('%s', '%s', '%s', '%s')", user.Username, hashedPass, salt, now)
	res, err := db.Exec(queryString)
	if err != nil {
		return err
	}
	ids, _ := res.LastInsertId()
	rws, _ := res.RowsAffected()

	log.Printf("CreateUser result: ids: %d rows: %d", ids, rws)

	return nil
}

func generateSalt() string {
	var alphabet []rune = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	alphabetSize := len(alphabet)
	var sb strings.Builder

	for i := 0; i < 10; i++ {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}

	s := sb.String()
	return s
}
