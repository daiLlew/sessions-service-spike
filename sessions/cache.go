package sessions

import (
	"sync"
	"time"

	"github.com/ONSdigital/log.go/log"
)

var mutex = &sync.Mutex{}

type Cache struct {
	ttl      time.Duration
	interval time.Duration
	store    map[string]*Session
}

func NewCache(interval time.Duration, ttl time.Duration) *Cache {
	cache := &Cache{
		ttl:      ttl,
		interval: interval,
		store:    map[string]*Session{},
	}

	go runPurger(cache)

	log.Event(nil, "session cache created")
	return cache
}

func (c *Cache) purge() {
	log.Event(nil, "executing session cache purge")
	if len(c.store) == 0 {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for id, sess := range c.store {
		if c.isExpired(sess) {
			log.Event(nil, "purging expired session", log.Data{"email": sess.Email})
			delete(c.store, id)
		}
	}
}

func (c *Cache) GetByID(id string) (*Session, error) {
	mutex.Lock()
	defer mutex.Unlock()

	sess, ok := c.store[id]
	if !ok {
		log.Event(nil, "session not found")
		return nil, nil
	}

	if c.isExpired(sess) {
		log.Event(nil, "session expired")
		delete(c.store, id)
		return nil, nil
	}

	sess.LastAccessed = time.Now()
	c.store[id] = sess
	return sess, nil
}

func (c *Cache) GetByEmail(email string) (*Session, error) {
	mutex.Lock()
	defer mutex.Unlock()

	findByEmail := func(s *Session) bool {
		return s.Email == email
	}

	sess := c.findSessionBy(findByEmail)
	if sess == nil {
		return nil, nil
	}

	if c.isExpired(sess) {
		log.Event(nil, "session expired")
		delete(c.store, sess.ID)
		return nil, nil
	}

	return sess, nil
}

func (c *Cache) Set(sess *Session) {
	mutex.Lock()
	defer mutex.Unlock()
	log.Event(nil, "adding session to cache")
	c.store[sess.ID] = sess
}

func (c *Cache) isExpired(sess *Session) bool {
	return time.Since(sess.LastAccessed) >= c.ttl
}

func (c *Cache) findSessionBy(filterFunc func(s *Session) bool) *Session {
	for _, sess := range c.store {
		if filterFunc(sess) {
			return sess
		}
	}
	return nil
}

func runPurger(cache *Cache) {
	purgeTicker := time.NewTicker(cache.interval)

	for {
		select {
		case <-purgeTicker.C:
			cache.purge()
		}
	}
}
