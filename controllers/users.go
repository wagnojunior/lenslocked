package controllers

import (
	"net/http"

	"github.com/wagnojunior/lenslocked/views"
)

// Type Users holds a template struct that stores all the templates needed to render
// different pages
type Users struct {
	Templates struct {
		New views.Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil)
}
