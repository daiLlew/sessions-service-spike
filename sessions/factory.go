package sessions

import (
	"math/rand"
	"time"
)

var (
	idChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

type Factory struct {
	random *rand.Rand
}

func NewFactory() *Factory {
	source := rand.NewSource(time.Now().UnixNano())
	return &Factory{random: rand.New(source)}
}

func (factory *Factory) NewSession(email string) *Session {
	now := time.Now()

	return &Session{
		ID:           factory.NewID(),
		Email:        email,
		Start:        now,
		LastAccessed: now,
	}
}

func (factory *Factory) NewID() string {
	id := ""
	for len(id) <= 64 {
		id += string(idChars[factory.random.Intn(len(idChars))])
	}
	return id
}
