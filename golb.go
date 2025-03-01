package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gopkg.in/yaml.v3"
)

func main() {
	_, err := checkArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	basePath := os.Args[1]
	pages, err := readMdFiles(basePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pages = extractFrontMatter(pages)
	pages, headerLinks := setHeaderLinks(pages)
	pages = generateHtmlContent(pages)
	err = writePagesToFiles(pages, basePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = buildIndexPage(pages, headerLinks, basePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = writeAtomFeed(pages, basePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkArgs() (bool, error) {
	if len(os.Args) < 2 {
		return false, errors.New("No build directory provided")
	}

	fileInfo, err := os.Stat(os.Args[1])
	if err != nil {
		return false, err
	}

	if !fileInfo.IsDir() {
		return false, errors.New("Provided path is not a directory")
	}

	return true, nil
}

func mdToHTML(md []byte) string {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

func readMdFiles(basePath string) ([]Page, error) {
	var pages []Page
	markdownPath := basePath + "/markdown/"

	files, err := ioutil.ReadDir(markdownPath)
	if err != nil {
		return nil, err
	}

	// Loop through each file in the directory
	for _, file := range files {
		// Check if it's a regular file (not a directory)
		if file.Mode().IsRegular() {
			// Get the file path
			filePath := markdownPath + file.Name()

			// Read the file
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}

			pages = append(pages, Page{MdContent: data})
		}
	}

	return pages, nil
}

func extractFrontMatter(pages []Page) []Page {
	re := regexp.MustCompile(`(?s)^---\n(.+?)\n---\n(.*)`)
	var newPages []Page

	for _, page := range pages {
		match := re.FindSubmatch(page.MdContent)
		err := yaml.Unmarshal(match[1], &page)
		if err != nil {
			fmt.Println(err)
		}
		page.MdContent = match[2]
		newPages = append(newPages, page)
	}
	return newPages
}

func generateHtmlContent(pages []Page) []Page {
	var newPages []Page
	for _, page := range pages {
		page.HtmlContent = template.HTML(mdToHTML(page.MdContent))
		newPages = append(newPages, page)
	}
	return newPages
}

func writePagesToFiles(pages []Page, basePath string) error {
	os.RemoveAll(basePath + "/site")
	err := os.Mkdir(basePath+"/site", 0750)
	if err != nil {
		return err
	}

	t := template.Must(template.ParseFiles(basePath + "/templates/page.html"))

	for _, page := range pages {
		filePath := fmt.Sprintf("%v/site/%v.html", basePath, page.Path)
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		err = t.Execute(f, page)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildIndexPage(pages []Page, headerLinks []Page, basePath string) error {
	listTemplate := template.Must(template.ParseFiles(basePath + "/templates/page_list.html"))
	pageTemplate := template.Must(template.ParseFiles(basePath + "/templates/page.html"))

	listBuffer := &bytes.Buffer{}
	var listPages []Page
	for _, page := range pages {
		if !page.IsHeaderLink {
			listPages = append(listPages, page)
		}
	}
	sort.SliceStable(listPages, func(i, j int) bool {
		return listPages[i].Date.After(listPages[j].Date)
	})

	listTemplate.Execute(listBuffer, struct{ Pages []Page }{listPages})

	indexPage := Page{
		Title:       "Josh Marchello",
		HtmlContent: template.HTML(listBuffer.String()),
		HeaderLinks: headerLinks,
	}

	filePath := fmt.Sprintf("%v/site/index.html", basePath)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = pageTemplate.Execute(f, indexPage)
	if err != nil {
		return err
	}

	return nil
}

func setHeaderLinks(pages []Page) ([]Page, []Page) {
	var headerLinks []Page
	for _, page := range pages {
		if page.IsHeaderLink {
			headerLinks = append(headerLinks, page)
		}
	}

	var newPages []Page
	for _, page := range pages {
		page.HeaderLinks = headerLinks
		newPages = append(newPages, page)
	}

	return newPages, headerLinks
}

func writeAtomFeed(pages []Page, basePath string) error {
	// TODO: get these values from a config or something
	feed := Feed{
		Title:   "Josh Marchello",
		Link:    "https://jmarchello.com",
		Updated: time.Now(),
		Author:  "Josh Marchello",
	}

	for _, page := range pages {
		entry := FeedEntry{
			Title:   page.Title,
			Link:    feed.Link + "/" + page.Path,
			Updated: page.Date,
			Content: page.HtmlContent,
		}
		feed.Entries = append(feed.Entries, entry)
	}
	err := feed.Render(basePath + "/site/feed.xml")
	return err
}
