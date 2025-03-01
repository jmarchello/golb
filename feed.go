package main

import (
	"errors"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type Feed struct {
	Title   string
	Link    string
	Updated time.Time
	Author  string
	Entries []FeedEntry
}

func (f *Feed) UpdatedString() string {
	return f.Updated.Format(time.RFC3339)
}

func (f *Feed) Render(path string) error {
	// TODO: Maybe do this with marshaling instead
	feedTemplateString := `<?xml version="1.0" encoding="utf-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">

	  <title>{{.Title}}</title>
	  <link href="{{.Link}}"/>
	  <updated>{{.UpdatedString}}</updated>
	  <author>
	    <name>{{.Author}}</name>
	  </author>
	  <id>{{.Link}}/feed.xml</id>

	  {{range .Entries}}
	    <entry>
	     <title>{{.Title}}</title>
	     <link href="{{.Link}}"/>
	     <id>{{.Link}}</id>
	     <updated>{{.UpdatedString}}</updated>
	     <content type="html">{{.EscapedContent}}</content>
	    </entry>
	  {{end}}

	</feed>
	`

	parentDir := filepath.Dir(path)
	fileInfo, err := os.Stat(parentDir)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return errors.New("Provided path does not exists")
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	feedTemplate := template.Must(template.New("feed").Parse(feedTemplateString))
	err = feedTemplate.Execute(file, f)
	if err != nil {
		return err
	}

	return nil
}
