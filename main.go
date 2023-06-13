package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/wagnojunior/lenslocked/controllers"
	"github.com/wagnojunior/lenslocked/migrations"
	"github.com/wagnojunior/lenslocked/models"
	"github.com/wagnojunior/lenslocked/templates"
	"github.com/wagnojunior/lenslocked/views"
)

// config defines all environment variables that this application requires to
// run
type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

// loadEnvConfig loads the environment variables and sets the config for this
// application
func loadEnvConfig() (config, error) {
	var cfg config

	// Loads the env variables
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	// TODO: PSQL
	cfg.PSQL = models.DefaultPostgresConfig()

	// TODO: SMTP
	cfg.SMTP.Host = os.Getenv("SMPT_HOST")
	portStr := os.Getenv("SMPT_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMPT_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMPT_PASSWORD")

	// TODO: CSRF
	cfg.CSRF.Key = "5YGEgDV0VAVTlV8wxfXdlCJSam82rvj1"
	cfg.CSRF.Secure = false

	// TODO: SERVER
	cfg.Server.Address = ":3000"

	return cfg, nil

}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Creates a connection to the DB
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Runs the migrations in the current directory of the filesystem
	// (thus the ".")
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Defines the services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	galleryService := &models.GalleryService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)

	// Creates an instance of the UserMiddleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	// Sets middleware
	csrfMW := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"), // Sets all cookies path to `/`
	)

	// Initializes the controller for the users `usersC`. This controller
	// receives the template that is parsed by `views.ParseFS`. `usersC.New` is
	// of type HandlerFunc, therefore it can be passed to `r.Get`. In here
	// `usersC.New` is passed as a type function, therefore no need to pass in
	// the arguments.
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS, "sign-up.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS, "sign-in.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignOut = views.Must(views.ParseFS(
		templates.FS, "me.gohtml", "tailwind.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS, "forgot-pw.gohtml", "tailwind.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(
		templates.FS, "check-your-email.gohtml", "tailwind.gohtml"))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(
		templates.FS, "reset-pw.gohtml", "tailwind.gohtml"))

	// Initializes the controller for the galleries `galleriesC`
	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}

	galleriesC.Templates.New = views.Must(views.ParseFS(
		templates.FS, "galleries/new.gohtml", "tailwind.gohtml"))

	// Creates a new chi router and applies the different middlewares
	r := chi.NewRouter()
	r.Use(csrfMW)
	r.Use(umw.SetUser)

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))
	r.Get("/faq", controllers.FAQ(views.Must(
		views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.Route("/galleries", func(r chi.Router) {
		// r.Group groups all paths to the same middleware
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/new", galleriesC.New)
		})

	})

	// Starts the server
	fmt.Printf("Starting the server on %s...", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
