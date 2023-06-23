package controllers

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/wagnojunior/lenslocked/context"
	"github.com/wagnojunior/lenslocked/models"
)

// Galleries holds the template struct that stores all the templates needed to
// render different pages. Also, it holds the necessary services
type Galleries struct {
	Templates struct {
		Show  Template
		New   Template
		Edit  Template
		Index Template
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

// Create handles the creation of a new gallery
func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
	}

	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")

	gallery, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}

	// If there is no errers the user is redirected to the `edit gallery page`
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

// Edit renders the `Edit` template where a gallery's title can be eddited. The
// gallery's ID is retrieved from the URL parameters, and a authorization check
// is performed to assess whether the requesting user owns the gallery.
func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	var data struct {
		ID    int
		Title string
	}

	data.ID = gallery.ID
	data.Title = gallery.Title

	// Renders the `Edit` page with the passed data
	g.Templates.Edit.Execute(w, r, data)

}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	// Retreives the new gallery name from the form
	gallery.Title = r.FormValue("title")
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "could not update the gallery", http.StatusInternalServerError)
		return
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)

}

// Show shows the images in a gallery
func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var data struct {
		ID     int
		Title  string
		Images []string
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	for i := 0; i < 20; i++ {
		w, h := rand.Intn(500)+200, rand.Intn(500)+200
		catImageURL := fmt.Sprintf("https://placekitten.com/%d/%d", w, h)
		data.Images = append(data.Images, catImageURL)
	}

	g.Templates.Show.Execute(w, r, data)
}

// Index looks up all of a user's galleries and sends this information to be
// rendered in a template
func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID    int
		Title string
	}

	var data struct {
		Galleries []Gallery
	}

	user := context.User(r.Context())
	galleries, err := g.GalleryService.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// Needs to translate the gallery which is stored in the DB to the
	// Gallery type created in this handler
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}

	g.Templates.Index.Execute(w, r, data)

}

// galleryOpt defines a functional option. Functions that have this signature
// are of type galleryOpt
type galleryOpt func(http.ResponseWriter, *http.Request, *models.Gallery) error

// galleryByID is a helper method that gets a gallery by ID and returns it. It
// receives a functional options, which are set the the caller of galleryByID
func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid ID", http.StatusNotFound)
		return nil, err
	}

	// Gets the gallery by the provided ID
	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrInvalidGallery) {
			http.Error(w, "gallery not found", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "something went wrong", http.StatusNotFound)
		return nil, err
	}

	// Loops through all functional options and returns an error, if any. This
	// erros is subsequently handled by the underlying function
	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}

	return gallery, nil
}

// userMustOwnGallery is a functional option which determines that a user must
// own a gallery
func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	// Checks whether the retrieved gallery belongs to the user that requested
	// it
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "gallery not found", http.StatusNotFound)
		return fmt.Errorf("user does not have access to this gallery")
	}

	return nil
}
