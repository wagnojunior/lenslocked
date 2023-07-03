package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
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

// Create handles the creation of a new gallery. Galleries are created
// unpublished by default
func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
		Status models.PublicationStatus
	}

	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")
	data.Status = models.Unpublished

	gallery, err := g.GalleryService.Create(data.Title, data.Status, data.UserID)
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

	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID     int
		Title  string
		Status models.PublicationStatus
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	data.Status = gallery.Status
	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}

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
	gallery, err := g.galleryByID(w, r, galleryMustBeVisible)
	if err != nil {
		return
	}

	// Creates a custom type Image to help construct the URL. This information
	// will be sent to the front-end, so it is a good idea to send only the
	// strictly necessary information. That is why the Image object is
	// constructed here, and not returned
	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID     int
		Title  string
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}

	g.Templates.Show.Execute(w, r, data)

}

// Publish handles the change of status of a gallery from unpublished to
// published
func (g Galleries) Publish(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	err = g.GalleryService.Publish(gallery)
	if err != nil {
		http.Error(w, "could not publish the gallery", http.StatusInternalServerError)
		return
	}

	// Redirect the user to the page them came from (refresh page)
	editPath := r.Header.Get("Referer")
	http.Redirect(w, r, editPath, http.StatusFound)
}

// Unpublish handles the change of status of a gallery from published to
// unpublished
func (g Galleries) Unpublish(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	err = g.GalleryService.Unpublish(gallery)
	if err != nil {
		http.Error(w, "could not unpublish the gallery", http.StatusInternalServerError)
		return
	}

	// Redirect the user to the page them came from (refresh page)
	editPath := r.Header.Get("Referer")
	http.Redirect(w, r, editPath, http.StatusFound)
}

// Delete deletes a gallery
func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	err = g.GalleryService.Delete(gallery.ID)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)

}

// Index looks up all of a user's galleries and sends this information to be
// rendered in a template
func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID     int
		Title  string
		Status models.PublicationStatus
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
			ID:     gallery.ID,
			Title:  gallery.Title,
			Status: gallery.Status,
		})
	}

	g.Templates.Index.Execute(w, r, data)

}

// Image handles HTTP requests to show an image.
func (g Galleries) Image(w http.ResponseWriter, r *http.Request) {
	filename := g.filename(w, r)
	galleryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid gallery ID", http.StatusNotFound)
		return
	}

	image, err := g.GalleryService.Image(galleryID, filename)
	if err != nil {
		if errors.Is(err, models.ErrImageNotFound) {
			http.Error(w, "image not found", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, image.Path)
}

// DeleteImage handles the HTTP request to delete an image
func (g Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	// Verifies that the user actually owns the gallery
	filename := g.filename(w, r)
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	err = g.GalleryService.DeleteImage(gallery.ID, filename)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// Redirects the user to the edit page
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)

}

// UploadImage handlers the HTTP POST request to upload an image
func (g Galleries) UploadImage(w http.ResponseWriter, r *http.Request) {
	// Verifies that the user actually owns the gallery
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	// Parses the request body as multipart/form-data. The argument is the
	// maximum memory allocated for this action. The parts of the request body
	// that do not fit in the allocated memory are stored in disk as temporary
	// files
	err = r.ParseMultipartForm(5 << 20) // 5MB
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// Accesses the file headers that were uploaded under the name `images` in
	// the HTML form. Since many files can be uploaded simulteneously, a slice
	// of file headers is returned
	fileHeaders := r.MultipartForm.File["images"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}
		defer file.Close()

		err = g.GalleryService.CreateImage(gallery.ID, fileHeader.Filename, file)
		if err != nil {
			var fileErr models.FileError
			if errors.As(err, &fileErr) {
				msg := fmt.Sprintf("%v has an invalid content type or extension. Only .png, .jpg, .jpeg, and .gif files can be uploaded.", fileHeader.Filename)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}

	// Redirects the user to the edit page
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

// /////////////////////////////////////////////////////////////////////////////
// HELPER FUNCTIONS
// /////////////////////////////////////////////////////////////////////////////

// filename returns the base of the provided filename. This is done to avoid
// malicious inputs from users
func (g Galleries) filename(w http.ResponseWriter, r *http.Request) string {
	filename := chi.URLParam(r, "filename")
	filename = filepath.Base(filename)

	return filename
}

// /////////////////////////////////////////////////////////////////////////////
// FUNCTIONAL OPTIONS
// /////////////////////////////////////////////////////////////////////////////

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

// galleryMustBeVisible checks if a user has acccess to the given gallery. If a
// user does not own a gallery and it is set to UNPUBLISHED, then access to the
// gallery is denied. Otherwise, access is granted
func galleryMustBeVisible(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())

	var thirdPartyGallery bool = (gallery.UserID != user.ID)
	var unpublished bool = (gallery.Status == models.Unpublished)
	if thirdPartyGallery && unpublished {
		http.Error(w, "gallery not found", http.StatusNotFound)
		return fmt.Errorf("gallery is not published and user does not have access to it")
	}

	return nil
}
