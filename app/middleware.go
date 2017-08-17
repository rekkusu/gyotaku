package app

import (
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
			if session, ok := Sessions[c.Value]; ok {
				ctx = context.WithValue(ctx, SessionKey, session)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		for {
			session, _ := NewSession()
			if _, ok := Sessions[session.Token]; !ok {
				Sessions[session.Token] = session
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
