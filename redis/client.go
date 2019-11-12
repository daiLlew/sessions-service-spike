package redis

import (
	"errors"
	"time"

	"github.com/ONSdigital/log.go/log"
	"github.com/daiLlew/sessions-service-spike/sessions"
	"github.com/gomodule/redigo/redis"
)

const TTL = time.Second * 30

type Client struct {
	pool redis.Pool
}

func NewCli() *Client {
	return &Client{pool: redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 1 * time.Minute,
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp", ":6379")
		},
	}}
}

func (c *Client) GetTTL(id string) (int, error) {
	conn := c.pool.Get()
	defer close(conn)

	return redis.Int(conn.Do("TTL", id))
}

func (c *Client) GetByID(id string) (*sessions.Session, error) {
	conn := c.pool.Get()
	defer close(conn)

	s, err := get(conn, id)
	if err != nil {
		return nil, err
	}

	if err := c.updateExpireAt(conn, s); err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Client) GetByEmail(email string) (*sessions.Session, error) {
	conn := c.pool.Get()
	defer close(conn)

	s, err := get(conn, email)
	if err != nil {
		return nil, err
	}

	if err := c.updateExpireAt(conn, s); err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Client) Set(sess *sessions.Session) error {
	if sess == nil {
		return errors.New("session was nil")
	}

	conn := c.pool.Get()
	defer close(conn)

	expireAt := time.Now().Add(TTL).Unix()

	conn.Send("MULTI")

	conn.Send("HSET", sess.ID, "id", sess.ID, "email", sess.Email)
	conn.Send("EXPIREAT", sess.ID, expireAt)

	conn.Send("HSET", sess.Email, "email", sess.Email, "id", sess.ID)
	conn.Send("EXPIREAT", sess.Email, expireAt)

	_, err := conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

func get(conn redis.Conn, key interface{}) (*sessions.Session, error) {
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return nil, err
	}

	if !exists {
		log.Event(nil, "session not found")
		return nil, sessions.SessionNotFoundErr
	}

	values, err := redis.Values(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}

	sess, err := redis.Values(values, nil)
	if err != nil {
		return nil, err
	}

	var s sessions.Session
	err = redis.ScanStruct(sess, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (c *Client) updateExpireAt(conn redis.Conn, s *sessions.Session) error {
	expireAt := time.Now().Add(TTL).Unix()

	conn.Send("MULTI")
	conn.Send("EXPIREAT", s.ID, expireAt)
	conn.Send("EXPIREAT", s.Email, expireAt)

	_, err := conn.Do("EXEC")
	return err
}

func close(conn redis.Conn) {
	if err := conn.Close(); err != nil {
		log.Event(nil, "error closing redis conn", log.Error(err))
	}
}
