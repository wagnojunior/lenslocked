// Added comment to rebuild
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

	// PSQL configuration
	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	var serverNotConfig bool = (cfg.PSQL.Host == "" && cfg.PSQL.Port == "")
	if serverNotConfig {
		return cfg, fmt.Errorf("no PSQL config provided")
	}

	// SMTP configuration
	cfg.SMTP.Host = os.Getenv("SMPT_HOST")
	portStr := os.Getenv("SMPT_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMPT_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMPT_PASSWORD")

	// SCRF configuration
	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = (os.Getenv("CSRF_SECURE") == "true")

	// SERVER configuration
	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil

}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	err = run(cfg)
	if err != nil {
		panic(err)
	}
}

// run runs the application. This is apart from main to make it easier to test
// the main function
func run(cfg config) error {

	// Creates a connection to the DB
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Runs the migrations in the current directory of the filesystem
	// (thus the ".")
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		return err
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
		DB:         db,
		ImagesDir:  "",                // Use default value if not set
		ImagesExt:  make([]string, 0), // Use default value if not set
		ImagesCont: make([]string, 0), // Use default value if not set
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
	galleriesC.Templates.Edit = views.Must(views.ParseFS(
		templates.FS, "galleries/edit.gohtml", "tailwind.gohtml"))
	galleriesC.Templates.Index = views.Must(views.ParseFS(
		templates.FS, "galleries/index.gohtml", "tailwind.gohtml"))
	galleriesC.Templates.Show = views.Must(views.ParseFS(
		templates.FS, "galleries/show.gohtml", "tailwind.gohtml"))

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
		r.Get("/{id}", galleriesC.Show)
		r.Get("/{id}/images/{filename}", galleriesC.Image)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/new", galleriesC.New)
			r.Get("/", galleriesC.Index)
			r.Post("/", galleriesC.Create)
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Post("/{id}", galleriesC.Update)
			r.Post("/{id}/delete", galleriesC.Delete)
			r.Post("/{id}/publish", galleriesC.Publish)
			r.Post("/{id}/unpublish", galleriesC.Unpublish)
			r.Post("/{id}/images/{filename}/delete", galleriesC.DeleteImage)
			r.Post("/{id}/images", galleriesC.UploadImage)
		})

	})
	// Serve static files from the folder `assets`
	assetsHandler := http.FileServer(http.Dir("assets"))
	r.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	// Starts the server
	fmt.Printf("Starting the server on %s...", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		return err
	}

	return nil
}
