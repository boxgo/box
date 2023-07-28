package ginsession

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/boxgo/redisstore/v2"
	"github.com/boxgo/redisstore/v2/serializer"
	"github.com/gorilla/sessions"
)

// MockSet set cookie and return redis sid key and cookie string
func (s *GinSession) MockSet(value interface{}) (sid string, cookie string, err error) {
	var (
		store   *redisstore.RedisStore
		session *sessions.Session
		req     *http.Request
		rsp     = httptest.NewRecorder()
	)

	if store, err = redisstore.NewStoreWithUniversalClient(
		s.client.Client(),
		redisstore.WithMaxLength(s.cfg.MaxLen),
		redisstore.WithKeyPrefix(s.cfg.KeyPrefix),
		redisstore.WithKeyPairs([]byte(s.cfg.KeyPair)),
		redisstore.WithSerializer(serializer.JSONSerializer{}),
	); err != nil {
		return
	}

	if req, err = http.NewRequest("GET", "http://ginsession_mock", nil); err != nil {
		return
	}

	if session, err = store.New(req, s.CookieName()); err != nil {
		return
	}

	session.Values["user"] = value

	if err = session.Save(req, rsp); err != nil {
		return
	}

	return fmt.Sprintf("%s%s", s.cfg.KeyPrefix, session.ID), rsp.Header().Get("Set-Cookie"), nil
}

// MockDel del session by sid
func (s *GinSession) MockDel(sid string) (err error) {
	return s.client.Client().Del(context.TODO(), sid).Err()
}
