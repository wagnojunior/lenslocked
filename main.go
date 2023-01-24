package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/wagnojunior/lenslocked/controllers"
	"github.com/wagnojunior/lenslocked/migrations"
	"github.com/wagnojunior/lenslocked/models"
	"github.com/wagnojunior/lenslocked/templates"
	"github.com/wagnojunior/lenslocked/views"
)

func main() {
	// Creates a new chi router
	r := chi.NewRouter()

	// Parses the home templates before the server starts. `views.ParseFS` returns a Template and an error. This fits the scope of `views.Must`
	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	// Parses the contact templates before the server starts. `views.ParseFS` returns a Template and an error. This fits the scope of views.Must
	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	// Parses the faq templates before the server starts. `views.ParseFS` returns a Template and an error. This fits the scope of views.Must
	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	// Creates a connection to the DB
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Runs the migrations in the current directory of the filesystem (thus the ".")
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Creates a UserService
	userService := models.UserService{
		DB: db,
	}

	// Creates the SessionService
	sessionService := models.SessionService{
		DB: db,
	}

	// Initializes the controller for the users `usersC`. This controller receives the template that is parsed by `views.ParseFS`. `usersC.New` is of type HandlerFunc, therefore it can be passed to `r.Get`. In here, `usersC.New` is passed as a type function, therefore no need to pass in the arguments.
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS, "signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignOut = views.Must(views.ParseFS(
		templates.FS, "me.gohtml", "tailwind.gohtml"))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/users/me", usersC.CurrentUser)

	// Set middleware
	csrfKey := "5YGEgDV0VAVTlV8wxfXdlCJSam82rvj1"
	csrfMW := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false), // TODO: change false -> true for deployment
	)

	// Starts the server
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", csrfMW(r))
}
