package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"regexp"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gopkg.in/yaml.v3"
)

type Page struct {
	Path         string    `yaml:"path"`
	Title        string    `yaml:"title"`
	Date         time.Time `yaml:"date"`
	IsHeaderLink bool      `yaml:"is_header_link"`
	HtmlContent  []byte
	MdContent    []byte
	HeaderLinks  []*Page
}

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
	fmt.Printf("%+v\n", pages[0])
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

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
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

func extractFrontMatter(pages []Page) ([]Page) {
	re := regexp.MustCompile(`(?s)^---\n(.+?)\n---\n`)
	var newPages []Page

	for _, page := range pages {
		match := re.FindSubmatch(page.MdContent)
		err := yaml.Unmarshal(match[1], &page)
		if err != nil {
			fmt.Println(err)
		}
		newPages = append(newPages, page)
	}
	return newPages
}