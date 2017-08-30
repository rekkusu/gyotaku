package app

import (
	"log"
	"net/http"

	"golang.org/x/net/context"
)

const (
	SessionCookieName = "session"
	SessionKey        = "session"
)

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(SessionCookieName)
		var token string
		if err == nil {
			token = c.Value
		}
		ctx := r.Context()

		var session *Session
		if session, err = GetSession(token); err != nil {
			session, err = NewSession()
			if err != nil {
				panic(err)
			}
			http.SetCookie(w, &http.Cookie{
				Name:  SessionCookieName,
				Value: session.Token,
				Path:  "/",
			})
		}

		ctx = context.WithValue(ctx, SessionKey, session)

		next.ServeHTTP(w, r.WithContext(ctx))

		SaveSession(session)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}
