package models

import (
	"io/ioutil"
	"path"
)

var WikiFolder string = "wiki_pages"
var WikiPathRoot string = "."

type Page struct {
	Title string
	Body  []byte
}

// Saves a page on disk
func (p *Page) Save() error {
	return ioutil.WriteFile(p.getFileName(), p.Body, 0600)
}

func (p *Page) getFileName() string {
	return path.Join(WikiPathRoot, WikiFolder, p.Title+".wk")
}

// Load page from disk
func LoadPage(title string) (*Page, error) {
	page := Page{Title: title, Body: nil}
	body, err := ioutil.ReadFile(page.getFileName())

	if err != nil {
		return nil, err
	}

	page.Body = body

	return &page, nil
}
