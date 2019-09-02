package sessions

import (
	"math/rand"
)

var (
	idChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

type SessionIDGenerator struct {
	r    *rand.Rand
	size int
}

func NewIDGenerator(Random *rand.Rand) *SessionIDGenerator {
	return &SessionIDGenerator{
		r:    Random,
		size: 64,
	}
}

func (gen *SessionIDGenerator) NewID() string {
	id := ""
	for len(id) <= 64 {
		id += string(idChars[gen.r.Intn(len(idChars))])
	}
	return id
}
