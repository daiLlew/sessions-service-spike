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
	var status int
	var body string

	switch err {
	case sessions.SessionNotFoundErr:
		status = http.StatusNotFound
		body = err.Error()
	case sessions.SessionExpiredErr:
		status = http.StatusUnauthorized
		body = err.Error()
	default:
		status = http.StatusInternalServerError
		body = "internal server error"
	}

	log.Event(ctx, "get session request unsuccessful", log.Error(err), log.Data{"status": status})
	writeErrorResponse(ctx, w, body, status)
}

func handleFindSessionError(ctx context.Context, w http.ResponseWriter, err error) {
	var status int
	var body string

	switch err {
	case BadRequestErr:
		status = http.StatusBadRequest
		body = err.Error()
	case sessions.SessionNotFoundErr:
		status = http.StatusNotFound
		body = err.Error()
	case sessions.SessionExpiredErr:
		status = http.StatusUnauthorized
		body = err.Error()
	default:
		status = http.StatusInternalServerError
		body = "internal server error"
	}

	log.Event(ctx, "find session request unsuccessful", log.Error(err), log.Data{"status": status})
	writeErrorResponse(ctx, w, body, status)
}
