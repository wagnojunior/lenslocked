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
	tpl := views.Must(views.ParseFS(
		templates.FS,
		"home.gohtml",
		"tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	// Parses the contact templates before the server starts
	// views.Parse returns a Template and an error. This fits the scope of views.Must
	tpl = views.Must(views.ParseFS(templates.FS,
		"contact.gohtml",
		"tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	// Parses the faq templates before the server starts
	// views.Parse returns a Template and an error. This fits the scope of views.Must
	tpl = views.Must(views.ParseFS(
		templates.FS,
		"faq.gohtml",
		"tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	// Initializes the controller for the users `usersC`. This controller receives the
	// template that is parsed by `views.ParseFS`. `usersC.New` is a of type HandlerFunc,
	// therefore it can be passed to `r.Get`. In here, `usersC.New` is passed as a type
	// function, therefore no need to pass in the arguments.
	usersC := controllers.Users{}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml",
		"tailwind.gohtml"))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)

	// Starts the server
	// views.Parse returns a Template and an error. This fits the scope of views.Must
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
