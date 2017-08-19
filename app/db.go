package app

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"golang.org/x/sync/syncmap"
)

var (
	Sessions *syncmap.Map
	Pages    []*Page
)

type Session struct {
	Token   string
	User    *User
	Message string
}

type User struct {
	Pages []*Page
}

type Page struct {
	Url  string
	Body string
}

func DBInit() {
	Sessions = new(syncmap.Map)
}

func NewSession() (*Session, error) {
	r := make([]byte, 32)
	if _, err := rand.Read(r); err != nil {
		return nil, err
	}

	token := hex.EncodeToString(r)

	return &Session{
		Token: token,
		User: &User{
			Pages: make([]*Page, 0),
		},
	}, nil
}

func GetSession(r *http.Request) *Session {
	return r.Context().Value(SessionKey).(*Session)
}
