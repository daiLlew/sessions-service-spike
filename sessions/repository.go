package sessions

import "errors"

var (
	IDEmptyError      = errors.New("session id required but was empty")
	SessionEmptyError = errors.New("session required but was nil")
)

type Repository struct {
	Store map[string]*Session
}

func NewRepository() *Repository {
	return &Repository{Store: make(map[string]*Session, 0)}
}

func (repo *Repository) Save(sess *Session) error {
	if sess == nil {
		return SessionEmptyError
	}

	repo.Store[sess.ID] = sess
	return nil
}

func (repo *Repository) GetByID(id string) (*Session, error) {
	if len(id) == 0 {
		return nil, IDEmptyError
	}
	sess := repo.Store[id]
	return sess, nil
}
