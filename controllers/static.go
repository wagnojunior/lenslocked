package controllers

import (
	"html/template"
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

// FAQ handles the templating of the FAQ page
func FAQ(tpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML // This type indicates that it is safe to render the answers as HTML
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes! We offer a free trial for 30 days on any paid plans.",
		},
		{
			Question: "What are your support hours?",
			Answer:   "We have support staff answering emails 24/7, though response times may be a bit slower on weekends.",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us at <a href="mailto:support@panpancorp.com">support@panpancorp.com</a>`,
		},
		{
			Question: "Where is your office located?",
			Answer:   "Our entire time is remote.",
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
