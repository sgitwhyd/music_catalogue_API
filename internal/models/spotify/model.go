package spotify

import "gorm.io/gorm"

type (
	SearchResponse struct {
		Items  []SpotifyTrackObjectResponse `json:"items"`
		Limit  int                          `json:"limit"`
		Offset int                          `json:"offset"`
		Total  int                          `json:"total"`
	}

	SpotifyTrackObjectResponse struct {
		// album related fields
		AlbumType        string   `json:"album_type"`
		AlbumTotalTracks int      `json:"album_total_tracks"`
		AlbumImagesURL   []string `json:"album_image_url"`
		AlbumName        string   `json:"album_name"`

		// artis related fields
		ArtistsName []string `json:"artists_name"`

		// track related field
		Explicit bool   `json:"explicit"`
		Href     string `json:"href"`
		ID       string `json:"id"`
		Name     string `json:"name"`
		IsLiked  *bool	`json:"is_liked"`
	}
)

// track activities

type (
	TrackActivity struct {
		gorm.Model
		UserID 		uint `gorm:"not null"`
		SpotifyID string `gorm:"not null"`
		IsLiked 	*bool
	}

	TrackActivityRequest struct {
		SpotifyID string `json:"spotify_id" binding:"required"`
		IsLiked *bool `json:"is_liked"`
	}
)