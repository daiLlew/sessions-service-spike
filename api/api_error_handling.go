package api

import (
	"context"
	"net/http"

	"github.com/ONSdigital/log.go/log"
	"github.com/daiLlew/sessions-service-spike/sessions"
)

func handleCreateSessionError(ctx context.Context, w http.ResponseWriter, err error) {
	var status int
	var body string

	switch err {
	case BadRequestErr:
		status = http.StatusBadRequest
		body = err.Error()
	default:
		status = http.StatusInternalServerError
		body = "internal server error"
	}

	log.Event(ctx, "create session request unsuccessful", log.Error(err), log.Data{"status": status})
	writeErrorResponse(ctx, w, body, status)
}

func handleGetSessionError(ctx context.Context, w http.ResponseWriter, err error) {
	switch err {
	case sessions.SessionNotFoundErr:
		http.Error(w, "session not found", http.StatusNotFound)
	case sessions.SessionExpiredErr:
		http.Error(w, "session not found", http.StatusUnauthorized)
	default:
		log.Event(ctx, "internal server error ", log.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func handleFindSessionError(ctx context.Context, w http.ResponseWriter, err error) {
	switch err {
	case BadRequestErr:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case sessions.SessionNotFoundErr:
		http.Error(w, "session not found", http.StatusNotFound)
	case sessions.SessionExpiredErr:
		http.Error(w, "session not found", http.StatusUnauthorized)
	default:
		log.Event(ctx, "internal server error ", log.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func handleFlushAllError(ctx context.Context, w http.ResponseWriter, err error) {
	log.Event(ctx, "internal server error ", log.Error(err))
	http.Error(w, "internal server error", http.StatusInternalServerError)
}
