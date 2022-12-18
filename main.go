package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// homeHandler handles http requests to the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Sets the content type of the response header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Writes a html tag to the response writer w
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

// contactHandler handles the http requests to the contact page
func contactHandler(w http.ResponseWriter, r *http.Request) {
	// Sets the content type of the response header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Writes a html tag to the response writer w
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:wagnojunior@gmail.com\">wagnojunior@gmail.com</a>.")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	// Sets the content type of the response header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Writes a html tag to the respponse writer
	fmt.Fprintf(w, "<h1>FAQ Page</h1><p>Q: Question number one? <p>A: Answer number one.")
}

func main() {
	// Creates a new chi router
	r := chi.NewRouter()

	// Routes
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)

	// Starts the server
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
