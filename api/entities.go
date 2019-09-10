package api

type NewSessionDetails struct {
	Email string `json:"email"`
}

type SessionCreated struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}