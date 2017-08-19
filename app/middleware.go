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
		ctx := r.Context()
		if err == nil {
			if session, ok := Sessions.Load(c.Value); ok {
				ctx = context.WithValue(ctx, SessionKey, session)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		for {
			session, _ := NewSession()
			if _, ok := Sessions.Load(session.Token); !ok {
				Sessions.Store(session.Token, session)
				ctx = context.WithValue(ctx, SessionKey, session)
				http.SetCookie(w, &http.Cookie{
					Name:  SessionCookieName,
					Value: session.Token,
					Path:  "/",
				})
				break
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}
