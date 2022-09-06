package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jasonwvh/simple_iam/handlers"
	"github.com/jasonwvh/simple_iam/handlers/auth"
	"github.com/jasonwvh/simple_iam/handlers/users"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `users` (`uid` INTEGER PRIMARY KEY AUTOINCREMENT, `username` VARCHAR(255) NOT NULL, `password` VARCHAR(255) NOT NULL, `created` DATE NULL)")
	if err != nil {
		log.Fatal(err)
	}

	// _, err = db.Exec("INSERT INTO users(username, password, created) values('jasonwvh', 'pass', '2012-12-09')")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	env := &handlers.Env{DB: db}

	r := mux.NewRouter()
	usersRouter := r.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/", users.CreateUsers(env)).Methods("POST")
	usersRouter.HandleFunc("/{username}", users.GetUser(env)).Methods("GET")

	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.HandleFunc("/authenticate", auth.AuthenticationHandler(env)).Methods("POST")
	authRouter.HandleFunc("/authorize", auth.AuthorizationHandler(env)).Methods("POST")

	// create a new server
	s := http.Server{
		Addr:         "localhost:8080",  // configure the bind address
		Handler:      r,                 // set the default handler
		ErrorLog:     &log.Logger{},     // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		log.Println("Starting server on port 8080")

		err := s.ListenAndServe()
		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
