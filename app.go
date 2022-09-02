package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jasonwvh/simple_iam/storage"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	Router  *mux.Router
	Storage storage.Storage
}

func (a *App) Initialize() {
	log.Println("initialization")
	var err error
	a.Storage.DB, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = a.Storage.DB.Exec("CREATE TABLE IF NOT EXISTS `users` (`uid` INTEGER PRIMARY KEY AUTOINCREMENT, `username` VARCHAR(255) NOT NULL, `password` VARCHAR(255) NOT NULL, `created` DATE NULL)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = a.Storage.DB.Exec("INSERT INTO users(username, password, created) values('jasonwvh', 'pass', '2012-12-09')")
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/authenticate", a.AuthenticationHandler).Methods("POST")
	a.Router.HandleFunc("/authorize", a.AuthorizationHandler).Methods("POST")
}

func (a *App) AuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	var token string

	log.Print("auth handler")

	users, err := a.Storage.GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("users: %v", users)

	json.NewEncoder(w).Encode(token)
}

func (a *App) AuthorizationHandler(w http.ResponseWriter, r *http.Request) {
	var token string

	json.NewEncoder(w).Encode(token)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
