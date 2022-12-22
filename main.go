package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wagnojunior/lenslocked/controllers"
	"github.com/wagnojunior/lenslocked/templates"
	"github.com/wagnojunior/lenslocked/views"
)

func main() {
	// Creates a new chi router
	r := chi.NewRouter()

	// Parses the home templates before the server starts
	// views.Parse returns a Template and an error. This fits the scope of views.Must
	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	// Parses the contact templates before the server starts
	// views.Parse returns a Template and an error. This fits the scope of views.Must
	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	// Parses the faq templates before the server starts
	// views.Parse returns a Template and an error. This fits the scope of views.Must
	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml"))
	r.Get("/faq", controllers.StaticHandler(tpl))

	// Starts the server
	// views.Parse returns a Template and an error. This fits the scope of views.Must
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
