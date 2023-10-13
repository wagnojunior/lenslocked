package controllers

import (
	"html/template"
	"net/http"
)

// StaticHandler executes a template and returns a HandlerFunc.
// It is actually a closure, so it is possible to access variables outside of
// its scope (such as tpl).
func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

// FAQ handles the templating of the FAQ page
func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question template.HTML
		Answer   template.HTML // This type indicates that it is safe to render the answers as HTML
	}{
		{
			Question: "What is the purpose of this website?",
			Answer:   `The purpose of this website is to learn web development with Go. It is the final project of the course <span class="font-semibold">Web Development with Go v2</span> by Jon Calhoun.`,
		},
		{
			Question: "What did you lean during the development of this application?",
			Answer:   "I learned about the basics of web development (HTTP methods, handlers, routers, cookies, sessions, etc.), MVC design pattern, database (PostgreSQL), schema migrations, user authentication, Go templates, sending emails from the application, error handling, and deployment with Docker.",
		},
		{
			Question: "What skills did you develop during the development of this application?",
			Answer:   `I developed proficiency in <span class="font-semibold">Vim motions</span>, which significantly increased my typing speed. Moreover, I improved my undersdanting of UI/UX and the importance of responsive design (using Tailwind).`,
		},
		{
			Question: "Did you apply what you leaned on other projects?",
			Answer:   `Yes! I am already developing a commercial website for a consulting company in South Korea.`,
		},
		{
			Question: `Do you need previous knowledge of Go to take the course <span class="font-semibold">Web Development with Go v2</span> by Jon Calhoun?`,
			Answer:   `Yes! A basic understanding of Go is required. I recommend the Udemy course <span class="font-semibold">Go: The Complete Developer's Guide</span> by Stepehn Grider.`,
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, questions)
	}
}
