package controllers

import (
	"net/http"
)

// Type Users holds a template struct that stores all the templates needed to render
// different pages
type Users struct {
	Templates struct {
		New Template
	}
}

// New executes the template `New` that is stored in `u.Templates`
func (u Users) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil)
}
