package controllers

import (
	"fmt"
	"net/http"

	"github.com/wagnojunior/lenslocked/models"
)

// Type Users holds a template struct that stores all the templates needed to render different pages
type Users struct {
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

// New executes the template `New` that is stored in `u.Templates`
func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

// Create creates a new user when the sign up form is submited
func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	// r.FormValue("KEY_NAME") where KEY_NAME is defined in the form
	fmt.Fprint(w, "Email: ", r.FormValue("email"))
	fmt.Fprint(w, "Password: ", r.FormValue("password"))

}
