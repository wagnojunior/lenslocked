package models

import (
	"database/sql"
	"fmt"
)

// Gallery defines the gallery model according to the `gallery` SQL table
type Gallery struct {
	ID     int
	UserID int
	Title  string
}

// GalleryService defines the connection to the `gallery` DB
type GalleryService struct {
	DB *sql.DB
}

// Create creates a new gallery with the given title and associated with the
// given user
func (service *GalleryService) Create(title string, userID int) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userID,
	}

	row := service.DB.QueryRow(`
		INSERT INTO galleries (title, user_id)
		VALUES ($1, $2) RETURNING id`,
		title, userID)

	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}

	return &gallery, nil
}
