package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"

	"github.com/salpreh/go-wiki/models"
)

var templatesPath string = "views"

var templatesName []string = []string{
	"edit.html",
	"view.html",
}

var templates *template.Template
var validPath *regexp.Regexp

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := models.LoadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", page)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := models.LoadPage(title)

	if err != nil {
		page = &models.Page{Title: title}
	}

	renderTemplate(w, "edit", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
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
	err := templates.ExecuteTemplate(w, templateN+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTemplates() []string {
	var tmplPaths []string
	for _, tmpl := range templatesName {
		tmplPaths = append(tmplPaths, path.Join(templatesPath, tmpl))
	}

	return tmplPaths
}

func getWikiPageTitleOr404(w http.ResponseWriter, r *http.Request) (string, error) {
	title, err := getWikiPageTitle(*r.URL)
	if err != nil {
		http.NotFound(w, r)
	}

	return title, err
}

func getWikiPageTitle(url url.URL) (string, error) {
	m := validPath.FindStringSubmatch(url.Path)
	if m == nil {
		return "", errors.New("Invalid Page title")
	}

	return m[2], nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := getWikiPageTitleOr404(w, r)
		if err != nil {
			return
		}

		fn(w, r, title)
	}
}

func main() {
	templates = template.Must(template.ParseFiles(getTemplates()...))
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
