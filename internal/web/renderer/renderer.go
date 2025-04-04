package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/yuin/goldmark"
)

func LoadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"markdown": func(input string) template.HTML {
			var buf bytes.Buffer
			if err := goldmark.Convert([]byte(input), &buf); err != nil {
				return template.HTML(input)
			}
			return template.HTML(buf.String())
		},
		"successRate": func(success, total int) string {
			if total == 0 {
				return "0"
			}
			return fmt.Sprintf("%.2f", float64(success)/float64(total)*100)
		},
		"status": func(status string) string {
			status = strings.ToLower(status)
			if status == "accepted" || status == "pending" {
				return status
			}
			if status == "compile error" {
				return "compile_error"
			}
			return "error"
		},
		"initial": func(input string) string {
			return input[:1]
		},
	}

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	partials, err := filepath.Glob(templatesDir + "/partials/*.html")
	if err != nil {
		panic(err.Error())
	}

	pages, err := filepath.Glob(templatesDir + "/pages/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, page := range pages {
		files := append(layouts, page)
		files = append(files, partials...)
		r.AddFromFilesFuncs(filepath.Base(page), funcMap, files...)
	}
	return r
}
