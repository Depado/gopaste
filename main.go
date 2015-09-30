package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/Depado/gopaste/configuration"
	"github.com/GeertJohan/go.rice"
	"github.com/goji/param"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nu7hatch/gouuid"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var templateBox = rice.MustFindBox("templates")
var staticBox = rice.MustFindBox("assets")
var db gorm.DB

// Paste is a paste
type Paste struct {
	Title    string
	Content  string
	Markdown bool
}

// PasteEntry is the db model
type PasteEntry struct {
	gorm.Model
	Paste
	Key string
}

func loadTemplate(path string) (*template.Template, error) {
	templateString, err := templateBox.String(path)
	if err != nil {
		return nil, err
	}
	t, err := template.New(path).Parse(templateString)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func index(c web.C, w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method == "POST" {
		err = r.ParseForm()
		if err != nil {
			http.Error(w, "No good!", 400)
			return
		}
		p := Paste{}
		err = param.Parse(r.PostForm, &p)
		if err != nil {
			http.Error(w, "You're doing it wrong. In some way at least.", 500)
			return
		}
		pe := PasteEntry{}
		pe.Paste = p
		u, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "Something went wrong.", 500)
		}
		pe.Key = u.String()
		db.Create(&pe)
		http.Redirect(w, r, fmt.Sprintf("/paste/%s", pe.Key), http.StatusFound)
		return
	}

	t, err := loadTemplate("index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, c.Env)
}

func getPaste(c web.C, w http.ResponseWriter, r *http.Request) {
	pe := PasteEntry{}
	db.Where(&PasteEntry{Key: c.URLParams["key"]}).First(&pe)
	if pe.Key == "" {
		http.Error(w, "Not found", 404)
	}
	fmt.Fprint(w, pe.Content)
}

func main() {
	var err error

	err = configuration.Load("conf.yml")
	if err != nil {
		log.Fatal(err)
	}

	db, err = gorm.Open("sqlite3", "my.db")
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&PasteEntry{})

	goji.Get("/paste/:key", getPaste)
	goji.Get("/static/*", http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox())))
	goji.Get("/", index)
	goji.Post("/", index)

	flag.Set("bind", ":"+configuration.Config.Port)
	goji.Serve()
}
