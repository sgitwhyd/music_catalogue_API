package spotify

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/music-catalogue/internal/models/spotify"
	spotifyRepo "github.com/sgitwhyd/music-catalogue/internal/repositorys/spotify"
	"gorm.io/gorm"
)

//go:generate mockgen -source=service.go -destination=../../handlers/spotify/handler_mock_test.go -package=spotify
type SpotifyService interface {
	Search(ctx context.Context, query string, pageSize, pageIndex int,  userID uint) (*spotify.SearchResponse, error)
	UpSertActivity(ctx context.Context, userID uint, request spotify.TrackActivityRequest) error
}

type spotifyService struct {
	spotifyOutbond spotifyRepo.SpotifyOutbond
	spotifyRepo spotifyRepo.SpotifyRepository
}

func NewSpotifyServie(spotifyOutbond spotifyRepo.SpotifyOutbond, spotifyRepo spotifyRepo.SpotifyRepository) *spotifyService {
	return &spotifyService{
		spotifyOutbond: spotifyOutbond,
		spotifyRepo: spotifyRepo,
	}
}

func (s *spotifyService) Search(ctx context.Context, query string, pageSize, pageIndex int, userID uint) (*spotify.SearchResponse, error) {
	limit := pageSize
	offset := (pageIndex - 1) * pageSize


	trackDetails, err := s.spotifyOutbond.Search(ctx, query, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("error search track spotify")
		return nil, err
	}

	trackIDs := make([]string, len(trackDetails.Tracks.Items))
	for idx, track := range trackDetails.Tracks.Items {
		trackIDs[idx] = track.ID
	}

	trackActivities, err := s.spotifyRepo.GetBulkSpotifyIDs(ctx, userID, trackIDs)
	if err != nil {
		log.Error().Err(err).Msg("error get track activities from db")
		return nil, err
	}


	return modelToResponse(trackDetails, trackActivities), nil
}

func modelToResponse(data *spotifyRepo.SpotifySearchResponse, mapTrackActivities map[string]spotify.TrackActivity) *spotify.SearchResponse {
	if data == nil {
		return nil
	}

	items := make([]spotify.SpotifyTrackObjectResponse, len(data.Tracks.Items))

	for i, item := range data.Tracks.Items{

		artisName := make([]string, len(item.Artists))
		for idx, artist := range item.Artists {
			artisName[idx] = artist.Name
		}

		imageUrl := make([]string, len(item.Album.Images))
		for idx, image := range item.Album.Images {
			imageUrl[idx] = image.URL
		}

		items[i] = spotify.SpotifyTrackObjectResponse{
		// album related fields
		AlbumType       : item.Album.AlbumType,
		AlbumTotalTracks : item.Album.TotalTracks,
		AlbumImagesURL    : imageUrl,
		AlbumName        : item.Album.Name,

		// artists related field
		ArtistsName : artisName,

		// track related field
		Explicit : item.Explicit,
		Href     : item.Href,
		ID      : item.ID,
		Name     : item.Name,
		IsLiked: mapTrackActivities[item.ID].IsLiked,
		}
	}

	return &spotify.SearchResponse{
		Limit: data.Tracks.Limit,
		Offset: data.Tracks.Offset,
		Total: data.Tracks.Total,
		Items: items,
	}
}

func (s *spotifyService) UpSertActivity(ctx context.Context, userID uint, request spotify.TrackActivityRequest) error {

	foundedActivity, err := s.spotifyRepo.Get(ctx, userID, request.SpotifyID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Msg("service: error get record from db")
		return err
	}

	if err == gorm.ErrRecordNotFound || foundedActivity == nil {
		err = s.spotifyRepo.Create(ctx, spotify.TrackActivity{
			UserID: userID,
			SpotifyID: request.SpotifyID,
			IsLiked: request.IsLiked,
		})

		if err != nil {
			log.Error().Err(err).Msg("service: error create record from db")
			return err
		}

		return nil
	} 

	foundedActivity.IsLiked = request.IsLiked
	err = s.spotifyRepo.Update(ctx, *foundedActivity)
	if err != nil {
		log.Error().Err(err).Msg("service: error update record from db")
		return err
	}
	
	return nil
}