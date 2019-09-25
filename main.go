package main

import (
	"net/http"
	"os"
	"time"

	"github.com/ONSdigital/log.go/log"
	"github.com/daiLlew/sessions-service-spike/api"
	"github.com/daiLlew/sessions-service-spike/sessions"
	"github.com/gorilla/mux"
)

func main() {
	purgeInterval := time.Second * 15
	sessionTTL := time.Second * 30

	sessionCache := sessions.NewCache(purgeInterval, sessionTTL)
	factory := sessions.NewFactory()
	router := mux.NewRouter()

	router.HandleFunc("/session/hi", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("working"))
	}).Methods("GET")

	createSessionHandler := api.CreateSessionHandler(factory, sessionCache)
	router.HandleFunc("/session", createSessionHandler).Methods("POST")

	getSessionHandler := api.GetSessionHandler(sessionCache)
	router.HandleFunc("/session/{id}", getSessionHandler).Methods("GET")

	findSessionHandler := api.FindSessionHandler(sessionCache)
	router.HandleFunc("/search", findSessionHandler).Methods("GET")

	if err := http.ListenAndServe(":6666", router); err != nil {
		log.Event(nil, "error starting http server", log.Error(err))
		os.Exit(1)
	}
}
