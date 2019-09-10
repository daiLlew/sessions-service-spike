package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/log.go/log"
	"github.com/daiLlew/sessions-service-spike/sessions"
	"github.com/gorilla/mux"
)

type SessionCache interface {
	Get(id string) (*sessions.Session, error)
	Set(*sessions.Session)
}

func CreateSessionHandler(factory *sessions.Factory, cache *sessions.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		details, err := getNewSessionDetails(ctx, r.Body)
		if err != nil {
			writeErrorResponse(ctx, w, err.Error(), http.StatusBadRequest)
		} else {

			sess := factory.NewSession(details.Email)
			cache.Set(sess)

			created := SessionCreated{
				URI: fmt.Sprintf("/session/%s", sess.ID),
				ID:  sess.ID,
			}

			writeResponse(ctx, w, created, http.StatusCreated)
		}
	}
}

func GetSessionHandler(cache *sessions.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessID := mux.Vars(r)["id"]

		sess, err := cache.Get(sessID)
		if err != nil {
			writeErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
		} else if sess == nil {
			writeErrorResponse(ctx, w, "session not found", http.StatusNotFound)
		} else {
			writeResponse(ctx, w, sess, http.StatusOK)
		}
	}
}

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
