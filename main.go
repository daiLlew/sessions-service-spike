package main

import (
	"context"
	"net/http"
	"os"

	"github.com/ONSdigital/log.go/log"
	"github.com/daiLlew/sessions-service-spike/api"
	"github.com/daiLlew/sessions-service-spike/redis"
	"github.com/daiLlew/sessions-service-spike/sessions"
)

func main() {
	log.Event(context.Background(), "starting session service")

	router := api.Initialize(sessions.NewFactory(), redis.NewCli())

	if err := http.ListenAndServe(":6666", router); err != nil {
		log.Event(nil, "error starting http server", log.Error(err))
		os.Exit(1)
	}
}
