package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/music-catalogue/internal/configs"
	"github.com/sgitwhyd/music-catalogue/pkg/httpclient"
)

type outbond struct {
	cfg *configs.Config
	client httpclient.HTTPClient
	AccessToken string
	TokenType string
	ExpiredAt time.Time
}



func NewSpotifyOutbond(cfg *configs.Config, client httpclient.HTTPClient) *outbond {
		return &outbond{
			cfg: cfg,
			client: client,
		}
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
