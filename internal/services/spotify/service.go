package spotify

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/music-catalogue/internal/models/spotify"
	spotifyRepo "github.com/sgitwhyd/music-catalogue/internal/repositorys/spotify"
)

//go:generate mockgen -source=service.go -destination=service_mock_test.go -package=spotify
type SpotifyRepository interface {
	spotifyRepo.SpotifyRepository
}

//go:generate mockgen -source=service.go -destination=../../handlers/spotify/handler_mock_test.go -package=spotify
type SpotifyService interface {
	Search(ctx context.Context, query string, pageSize, pageIndex int) (*spotify.SearchResponse, error)
}

type spotifyService struct {
	spotifyOutbond spotifyRepo.SpotifyRepository
}

func NewSpotifyServie(spotifyOutbond spotifyRepo.SpotifyRepository) *spotifyService {
	return &spotifyService{
		spotifyOutbond: spotifyOutbond,
	}
}

func (s *spotifyService) Search(ctx context.Context, query string, pageSize, pageIndex int) (*spotify.SearchResponse, error) {
	limit := pageSize
	offset := (pageIndex - 1) * pageSize


	trackDetails, err := s.spotifyOutbond.Search(ctx, query, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("error search track spotify")
		return nil, err
	}

	return modelToResponse(trackDetails), nil
}

func modelToResponse(data *spotifyRepo.SpotifySearchResponse) *spotify.SearchResponse {
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
		}
	}

	return &spotify.SearchResponse{
		Limit: data.Tracks.Limit,
		Offset: data.Tracks.Offset,
		Total: data.Tracks.Total,
		Items: items,
	}
}