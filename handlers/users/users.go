package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jasonwvh/simple_iam/handlers"
	"github.com/jasonwvh/simple_iam/storage"
)

func GetUser(env *handlers.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		user, err := storage.GetUser(env.DB, username)
		if err != nil {
			handlers.HandleError(w, http.StatusInternalServerError, err.Error())
			log.Print(err.Error())
			return
		}

		handlers.HandleSuccess(w, http.StatusOK, user)
	}
}

func CreateUsers(env *handlers.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user *storage.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = storage.CreateUser(env.DB, user)
		if err != nil {
			handlers.HandleError(w, http.StatusInternalServerError, err.Error())
			log.Print(err.Error())
			return
		}

		handlers.HandleSuccess(w, http.StatusOK, "New user created")
	}
}
