package controllers

import (
	"net/http"

	"github.com/wagnojunior/lenslocked/models"
)

// Galleries holds the template struct that stores all the templates needed to
// render different pages. Also, it holds the necessary services
type Galleries struct {
	Templates struct {
		New Template
	}
	GalleryService *models.GalleryService
}

// New executes the template `New` that is stored in `g.Template`
func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}
