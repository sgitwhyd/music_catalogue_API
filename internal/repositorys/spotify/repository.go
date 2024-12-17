package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rs/zerolog/log"
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


type SpotifyRepository interface {
	Search(ctx context.Context, query string, limit, offset int) (*SpotifySearchResponse, error)
}

func (o *outbond) Search(ctx context.Context, query string, limit, offset int) (*SpotifySearchResponse, error) {
	// set url params
	params := url.Values{}
	params.Set("q", query)
	params.Set("type", "track")
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))

	BASE_URL := "https://api.spotify.com/v1/search"
	SEARCH_ENDPOINT := fmt.Sprintf(`%s?%s`, BASE_URL, params.Encode())

	// get token GetTokenDetails
	accessToken, tokenType, err := o.GetTokenDetails()
	if err != nil {
		return nil, err
	}

	BEARER_TOKEN := fmt.Sprintf("%s %s", tokenType, accessToken)

	req, err := http.NewRequest(http.MethodGet, SEARCH_ENDPOINT, nil)
	if err != nil {
		log.Error().Err(err).Msg("error create request spotify search")
		return nil, err
	}

	req.Header.Set("Authorization", BEARER_TOKEN)
	resp, err := o.client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("error execute search spotify")
		return nil, err
	}

	defer resp.Body.Close()

	var response SpotifySearchResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Error().Err(err).Msg("error decoded spotify search response")
		return nil, err
	}

	return &response, nil

}