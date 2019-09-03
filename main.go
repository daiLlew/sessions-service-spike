package main

import (
	"net/http"

	"github.com/daiLlew/sessions-service-spike/api"
	"github.com/daiLlew/sessions-service-spike/sessions"
	"github.com/gorilla/mux"
)

func main() {
	repo := sessions.NewRepository()
	factory := sessions.NewFactory()

	router := mux.NewRouter()

	createSessionHandler := api.CreateSessionHandler(factory, repo)
	router.HandleFunc("/session", createSessionHandler).Methods("POST")

	getSessionHandler := api.GetSessionHandler(repo)
	router.HandleFunc("/session/{id}", getSessionHandler).Methods("GET")

	http.ListenAndServe(":8080", router)
}
