package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Custom template type that wraps around the native template type
type Template struct {
	HTMLTpl *template.Template
}

// Must  wraps a call to a function returning (Template, error) and panics if the error is non-nil
func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

// Parse parses the template located in the filepath
func Parse(filepath string) (Template, error) {
	// If there is an error parsing, it will be handled here (i.e. invalid function in the template)
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	// If there is no error, then return
	return Template{HTMLTpl: tpl}, nil
}

// Execute executes a template of type <Template> that is already parsed
func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	// Sets the content type of the response header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Executes the template t without any data (nil)
	// If there is an error rendering, it will be handled here (i.e. invalid field in the template)
	// This approach writes to the response writer until an error is detected (if any). If an error
	// is detected half-way through the execution, then the webpage will be half rendered
	err := t.HTMLTpl.Execute(w, nil)
	if err != nil {
		log.Printf("executing templates: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}
