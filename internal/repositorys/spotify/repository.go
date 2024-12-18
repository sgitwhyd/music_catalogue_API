package spotify

import (
	"context"

	"github.com/sgitwhyd/music-catalogue/internal/models/spotify"
	"gorm.io/gorm"
)

type (
	SpotifySearchResponse struct {
		Tracks SpotifyTrack `json:"tracks"`
	}

	SpotifyTrack struct {
		Href 			string 								`json:"href"`
		Limit 		int 									`json:"limit"`
		Next 			*string									`json:"next"`
		Offset 		int 									`json:"offset"`
		Previous 	*string 							`json:"previous"`
		Total 		int 									`json:"total"`
		Items 		[]SpotifyTrackObject 	`json:"items"`
	}

	SpotifyTrackObject struct {
		Album 		SpotifyAlbumObject 			`json:"album"`
		Artists 	[]SpotifyArtisObject 		`json:"artists"`
		Explicit 	bool 										`json:"explicit"`
		Href 			string 									`json:"href"`
		ID 				string 									`json:"id"`
		Name 			string 									`json:"name"`
	}

	SpotifyArtisObject struct {
		Name 					string 								`json:"name"`
		Href 					string 								`json:"href"`
	}

	SpotifyAlbumObject struct {
		AlbumType 		string 								`json:"album_type"`
		TotalTracks 	int 								`json:"total_tracks"`
		Images 				[]SpotifyImagesObject `json:"images"`
		Name 					string 								`json:"name"`
	}

	SpotifyImagesObject struct {
		URL string `json:"url"`
	}
)


type spotifyRepository struct {
	db *gorm.DB
}

func NewSpotifyRepository (db *gorm.DB) *spotifyRepository {
	return &spotifyRepository{
		db: db,
	}
}

//go:generate mockgen -source=repository.go -destination=../../services/spotify/service_mock_test.go -package=spotify
type SpotifyOutbond interface {
	Search(ctx context.Context, query string, limit, offset int) (*SpotifySearchResponse, error)
}
type SpotifyRepository interface {
	Create(ctx context.Context, model spotify.TrackActivity) error
	Update(ctx context.Context, model spotify.TrackActivity) error
	Get(ctx context.Context, UserID uint, spotifyID string) (*spotify.TrackActivity, error)
	GetBulkSpotifyIDs(ctx context.Context, UserID uint, spotifyIDs []string) (map[string]spotify.TrackActivity, error)
}

func (r *spotifyRepository) Create(ctx context.Context, model spotify.TrackActivity) error {
	return r.db.Create(&model).Error
}

func (r *spotifyRepository) Update(ctx context.Context, model spotify.TrackActivity) error {
	return r.db.Save(&model).Error
}

func (r *spotifyRepository) Get(ctx context.Context, UserID uint, spotifyID string) (*spotify.TrackActivity, error) {
	activity := spotify.TrackActivity{}

	response := r.db.Where("user_id = ?", UserID).Where("spotify_id = ?", spotifyID).First(&activity)
	if response.Error != nil {
		return nil, response.Error
	}

	return &activity, nil
}

func (r *spotifyRepository) GetBulkSpotifyIDs(ctx context.Context, UserID uint, spotifyIDs []string) (map[string]spotify.TrackActivity, error) {
	activities := []spotify.TrackActivity{}

	response := r.db.Where("user_id = ?", UserID).Where("spotify_id IN ?", spotifyIDs).Find(&activities)
	if response.Error != nil {
		return nil, response.Error
	}

	result := make(map[string]spotify.TrackActivity, 0)
	for _,activity:=range activities {
		result[activity.SpotifyID] = activity
	}

	return result, nil
}
