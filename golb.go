package main

import (
  "fmt"
  "os"
  "errors"
  "time"
  "io/ioutil"

  "github.com/gomarkdown/markdown"
  "github.com/gomarkdown/markdown/html"
  "github.com/gomarkdown/markdown/parser"
)

type Page struct {
  Path string
  Title string
  Date time.Time
  Content []byte
  IsHeaderLink bool
  HeaderLinks []*Page
}

func main() {
  _, err := checkArgs()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  basePath := os.Args[1]
  markdownData, err := readMdFiles(basePath)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }


  fmt.Println(len(markdownData))
  fmt.Println(markdownData[0])
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

func readMdFiles(basePath string) ([][]byte, error) {
    var fileContents [][]byte
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

        fileContents = append(fileContents, data)
      }
    }

    return fileContents, nil
}
