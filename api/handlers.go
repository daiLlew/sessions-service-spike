package api

import (
	"encoding/json"
	"net/http"

	"github.com/daiLlew/sessions-service-spike/sessions"
	"github.com/gorilla/mux"
)

func CreateSessionHandler(factory *sessions.Factory, repository *sessions.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := factory.NewSession("test@test.com")
		repository.Save(sess)

		w.Write([]byte(sess.ID))
		w.WriteHeader(201)
	}
}

func GetSessionHandler(repository *sessions.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessID := mux.Vars(r)["id"]
		sess, _ := repository.GetByID(sessID)

		b, _ := json.Marshal(sess)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}
