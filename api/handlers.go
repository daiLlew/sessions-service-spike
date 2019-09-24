package api

import (
	"fmt"
	"net/http"

	"github.com/daiLlew/sessions-service-spike/sessions"
	"github.com/gorilla/mux"
)

type SessionCache interface {
	GetByID(id string) (*sessions.Session, error)
	GetByEmail(email string) (*sessions.Session, error)
	Set(*sessions.Session)
}

func CreateSessionHandler(factory *sessions.Factory, cache *sessions.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		details, err := getNewSessionDetails(ctx, r.Body)
		if err != nil {
			writeErrorResponse(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}

		sess := factory.NewSession(details.Email)
		cache.Set(sess)

		created := SessionCreated{
			URI: fmt.Sprintf("/session/%s", sess.ID),
			ID:  sess.ID,
		}

		writeResponse(ctx, w, created, http.StatusCreated)
	}
}

func GetSessionHandler(cache *sessions.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessID := mux.Vars(r)["id"]

		sess, err := cache.GetByID(sessID)
		if err != nil {
			writeErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
		} else if sess == nil {
			writeErrorResponse(ctx, w, "session not found", http.StatusNotFound)
		} else {
			writeResponse(ctx, w, sess, http.StatusOK)
		}
	}
}

func FindSessionHandler(cache *sessions.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userEmail := r.URL.Query().Get("email")
		if len(userEmail) == 0 {
			writeErrorResponse(ctx, w, "bad request", http.StatusBadRequest)
			return
		}

		sess, err := cache.GetByEmail(userEmail)
		if err != nil {
			writeErrorResponse(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}

		if sess == nil {
			writeErrorResponse(ctx, w, "not found", http.StatusNotFound)
			return
		}

		writeResponse(ctx, w, sess, http.StatusOK)
	}
}
