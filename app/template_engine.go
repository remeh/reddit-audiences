package app

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

// ReadTemplates looks and reads in the configured directory
// for all available templates.
func ReadTemplates(app *App) (*template.Template, error) {
	rv, err := lookForTemplates(app.Config.TemplatesDir)
	if err != nil {
		return nil, err
	}

	if len(rv) == 0 {
		return nil, fmt.Errorf("no templates available in %s", app.Config.TemplatesDir)
	}

	// assign custom functions and parse templates.
	t, err := template.New("").Funcs(TemplateHelpers()).ParseFiles(rv...)
	if err != nil {
		return nil, err
	}

	// attach method helpers
	return t, nil
}

func lookForTemplates(path string) ([]string, error) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	rv := make([]string, 0)

	dir, err := os.Open(path)
	if err != nil {
		return rv, err
	}

	if stat, err := dir.Stat(); err != nil {
		return rv, err
	} else {
		if !stat.IsDir() {
			return rv, fmt.Errorf("the templates directory isn't a directory!")
		}
	}

	fileInfos, err := dir.Readdir(0)
	if err != nil {
		return rv, err
	}

	for _, fi := range fileInfos {
		// ignore directory
		if fi.IsDir() {
			continue
		}

		if strings.HasSuffix(fi.Name(), ".html") {
			rv = append(rv, path+fi.Name())
		}
	}

	return rv, nil
}
