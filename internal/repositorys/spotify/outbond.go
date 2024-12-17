package spotify

import (
	"time"

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