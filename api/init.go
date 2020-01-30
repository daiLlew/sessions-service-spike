package api

import (
	"github.com/daiLlew/sessions-service-spike/sessions"
	"github.com/gorilla/mux"
)

type SessionCache interface {
	GetByID(id string) (*sessions.Session, error)
	GetByEmail(email string) (*sessions.Session, error)
	Set(*sessions.Session) error
	FlushAll() error
}

func Initialize(sessionFactory *sessions.Factory, cache SessionCache) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/session", CreateSessionHandler(sessionFactory, cache)).Methods("POST")
	router.HandleFunc("/session/{id}", GetSessionHandler(cache)).Methods("GET")
	router.HandleFunc("/search", FindSessionHandler(cache)).Methods("GET")
	router.HandleFunc("/sessions", FlushCashHandler(cache)).Methods("DELETE")

	return router
}
