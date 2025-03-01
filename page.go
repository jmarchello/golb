package main

import (
	"html/template"
	"time"
)

type Page struct {
	Path         string    `yaml:"path"`
	Title        string    `yaml:"title"`
	Date         time.Time `yaml:"date"`
	IsHeaderLink bool      `yaml:"is_header_link"`
	HtmlContent  template.HTML
	MdContent    []byte
	HeaderLinks  []Page
}

func (p *Page) DisplayDate() string {
	return p.Date.Format("01.02.2006")
}
