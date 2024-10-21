package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"sort"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	m "github.com/jmarchello/golb/internal/model"
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

func readMdFiles(basePath string) ([]m.Page, error) {
	var pages []m.Page
	markdownPath := basePath + "/markdown/"

	files, err := os.ReadDir(markdownPath)
	if err != nil {
		return nil, err
	}

	// Loop through each file in the directory
	for _, file := range files {
		// Check if it's a regular file (not a directory)
		if !file.IsDir() {
			// Get the file path
			filePath := markdownPath + file.Name()

			// Read the file
			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}

			pages = append(pages, m.Page{MdContent: data})
		}
	}

	return pages, nil
}

func extractFrontMatter(pages []m.Page) []m.Page {
	re := regexp.MustCompile(`(?s)^---\n(.+?)\n---\n(.*)`)
	var newPages []m.Page

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

func generateHtmlContent(pages []m.Page) []m.Page {
	var newPages []m.Page
	for _, page := range pages {
		page.HtmlContent = template.HTML(mdToHTML(page.MdContent))
		newPages = append(newPages, page)
	}
	return newPages
}

func writePagesToFiles(pages []m.Page, basePath string) error {
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

func buildIndexPage(pages []m.Page, headerLinks []m.Page, basePath string) error {
	listTemplate := template.Must(template.ParseFiles(basePath + "/templates/page_list.html"))
	pageTemplate := template.Must(template.ParseFiles(basePath + "/templates/page.html"))

	listBuffer := &bytes.Buffer{}
	var listPages []m.Page
	for _, page := range pages {
		if !page.IsHeaderLink {
			listPages = append(listPages, page)
		}
	}
	sort.SliceStable(listPages, func(i, j int) bool {
		return listPages[i].Date.After(listPages[j].Date)
	})

	listTemplate.Execute(listBuffer, struct{ Pages []m.Page }{listPages})

	indexPage := m.Page{
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

func setHeaderLinks(pages []m.Page) ([]m.Page, []m.Page) {
	var headerLinks []m.Page
	for _, page := range pages {
		if page.IsHeaderLink {
			headerLinks = append(headerLinks, page)
		}
	}

	var newPages []m.Page
	for _, page := range pages {
		page.HeaderLinks = headerLinks
		newPages = append(newPages, page)
	}

	return newPages, headerLinks
}
