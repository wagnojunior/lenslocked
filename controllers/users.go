package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/wagnojunior/lenslocked/context"
	"github.com/wagnojunior/lenslocked/models"
)

// Type Users holds a template struct that stores all the templates needed to render different pages
type Users struct {
	Templates struct {
		New            Template
		SignIn         Template
		SignOut        Template
		ForgotPassword Template
		CheckYourEmail Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
}

// New executes the template `New` that is stored in `u.Templates`
func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

// Create creates a new user when the sign up form is submited
func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	// email := r.FormValue("email")
	// password := r.FormValue("password")

	user, err := u.UserService.Create(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Creates a session after creating an user, since it is unecessary to ask a user to login immediately after they have signed up
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		// TODO: long term-> show a warning about not being able to sign in
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)

	fmt.Fprintf(w, "User created: %+v", user)
}

// SignIn executes the template `SignIn` that is stored in `u.Templates`
func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

// ProcessSignIn executes the template `SignIn` that is stored in `u.Templates`
func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	// Authenticate user
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Proper location to set cookies is after authentication and before writing to the response writer
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)

	fmt.Fprintf(w, "User authenticated: %+v", user)
}

// CurrentUser retrieves the user from the context
func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	// Sets data to be passed to the template
	var data struct {
		Email string
	}
	data.Email = user.Email

	u.Templates.SignOut.Execute(w, r, data)
}

// ProcessSignOut signs out a user, deletes the session from the DB, and
// redirect the user to the sign in page
func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	// Retrieve the session cookie and redirects to the signin page in case the
	// session is not set
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	// Delete the session from the DB
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

// ForgotPassword executes the template `ForgotPassword` stored in `u.Templates`
func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")

	// Create a new password reset token
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO: handle other cases in the future.
		// 1. user does not exists with a certain email address
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// Create the URL that is sent to the user. Since this application is not
	// in production mode, the URL is hard coded. However, in the futere it
	// will be associated with the user model
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetURL := "https://www.lenslocked.com/reset-pw?" + vals.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// Executes the `CheckYourEmail` template.
	// NOTE: do not render the reset token here! We need the user to confirm
	// access to the registered email to verify their email
	u.Templates.CheckYourEmail.Execute(w, r, data)

}

///////////////////////////////////////////////////////////////////////////////

// UserMiddleware defines a new type to handle the user middleware
type UserMiddleware struct {
	SessionService *models.SessionService
}

// SetUser looks up a token session from the cookie session, retrieves the user associated with this token, and sets this user to the current request
func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r) // proceed with the request assuming the user is not logged in
			return
		}

		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r) // proceed with the request assuming the user is not logged in
			return
		}

		// Gets the context from the request, overwrites the context with the user, and updates the request with the context
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// RequireUser checks if an user is signed in and redirects to the signin page if it ins't
func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
