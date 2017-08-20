package app

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/rekkusu/gyotaku/crawler"

	"goji.io/pat"
)

type Handler struct {
	Template *template.Template
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	session := GetSession(r)
	h.Template.ExecuteTemplate(w, "index.html", struct {
		Session *Session
	}{
		Session: session,
	})
	session.Message = ""
}

func (h *Handler) NewPage(w http.ResponseWriter, r *http.Request) {
	session := GetSession(r)

	url := r.FormValue("url")

	if strings.HasPrefix(strings.ToLower(url), "file://") {
		session.Message = "Invalid URL"
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	page := &Page{
		Url:  url,
		Body: "Loading",
	}

	go func() {
		body := GetWebPage(page.Url)

		if len([]rune(body)) > 8192 {
			body = string([]rune(body)[:8192])
		}

		page.Body = body

		Pages = append(Pages, page)
		crawler.CrawlQueue <- len(Pages) - 1
	}()

	session.Message = "Added"
	session.User.Pages = append(session.User.Pages, page)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(pat.Param(r, "id"))
	session := GetSession(r)

	var tpl struct {
		Session *Session
		Url     template.HTML
		Body    string
	}

	tpl.Session = session
	tpl.Url = template.HTML("Not Found")

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)

		h.Template.ExecuteTemplate(w, "view.html", tpl)
		return
	}

	if id >= len(session.User.Pages) || id < 0 {
		w.WriteHeader(http.StatusNotFound)
		h.Template.ExecuteTemplate(w, "view.html", tpl)
		return
	}

	page := session.User.Pages[id]
	tpl.Url = template.HTML(page.Url)
	tpl.Body = page.Body

	h.Template.ExecuteTemplate(w, "view.html", tpl)
}

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	session := GetSession(r)
	id, err := strconv.Atoi(pat.Param(r, "id"))

	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	var tpl struct {
		Session *Session
		Url     template.HTML
		Body    string
	}

	tpl.Session = session
	tpl.Url = template.HTML("Not Found")

	p := Pages[id]
	tpl.Body = p.Body
	tpl.Url = template.HTML(p.Url)

	h.Template.ExecuteTemplate(w, "view.html", tpl)
}
