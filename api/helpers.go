package api

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/log.go/log"
)

func getNewSessionDetails(ctx context.Context, body io.ReadCloser) (*NewSessionDetails, error) {
	defer body.Close()

	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Event(ctx, "failed to read request body", log.Error(err))
		return nil, err
	}

	var details NewSessionDetails
	err = json.Unmarshal(b, &details)
	if err != nil {
		log.Event(ctx, "failed to unmarshal request body", log.Error(err))
		return nil, err
	}
	return &details, nil
}

func writeResponse(ctx context.Context, w http.ResponseWriter, entity interface{}, status int) {
	b, err := json.Marshal(entity)
	if err != nil {
		log.Event(ctx, "error marshalling response body", log.Error(err))
		writeErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, body string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(body))
}
