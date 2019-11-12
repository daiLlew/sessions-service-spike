package sessions

import (
	"encoding/json"
	"time"
)

const (
	dateTimeFMT = "2006-01-02T15:04:05.000Z"
)

type IDGenerator interface {
	NewID() string
}

type Session struct {
	ID           string    `json:"id" redis:"id"`
	Email        string    `json:"email" redis:"email"`
	Start        time.Time `json:"start" redis:"start"`
	LastAccessed time.Time `json:"lastAccess" redis:"lastAccess"`
}

type jsonModel struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Start        string `json:"start"`
	LastAccessed string `json:"lastAccess"`
}

func (sess *Session) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonModel{
		ID:           sess.ID,
		Email:        sess.Email,
		Start:        sess.Start.Format(dateTimeFMT),
		LastAccessed: sess.LastAccessed.Format(dateTimeFMT),
	})
}
