package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

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
		Template: template.New("templates"),
	}

	for _, name := range AssetNames() {
		handler.Template.New(name).Parse(string(MustAsset(name)))
	}

	mux.HandleFunc(pat.New("/"), handler.Index)
	mux.HandleFunc(pat.New("/new"), handler.NewPage)
	mux.HandleFunc(pat.Get("/view/:id"), handler.View)

	// For crawlers
	mux.HandleFunc(pat.Get(app.SecretURL), handler.Admin)

	conf := config()

	crawler.ChromePath = conf.Chrome
	crawler.Flag = conf.Flag
	crawler.StartCrawler(conf.CrawlerJobs)

	log.Printf("Listening on %s\n", conf.Listen)

	http.ListenAndServe(conf.Listen, mux)
}

type Config struct {
	Listen      string
	Chrome      string
	CrawlerJobs int
	Flag        string
}

func config() Config {
	conf := Config{}

	if conf.Listen = os.Getenv("LISTEN"); conf.Listen == "" {
		conf.Listen = "127.0.0.1:9999"
	}

	if conf.Chrome = os.Getenv("CHROME"); conf.Chrome == "" {
		conf.Chrome = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	}

	jobs, err := strconv.Atoi(os.Getenv("CRAWLER_JOBS"))
	if err == nil {
		conf.CrawlerJobs = jobs
	} else {
		conf.CrawlerJobs = 1
	}

	conf.Flag = app.Flag

	return conf
}
