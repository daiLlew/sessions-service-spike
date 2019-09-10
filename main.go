package main

import (
	"net/http"
	"time"

	"github.com/daiLlew/sessions-service-spike/api"
	"github.com/daiLlew/sessions-service-spike/sessions"
	"github.com/gorilla/mux"
)

func main() {
	sessionCache := sessions.NewCache(time.Second * 30, time.Minute * 1)
	factory := sessions.NewFactory()

	router := mux.NewRouter()

	createSessionHandler := api.CreateSessionHandler(factory, sessionCache)
	router.HandleFunc("/session", createSessionHandler).Methods("POST")

	getSessionHandler := api.GetSessionHandler(sessionCache)
	router.HandleFunc("/session/{id}", getSessionHandler).Methods("GET")

	http.ListenAndServe(":8080", router)
}
