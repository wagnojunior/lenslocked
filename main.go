package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	// Sets the content type of the response header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Parses the template home.gohtml located in the folder <templates>
	// If there is an error parsing, it will be handled here (invalid function in the template)
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("parsing templates: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}

	// Executes the template tpl without any data (nil)
	// If there is an error rendering, it will be handled here (invalid field in the template)
	// This approach writes to the response writer until an error is detected (if any). If an error
	// is detected half-way through the execution, then the webpage will be half rendered
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("executing templates: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
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
	// Sets the content type of the response header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Writes a html tag to the respponse writer
	fmt.Fprintf(w, "<h1>FAQ Page</h1><p>Q: Question number one? <p>A: Answer number one.")
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
