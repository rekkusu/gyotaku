package app

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/rekkusu/gyotaku/crawler"
	"github.com/rekkusu/gyotaku/secret"

	"goji.io/pat"
)

type Handler struct {
	Template *template.Template
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	session := GetDefaultSession(r)
	h.Template.ExecuteTemplate(w, "index.html", struct {
		Session *Session
	}{
		Session: session,
	})
	session.Message = ""
}

func (h *Handler) NewPage(w http.ResponseWriter, r *http.Request) {
	session := GetDefaultSession(r)

	if session.Token == secret.AdminToken {
		session.Message = "Admin cannot post a gyotaku"
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	url := r.PostFormValue("url")

	page := &Page{
		SessionID: session.Token,
		Url:       url,
		Body:      "Now loading...",
	}
	CreatePage(page)

	go func() {
		body := GetWebPage(page.Url)

		if len([]rune(body)) > 8192 {
			body = string([]rune(body)[:8192])
		}

		page.Body = body
		SavePage(page)

		//Pages = append(Pages, page)
		crawler.CrawlQueue <- page.Id
	}()

	session.Message = "Added"

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(pat.Param(r, "id"))
	session := GetDefaultSession(r)

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

	if id >= len(session.Pages) || id < 0 {
		w.WriteHeader(http.StatusNotFound)
		h.Template.ExecuteTemplate(w, "view.html", tpl)
		return
	}

	page := session.Pages[id]
	tpl.Url = template.HTML(page.Url)
	tpl.Body = page.Body

	h.Template.ExecuteTemplate(w, "view.html", tpl)
}

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	session := GetDefaultSession(r)
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

	p, err := GetPage(id)
	if err != nil {
		log.Fatalln(err)
	}

	tpl.Body = p.Body
	tpl.Url = template.HTML(p.Url)

	h.Template.ExecuteTemplate(w, "view.html", tpl)
}
