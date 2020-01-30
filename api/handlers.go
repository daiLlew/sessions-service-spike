package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ONSdigital/log.go/log"
	"github.com/daiLlew/sessions-service-spike/sessions"
	"github.com/gorilla/mux"
)

var (
	BadRequestErr   = errors.New("bad request")
	FlushAllSuccess = SimpleMessage{Message: "flushall completed successfully"}
)

func CreateSessionHandler(factory *sessions.Factory, cache SessionCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessCreated, err := createNewSession(r, factory, cache)
		if err != nil {
			handleCreateSessionError(r.Context(), w, err)
		} else {
			writeResponse(r.Context(), w, sessCreated, http.StatusCreated)
		}
	}
}

func GetSessionHandler(cache SessionCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessID := mux.Vars(r)["id"]

		sess, err := cache.GetByID(sessID)
		if err != nil {
			handleGetSessionError(ctx, w, err)
		} else {
			writeResponse(ctx, w, sess, http.StatusOK)
		}
	}
}

func FindSessionHandler(cache SessionCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sess, err := findSession(cache, r)
		if err != nil {
			handleFindSessionError(ctx, w, err)
		} else {
			writeResponse(ctx, w, sess, http.StatusOK)
		}
	}
}

func FlushCashHandler(cache SessionCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := cache.FlushAll()
		if err != nil {
			handleFlushAllError(ctx, w, err)
		} else {
			writeResponse(ctx, w, FlushAllSuccess, 200)
		}
	}
}

func createNewSession(r *http.Request, factory *sessions.Factory, cache SessionCache) (*SessionCreated, error) {
	ctx := r.Context()

	details, err := getNewSessionDetails(ctx, r.Body)
	if err != nil {
		return nil, err
	}

	log.Event(nil, "creating new session", log.Data{"details": details})

	sess := factory.NewSession(details.Email)
	if err := cache.Set(sess); err != nil {
		return nil, err
	}

	log.Event(nil, "session saved successfully")

	sessCreated := &SessionCreated{
		URI: fmt.Sprintf("/session/%s", sess.ID),
		ID:  sess.ID,
	}
	return sessCreated, nil
}

func findSession(cache SessionCache, r *http.Request) (*sessions.Session, error) {
	ctx := r.Context()
	userEmail := r.URL.Query().Get("email")
	logD := log.Data{"email": userEmail}

	log.Event(ctx, "finding session by email", logD)
	if len(userEmail) == 0 {
		return nil, BadRequestErr
	}

	return cache.GetByEmail(userEmail)
}
