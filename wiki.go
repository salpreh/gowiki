package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/salpreh/go-wiki/models"
)

var templatesPath string = "views"

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := getWikiPageTitle(*r.URL)
	page, err := models.LoadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view.html", page)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := getWikiPageTitle(*r.URL)
	page, err := models.LoadPage(title)

	if err != nil {
		page = &models.Page{Title: title}
	}

	renderTemplate(w, "edit.html", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := getWikiPageTitle(*r.URL)
	body := r.FormValue("body")

	page := models.Page{Title: title, Body: []byte(body)}
	err := page.Save()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, templateN string, p *models.Page) {
	t, err := template.ParseFiles(path.Join(templatesPath, templateN))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getWikiPageTitle(url url.URL) string {
	return strings.Join(strings.Split(url.Path, "/")[2:], "/")
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
