package spotify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/music-catalogue/internal/models/spotify"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_handler_Search(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSvc := NewMockSpotifyService(mockCtrl)

	tests := []struct {
		name string
		mockFn func()
		expectedCode int
		expectedResponse spotify.SearchResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			expectedCode: 200,
			wantErr: false,
			expectedResponse: 
				spotify.SearchResponse{
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
				mockSvc.EXPECT().Search(gomock.Any(), "bohemian rhapsody", 10, 1).Return(&spotify.SearchResponse{
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
				}, nil).Times(1)
			},
		},
		{
			name: "failed",
			expectedCode: 400,
			wantErr: true,
			expectedResponse: spotify.SearchResponse{},
			mockFn: func() {
				mockSvc.EXPECT().Search(gomock.Any(), "bohemian rhapsody", 10, 1).Return(nil, assert.AnError).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			r := gin.New()
			route := r.Group("/api/v1")

			h := &handler{
				RouterGroup:  route,
				service: mockSvc,
			}
			h.RegisterRoute()

			endpoint := "/api/v1/spotify/search?query=bohemian+rhapsody&pageIndex=1&pageSize=10"
			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			assert.NoError(t, err)

		
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
