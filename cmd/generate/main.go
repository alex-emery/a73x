package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func Run() error {
	t, err := template.ParseGlob("./templates/layouts/*.html")
	// t, err := template.ParseFiles("index.html", "header.html")
	if err != nil {
		return fmt.Errorf("Failed to parse layouts: %v", err)
	}

	for _, page := range []string{"index.html", "posts.html"} {
		file, err := os.Create(filepath.Join("public", page))
		if err != nil {
			return fmt.Errorf("Failed to create file: %v", err)
		}

		defer file.Close()

		foo, err := t.ParseFiles("./templates/" + page)
		if err != nil {
			return fmt.Errorf("Parse template file: %v", err)
		}

		err = foo.ExecuteTemplate(file, "index.html", nil)
		if err != nil {
			return fmt.Errorf("Failed to generate file: %v", err)
		}
	}

	return nil

}
func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
