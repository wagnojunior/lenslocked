package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// PublicationStatus defines a new type that wraps the native string type. It
// represents the publication status of a gallery (i.e.: published,
// unpublished). The default publication status of a gallery is unpublished
type PublicationStatus string

const (
	// standard directory where the images are stored
	stdImageDir string = "images"
)

// Defines the two publication status
var (
	Published   PublicationStatus = "published"
	Unpublished PublicationStatus = "unpublished"
)

// Gallery defines the gallery model according to the `gallery` SQL table
type Gallery struct {
	ID     int
	UserID int
	Title  string
	Status PublicationStatus
}

// Image defines a new type to represent an Image
type Image struct {
	GalleryID int
	Path      string
	Filename  string
}

// GalleryService defines the connection to the `gallery` DB
type GalleryService struct {
	DB *sql.DB
	// ImagesDir is used to tell the GalleryService where to store and locate
	// images. If not set, the GalleryService will default to using the "images"
	// directory
	ImagesDir string
}

// Create creates a new gallery with the given title, publication status and
// associated with the given user.
func (service *GalleryService) Create(title string, status PublicationStatus, userID int) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userID,
		Status: status,
	}

	row := service.DB.QueryRow(`
		INSERT INTO galleries (title, publication_status, user_id)
		VALUES ($1, $2, $3) RETURNING id`,
		title, status, userID)

	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}

	return &gallery, nil
}

// ByID query and returns a gallery by the given ID
func (service *GalleryService) ByID(id int) (*Gallery, error) {
	gallery := Gallery{
		ID: id,
	}

	row := service.DB.QueryRow(`
		SELECT title, publication_status, user_id
		FROM galleries
		WHERE id = $1`,
		id)

	err := row.Scan(&gallery.Title, &gallery.Status, &gallery.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidGallery
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}

	return &gallery, nil
}

// ByUserID query and returns all galleries associated with a user ID
func (service *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	rows, err := service.DB.Query(`
		SELECT id, title, publication_status
		FROM galleries
		WHERE user_id = $1`,
		userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidGallery
		}
		return nil, fmt.Errorf("query galleries by user id: %w", err)
	}

	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}

		err = rows.Scan(&gallery.ID, &gallery.Title, &gallery.Status)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user id: %w", err)
		}

		galleries = append(galleries, gallery)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user id: %w", err)
	}

	return galleries, nil
}

// Update updates the provided gallery
func (service *GalleryService) Update(gallery *Gallery) error {
	_, err := service.DB.Exec(`
		UPDATE galleries
		SET title = $2
		WHERE id = $1`,
		gallery.ID, gallery.Title)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}

	return nil
}

// Publish changes the publication status of a gallery from unpublished to
// publish
func (service *GalleryService) Publish(gallery *Gallery) error {
	newStatus := "published"

	_, err := service.DB.Exec(`
		UPDATE galleries
		SET publication_status = $2
		WHERE id = $1`,
		gallery.ID, newStatus)
	if err != nil {
		return fmt.Errorf("publish gallery: %w", err)
	}

	return nil
}

// Unpublish changes the publication status of a gallery from published to
// unpublish
func (service *GalleryService) Unpublish(gallery *Gallery) error {
	newStatus := "unpublished"

	_, err := service.DB.Exec(`
		UPDATE galleries
		SET publication_status = $2
		WHERE id = $1`,
		gallery.ID, newStatus)
	if err != nil {
		return fmt.Errorf("unpublish gallery: %w", err)
	}

	return nil
}

// Delete deletes a gallery by ID
func (service *GalleryService) Delete(id int) error {
	_, err := service.DB.Exec(`
		DELETE FROM galleries
		WHERE id = $1`,
		id)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}

	return nil
}

// /////////////////////////////////////////////////////////////////////////////
// IMAGES
// /////////////////////////////////////////////////////////////////////////////

// extensions return a list of image extensions supported by the server
func (service *GalleryService) extensions() []string {
	return []string{
		".png", ".jpg", ".jpeg", ".gif",
	}
}

// contentTypes returns a list of image content types supported by the server
func (service *GalleryService) contentTypes() []string {
	return []string{
		"image/png", "image/jpg", "image/jpeg", "image/gif",
	}
}

// Images returns a slice of Image in the given gallery.
func (service *GalleryService) Images(galleryID int) ([]Image, error) {
	// Firstly, the directory of the given gallery is retrieved. Secondly, a
	// glob pattern is constructed so that all files inside the directory are
	// returned. Thirdly, the extension of each file is checked to confirm if
	// it is amoung the supported extensions.

	galleryDir := service.galleryDir(galleryID)
	globPattern := filepath.Join(galleryDir, "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving images from gallery %d: %w", galleryID, err)
	}

	var images []Image
	supportedExt := service.extensions()
	for _, file := range allFiles {
		fileIsImage := hasExtension(file, supportedExt)
		if fileIsImage {
			images = append(images, Image{
				GalleryID: galleryID,
				Path:      file,
				Filename:  filepath.Base(file),
			})
		}
	}

	return images, nil
}

// Image returns the image defined by the given filename and given gallery. An
// error is returned in case the image does not exist
func (service *GalleryService) Image(galleryID int, filename string) (Image, error) {
	galleryDir := service.galleryDir(galleryID)
	imagePath := filepath.Join(galleryDir, filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrImageNotFound
		}

		return Image{}, fmt.Errorf("querying for image: %w", err)
	}

	image := Image{
		Filename:  filename,
		GalleryID: galleryID,
		Path:      imagePath,
	}
	return image, nil
}

// CreateImage creates an image from the provided contents and stores it in the
// respective gallery directory
func (service *GalleryService) CreateImage(galleryID int, filename string, contents io.ReadSeeker) error {
	// Checks whether the content-type and the extension of the file are
	// supported
	err := checkContentType(contents, service.contentTypes())
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}
	err = checkExtension(filename, service.extensions())
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}

	// Checks if the directory to which the image will be saved exists. In case
	// it does not, the directory is created
	galleryDir := service.galleryDir(galleryID)
	err = os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("creating gallery-%d image directory: %w", galleryID, err)
	}

	// Creates the image file
	imagePath := filepath.Join(galleryDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, contents)
	if err != nil {
		return fmt.Errorf("copying contents to image: %w", err)
	}

	return nil
}

// DeleteImage deletes the image defined by the given gallery ID and filename.
// It returnns nil if the image is successfully deleted, or an error otherwise
func (service *GalleryService) DeleteImage(galleryID int, filename string) error {
	image, err := service.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}

	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}

	return nil
}

// hasExtension returns true if the given file has one of the provided
// extensions
func hasExtension(file string, extension []string) bool {
	for _, ext := range extension {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)

		fileExt := filepath.Ext(file)
		sameExt := (fileExt == ext)
		if sameExt {
			return true
		}
	}

	return false
}

// galleryDir returns the directory where the images of the given gallery are
// stored. If no directory is specified (i.e.: empty string), then the standard
// directory, defined as the constant stdImageDir, is used.
func (service *GalleryService) galleryDir(id int) string {
	imagesDir := service.ImagesDir
	if imagesDir == "" {
		imagesDir = stdImageDir
	}

	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}
