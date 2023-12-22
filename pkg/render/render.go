package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"reservation_program/pkg/config"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string) {

	var tc map[string]*template.Template
	if app.UseCache {
		// Create a template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// Get the requested template
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	// Check for error in the value stored in the map
	err := t.Execute(buf, nil)
	if err != nil {
		log.Println(err)
	}

	// Render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// Get all the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// Range through all files ending with *.page.tmpl
	for _, page := range pages {
		// Get the name of the file
		name := filepath.Base(page)
		// Create a pointer to template, name the template, and parse the file page
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// Look for a layout file
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		// If there is a layout file, parse it
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		// Add the template to the map
		myCache[name] = ts
	}

	// Return the map
	return myCache, nil
}
