package controllers

import (
	"fmt"
	"net/http"

	"github.com/wagnojunior/lenslocked/models"
)

// Type Users holds a template struct that stores all the templates needed to render different pages
type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
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
	cookie := http.Cookie{
		Name:     "session",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
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
	cookie := http.Cookie{
		Name:     "session",     // name of the cookie
		Value:    session.Token, // value being stored
		Path:     "/",           // which paths on the server have access to this cookie
		HttpOnly: true,          // cookie not accessible via javascript
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/users/me", http.StatusFound)

	fmt.Fprintf(w, "User authenticated: %+v", user)
}

// CurrentUser retrieves the session cookie from the http request and uses it to authenticate the user. In case of error, redirect to the signin page
func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	user, err := u.SessionService.User(tokenCookie.Value)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}
