package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/wagnojunior/lenslocked/views"
)

// executeTemplate parses and executes a gohtml template
func executeTemplate(w http.ResponseWriter, filepath string) {
	// Parses the template located at the directory filepath
	t, err := views.Parse(filepath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}

	// Executes the parsed template
	t.Execute(w, nil)
}

// homeHandler handles http requests to the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Joins path to template to be parsed and executed
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath)
}

// contactHandler handles the http requests to the contact page
func contactHandler(w http.ResponseWriter, r *http.Request) {
	// Joins path to template to be parsed and executed
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath)
}

// faqHandler handles the http request to the faq page
func faqHandler(w http.ResponseWriter, r *http.Request) {
	// Joins path to template to be parsed and executed
	tplPath := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, tplPath)
}

// userHandler handles the http request to the user page
func userHandler(w http.ResponseWriter, r *http.Request) {
	// Sets the content type of the response header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Fetch the URL parameter userID
	userID := chi.URLParam(r, "userID")

	// Writes a html tag to the respponse writer
	fmt.Fprintf(w, "<h1> Welcome back, user %s</h1>", userID)
}

func main() {
	// Creates a new chi router
	r := chi.NewRouter()

	// Routes
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.Get("/user/{userID}", userHandler)

	// Starts the server
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
