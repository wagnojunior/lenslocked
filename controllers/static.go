package controllers

import (
	"net/http"

	"github.com/wagnojunior/lenslocked/views"
)

// StaticHandler executes a template and returns a HandlerFunc.
// It is actually a closure, so it is possible to access variables outside of
// its scope (such as tpl).
func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}
