package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// homeHandler handles http requests to the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	bio := `<script>alert("HaHa, you have been h4x0r3d!");</script>`

	// Sets the content type of the response header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Writes a html tag to the response writer w
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1><p>Bio:"+bio+"</p>")
}

// contactHandler handles the http requests to the contact page
func contactHandler(w http.ResponseWriter, r *http.Request) {
	// Sets the content type of the response header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Writes a html tag to the response writer w
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:wagnojunior@gmail.com\">wagnojunior@gmail.com</a>.")
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
