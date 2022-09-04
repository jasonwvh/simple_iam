package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Env struct {
	DB *sql.DB
}

func HandleError(w http.ResponseWriter, code int, message string) {
	HandleSuccess(w, code, map[string]string{"error": message})
}

func HandleSuccess(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
