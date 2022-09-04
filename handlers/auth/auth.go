package auth

import (
	"log"
	"net/http"

	"github.com/jasonwvh/simple_iam/handlers"
	"github.com/jasonwvh/simple_iam/storage"
)

func AuthenticationHandler(env *handlers.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var token string

		users, err := storage.GetUsers(env.DB)
		if err != nil {
			handlers.HandleError(w, http.StatusInternalServerError, err.Error())
			log.Print(err.Error())
			return
		}

		handlers.HandleSuccess(w, http.StatusOK, users)
	}
}

func AuthorizationHandler(env *handlers.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var token string

		users, err := storage.GetUsers(env.DB)
		if err != nil {
			handlers.HandleError(w, http.StatusInternalServerError, err.Error())
			log.Print(err.Error())
			return
		}

		handlers.HandleSuccess(w, http.StatusOK, users)
	}
}
