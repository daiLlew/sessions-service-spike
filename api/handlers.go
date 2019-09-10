package api

import (
	"encoding/json"
	"net/http"

	"github.com/daiLlew/sessions-service-spike/sessions"
	"github.com/gorilla/mux"
)

type SessionCache interface {
	Get(id string) (*sessions.Session, error)
	Set(*sessions.Session)
}

func CreateSessionHandler(factory *sessions.Factory, cache *sessions.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := factory.NewSession("test@test.com")
		cache.Set(sess)

		w.Write([]byte(sess.ID))
		w.WriteHeader(201)
	}
}

func GetSessionHandler(cache *sessions.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessID := mux.Vars(r)["id"]

		sess, err := cache.Get(sessID)
		if err != nil {
			writeErrorResponse(w, err.Error(), 500)
		} else if sess == nil {
			writeErrorResponse(w, "session not found", 404)
		} else {
			writeResponse(w, sess)
		}
	}
}

func writeResponse(w http.ResponseWriter, sess *sessions.Session) {
	b, err := json.Marshal(sess)
	if err != nil {
		writeErrorResponse(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func writeErrorResponse(w http.ResponseWriter, body string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(body))
}
