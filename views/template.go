package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/csrf"
	"github.com/wagnojunior/lenslocked/context"
	"github.com/wagnojunior/lenslocked/models"
)

// public defines a public error message
type public interface {
	Public() string
}

// Custom template type that wraps around the native template type
type Template struct {
	HTMLTpl *template.Template
}

// Must wraps a call to a function returning (Template, error) and panics if
// the error is non-nil
func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

// ParseFS parses the template located in the file system fs
func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	// We need to define the template functions BEFORE the templates are parsed.
	// To do that we first create an empty template with the name of the first
	// pattern. Then, we add a *placeholder* function to it. This temporaty
	// function will latter be rewritten when the templace is executed.
	tpl := template.New(path.Base(patterns[0]))
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("scrfField not implemented")

			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser not implemented")

			},
			"errors": func() []string {
				return nil
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	// If there is no error, then return
	return Template{HTMLTpl: tpl}, nil
}

// Execute executes a template of type <Template> that is already parsed
func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	// Clones the original template to avoid *race condition* where many users
	// could be requesting from the same template at the same time
	tpl, err := t.HTMLTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}

	// Gets all the error messages
	errMessages := errMessages(errs...)

	// This function replaces the placeholder function from `ParseFS`
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)

			},
			"currentUser": func() *models.User {
				return context.User(r.Context())

			},
			"errors": func() []string {
				return errMessages
			},
		},
	)

	// Sets the content type of the response header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Executes the template t with data. If there is an error rendering, it
	// will be handled here (i.e. invalid field in the template). This approach
	// writes to the buffer until an error is detected (if any). If an error is
	// detected half-way through the execution, then the webpage will not be
	// rendered
	var buf bytes.Buffer

	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing templates: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}

// errMessages retuns a slice of error messages. If an error has a `Public`
// message associated with it, then the public message is displayed. Otherwise,
// a standard message is used.
func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		// Checks if the error is of type `public` and displays it.
		// If not, displays a regular message.
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err)
			msgs = append(msgs, "Something went wrong.")
		}
	}
	return msgs
}
