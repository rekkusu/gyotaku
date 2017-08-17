package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"goji.io/pat"

	goji "goji.io"

	"github.com/rekkusu/gyotaku/app"
	"github.com/rekkusu/gyotaku/crawler"
)

func main() {
	mux := goji.NewMux()
	app.DBInit()

	mux.Use(app.SessionMiddleware)

	handler := &app.Handler{
		Template: template.Must(template.New("templates").ParseGlob("templates/*.html")),
	}

	mux.HandleFunc(pat.New("/"), handler.Index)
	mux.HandleFunc(pat.New("/new"), handler.NewPage)
	mux.HandleFunc(pat.Get("/view/:id"), handler.View)
	mux.HandleFunc(pat.Get("/admin_53cr37api/"), handler.Admin)

	host := os.Getenv("LISTEN")
	if host == "" {
		host = "127.0.0.1:9999"
	}

	crawler.StartCrawler(8)

	log.Printf("Listening on %s\n", host)

	http.ListenAndServe(host, mux)
}
