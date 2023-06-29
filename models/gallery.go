package models

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
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

// galleryDir returns the directory where the images of the given directory are
// stored. If no directory is specified (i.e.: empty string), then the standard
// directory, defined as the constant stdImageDir, is used.
func (service *GalleryService) galleryDir(id int) string {
	imagesDir := service.ImagesDir
	if imagesDir == "" {
		imagesDir = stdImageDir
	}

	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}
