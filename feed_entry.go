package main

import (
	"html"
	"html/template"
	"time"
)

type FeedEntry struct {
	Title   string
	Link    string
	Updated time.Time
	Content template.HTML
}

func (e *FeedEntry) EscapedContent() string {
	return html.EscapeString(string(e.Content))
}

func (e *FeedEntry) UpdatedString() string {
	return e.Updated.Format(time.RFC3339)
}
