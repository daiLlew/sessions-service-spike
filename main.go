package main

import (
	"net/http"
	"os"

	"github.com/ONSdigital/log.go/log"
	"github.com/daiLlew/sessions-service-spike/api"
	"github.com/daiLlew/sessions-service-spike/redis"
	"github.com/daiLlew/sessions-service-spike/sessions"
	"github.com/gorilla/mux"
)

func main() {

	redisCli := redis.NewCli()

	factory := sessions.NewFactory()
	router := mux.NewRouter()

	router.HandleFunc("/session/hi", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("working"))
	}).Methods("GET")

	createSessionHandler := api.CreateSessionHandler(factory, redisCli)
	router.HandleFunc("/session", createSessionHandler).Methods("POST")

	getSessionHandler := api.GetSessionHandler(redisCli)
	router.HandleFunc("/session/{id}", getSessionHandler).Methods("GET")

	findSessionHandler := api.FindSessionHandler(redisCli)
	router.HandleFunc("/search", findSessionHandler).Methods("GET")

	log.Event(nil, "starting session service")
	if err := http.ListenAndServe(":6666", router); err != nil {
		log.Event(nil, "error starting http server", log.Error(err))
		os.Exit(1)
	}
}
