package spotify

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/music-catalogue/internal/configs"
	"github.com/sgitwhyd/music-catalogue/internal/models/spotify"
	"github.com/sgitwhyd/music-catalogue/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_handler_Search(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSvc := NewMockSpotifyService(mockCtrl)

	tests := []struct {
		name             string
		mockFn           func()
		expectedCode     int
		expectedResponse spotify.SearchResponse
		wantErr          bool
	}{
		// TODO: Add test cases.
		{
			name:         "success",
			expectedCode: 200,
			wantErr:      false,
			expectedResponse: spotify.SearchResponse{
				Limit:  10,
				Offset: 0,
				Items: []spotify.SpotifyTrackObjectResponse{
					{
						AlbumType:        "album",
						AlbumTotalTracks: 22,
						AlbumImagesURL:   []string{"https://i.scdn.co/image/ab67616d0000b273e8b066f70c206551210d902b"},
						AlbumName:        "Bohemian Rhapsody (The Original Soundtrack)",
						ArtistsName:      []string{"Queen"},
						Explicit:         false,
						ID:               "3z8h0TU7ReDPLIbEnYhWZb",
						Name:             "Bohemian Rhapsody",
					},
				},
				Total: 905,
			},
			mockFn: func() {
				mockSvc.EXPECT().Search(gomock.Any(), "bohemian rhapsody", 10, 1, uint(1)).Return(&spotify.SearchResponse{
					Limit:  10,
					Offset: 0,
					Items: []spotify.SpotifyTrackObjectResponse{
						{
							AlbumType:        "album",
							AlbumTotalTracks: 22,
							AlbumImagesURL:   []string{"https://i.scdn.co/image/ab67616d0000b273e8b066f70c206551210d902b"},
							AlbumName:        "Bohemian Rhapsody (The Original Soundtrack)",
							ArtistsName:      []string{"Queen"},
							Explicit:         false,
							ID:               "3z8h0TU7ReDPLIbEnYhWZb",
							Name:             "Bohemian Rhapsody",
						},
					},
					Total: 905,
				}, nil)
			},
		},
		{
			name:             "failed",
			expectedCode:     400,
			wantErr:          true,
			expectedResponse: spotify.SearchResponse{},
			mockFn: func() {
				mockSvc.EXPECT().Search(gomock.Any(), "bohemian rhapsody", 10, 1, uint(1)).Return(nil, assert.AnError).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			config, err := configs.Init("../../configs", "env", "test.env")
			assert.NoError(t, err)
			
			gin.SetMode(gin.ReleaseMode)

			r := gin.New()
			route := r.Group("/api/v1")

			h := &handler{
				route: route,
				service:     mockSvc,
			}
			h.RegisterRoute()

			endpoint := "/api/v1/spotify/search?query=bohemian+rhapsody&pageIndex=1&pageSize=10"
			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			assert.NoError(t, err)

			token, err := jwt.CreateToken(uint(1), "username", config.SecretJWT)
			assert.NoError(t, err)

			req.Header.Set("Authorization", token)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := spotify.SearchResponse{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}

func Test_handler_UpsertActivity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSvc := NewMockSpotifyService(mockCtrl)

	isLiked := true

	tests := []struct {
		name string
		mockFn func()
		requestBody spotify.TrackActivityRequest
		expectedStatusCode int
	}{
		// TODO: Add test cases.
		{
			name: "success",
			requestBody: spotify.TrackActivityRequest{
				SpotifyID: "SpotifyID",
				IsLiked: &isLiked,
			},
			expectedStatusCode: 200,
			mockFn: func() {
				mockSvc.EXPECT().UpSertActivity(gomock.Any(), uint(1), spotify.TrackActivityRequest{
					SpotifyID: "SpotifyID",
					IsLiked: &isLiked,
				}).Return(nil)
			},
		},
		{
			name: "failed",
			requestBody: spotify.TrackActivityRequest{
				SpotifyID: "SpotifyID",
				IsLiked: &isLiked,
			},
			expectedStatusCode: 400,
			mockFn: func() {
				mockSvc.EXPECT().UpSertActivity(gomock.Any(), uint(1), spotify.TrackActivityRequest{
					SpotifyID: "SpotifyID",
					IsLiked: &isLiked,
				}).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			config, err := configs.Init("../../configs", "env", "test.env")
			assert.NoError(t, err)

			gin.SetMode(gin.ReleaseMode)

			r := gin.New()
			route := r.Group("/api/v1")

			h := &handler{
				route: route,
				service: mockSvc,
			}

			h.RegisterRoute()

			w := httptest.NewRecorder()

			endpoint := "/api/v1/spotify/activity"
			bodyBytes, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			httpReq, err := http.NewRequest(http.MethodPost, endpoint, io.NopCloser(bytes.NewBuffer(bodyBytes)))
			assert.NoError(t, err)

			token, err := jwt.CreateToken(uint(1), "username", config.SecretJWT)
			assert.NoError(t, err)

			httpReq.Header.Set("Authorization", token)

			r.ServeHTTP(w, httpReq)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
