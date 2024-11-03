package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"git.sr.ht/~a73x/home/pages"
)

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	pages, err := pages.Collect("content")
	if err != nil {
		return err
	}

	var errs []error
	for _, page := range pages {
		err = writeFile(path.Join("public", page.Path), []byte(page.Content))
		if err != nil {
			errs = append(errs, err)
		}
	}

	if errs != nil {
		return errors.Join(errs...)
	}

	return nil
}

func writeFile(name string, contents []byte) error {
	folders := path.Dir(name)
	_, err := os.Stat(folders)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(folders, 0744); err != nil {
			return fmt.Errorf("failed to mkdir %s\n%w", folders, err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat folder %s\n%w", folders, err)
	}

	err = os.WriteFile(name, contents, 0666)
	if err != nil {
		return fmt.Errorf("failed to write file %s\n%w", name, err)
	}

	return nil
}
