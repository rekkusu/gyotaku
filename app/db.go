package app

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"

	"golang.org/x/sync/syncmap"
)

var (
	Sessions *syncmap.Map
	Pages    []*Page
	db       *gorm.DB
)

type Session struct {
	Token   string `gorm:"primary_key"`
	Message string
	Pages   []Page `gorm:"ForeignKey:SessionID"`
}

type Page struct {
	Id        int    `gorm:"primary_key"`
	SessionID string `gorm:"index"`
	Url       string
	Body      string
}

func DBInit() {
	Sessions = new(syncmap.Map)
	var err error
	db, err = gorm.Open("sqlite3", "fbf04f53a0915945a74e9b7d5d2451c8.sqlite")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Session{}, &Page{})
	db.CreateTable(&Session{}, &Page{})
}

func NewSession() (*Session, error) {
	token := generateSessionKey()

	session := &Session{
		Token: token,
		Pages: make([]Page, 0),
	}

	if err := db.Save(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func generateSessionKey() string {
	for {
		r := make([]byte, 32)
		if _, err := rand.Read(r); err != nil {
			panic(err)
		}

		token := hex.EncodeToString(r)
		if err := db.First(&Session{}, "token = ?", token).Error; err != nil {
			return token
		}
	}

}

func GetSession(key string) (*Session, error) {
	session := &Session{}
	var pages []Page
	if err := db.First(session, "token = ?", key).Related(&pages, "SessionID").Error; err != nil {
		return nil, err
	}
	session.Pages = pages
	return session, nil
}

func SaveSession(session *Session) error {
	if err := db.Set("gorm:save_associations", false).Save(session).Error; err != nil {
		return err
	}
	return nil
}

func GetPage(id int) (*Page, error) {
	var page Page
	err := db.First(&page, id).Error
	return &page, err
}

func CreatePage(page *Page) error {
	return db.Create(page).Error
}

func SavePage(page *Page) error {
	return db.Save(page).Error
}

func GetDefaultSession(r *http.Request) *Session {
	return r.Context().Value(SessionKey).(*Session)
}
