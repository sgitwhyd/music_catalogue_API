package spotify

import (
	"context"
	"reflect"
	"testing"

	"github.com/sgitwhyd/music-catalogue/internal/models/spotify"
	spotifyRepo "github.com/sgitwhyd/music-catalogue/internal/repositorys/spotify"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_spotifyService_Search(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSpotifyOutbond := NewMockSpotifyOutbond(mockCtrl)
	mockSpotifyRepo := NewMockSpotifyRepository(mockCtrl)

	isLikedTrue := true
	isLikedFalse := false

	createMockResponse := func() *spotifyRepo.SpotifySearchResponse {
		return &spotifyRepo.SpotifySearchResponse{
			Tracks: spotifyRepo.SpotifyTrack{
				Href:   "https://api.spotify.com/v1/search?query=bohemian+rhapsody&type=track&market=ID&locale=en-US%2Cen%3Bq%3D0.9&offset=0&limit=10",
				Limit:  10,
				Next:   func(s string) *string { return &s }("https://api.spotify.com/v1/search?query=bohemian+rhapsody&type=track&market=ID&locale=en-US%2Cen%3Bq%3D0.9&offset=10&limit=10"),
				Offset: 0,
				Total:  905,
				Items: []spotifyRepo.SpotifyTrackObject{
					{
						Album: spotifyRepo.SpotifyAlbumObject{
							AlbumType:   "album",
							TotalTracks: 22,
							Images: []spotifyRepo.SpotifyImagesObject{
								{URL: "https://i.scdn.co/image/ab67616d0000b273e8b066f70c206551210d902b"},
								{URL: "https://i.scdn.co/image/ab67616d00001e02e8b066f70c206551210d902b"},
								{URL: "https://i.scdn.co/image/ab67616d00004851e8b066f70c206551210d902b"},
							},
							Name: "Bohemian Rhapsody (The Original Soundtrack)",
						},
						Artists: []spotifyRepo.SpotifyArtisObject{
							{Name: "Queen"},
						},
						Explicit: false,
						ID:       "3z8h0TU7ReDPLIbEnYhWZb",
						Name:     "Bohemian Rhapsody",
					},
					{
						Album: spotifyRepo.SpotifyAlbumObject{
							AlbumType:   "album",
							TotalTracks: 12,
							Images: []spotifyRepo.SpotifyImagesObject{
								{URL: "https://i.scdn.co/image/ab67616d0000b273e319baafd16e84f0408af2a0"},
								{URL: "https://i.scdn.co/image/ab67616d00001e02e319baafd16e84f0408af2a0"},
								{URL: "https://i.scdn.co/image/ab67616d00004851e319baafd16e84f0408af2a0"},
							},
							Name: "A Night At The Opera (2011 Remaster)",
						},
						Artists: []spotifyRepo.SpotifyArtisObject{
							{Name: "Queen"},
						},
						Explicit: false,
						ID:       "4u7EnebtmKWzUH433cf5Qv",
						Name:     "Bohemian Rhapsody - Remastered 2011",
					},
				},
			},
		}
	}

	type args struct {
		query     string
		pageSize  int
		pageIndex int
	}
	tests := []struct {
		name    string
		args    args
		want    *spotify.SearchResponse
		wantErr bool
		mockFn  func(args args)
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				query:     "bohemian rhapsody",
				pageSize:  10,
				pageIndex: 1,
			},
			want: &spotify.SearchResponse{
				Limit:  10,
				Offset: 0,
				Items: []spotify.SpotifyTrackObjectResponse{
					{
						AlbumType:        "album",
						AlbumTotalTracks: 22,
						AlbumImagesURL:   []string{"https://i.scdn.co/image/ab67616d0000b273e8b066f70c206551210d902b", "https://i.scdn.co/image/ab67616d00001e02e8b066f70c206551210d902b", "https://i.scdn.co/image/ab67616d00004851e8b066f70c206551210d902b"},
						AlbumName:        "Bohemian Rhapsody (The Original Soundtrack)",
						ArtistsName:      []string{"Queen"},
						Explicit:         false,
						ID:               "3z8h0TU7ReDPLIbEnYhWZb",
						Name:             "Bohemian Rhapsody",
						IsLiked:          &isLikedTrue,
					},
					{
						AlbumType:        "album",
						AlbumTotalTracks: 12,
						AlbumImagesURL:   []string{"https://i.scdn.co/image/ab67616d0000b273e319baafd16e84f0408af2a0", "https://i.scdn.co/image/ab67616d00001e02e319baafd16e84f0408af2a0", "https://i.scdn.co/image/ab67616d00004851e319baafd16e84f0408af2a0"},
						AlbumName:        "A Night At The Opera (2011 Remaster)",
						ArtistsName:      []string{"Queen"},
						Explicit:         false,
						ID:               "4u7EnebtmKWzUH433cf5Qv",
						Name:             "Bohemian Rhapsody - Remastered 2011",
						IsLiked:          &isLikedFalse,
					},
				},
				Total: 905,
			},
			wantErr: false,
			mockFn: func(args args) {
				mockSpotifyOutbond.EXPECT().Search(gomock.Any(), "bohemian rhapsody", 10, 0).Return(createMockResponse(), nil)
				mockSpotifyRepo.EXPECT().GetBulkSpotifyIDs(gomock.Any(), uint(1), []string{
					"3z8h0TU7ReDPLIbEnYhWZb",
					"4u7EnebtmKWzUH433cf5Qv",
				}).Return(map[string]spotify.TrackActivity{
					"3z8h0TU7ReDPLIbEnYhWZb": {
						IsLiked: &isLikedTrue,
					},
					"4u7EnebtmKWzUH433cf5Qv": {
						IsLiked: &isLikedFalse,
					},
				}, nil)
			},
		},
		{
			name: "failed",
			args: args{
				query:     "bohemian rhapsody",
				pageSize:  10,
				pageIndex: 1,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockSpotifyOutbond.EXPECT().Search(gomock.Any(), "bohemian rhapsody", 10, 0).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &spotifyService{
				spotifyOutbond: mockSpotifyOutbond,
				spotifyRepo:    mockSpotifyRepo,
			}

			got, err := s.Search(context.Background(), tt.args.query, tt.args.pageSize, tt.args.pageIndex, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("spotifyService.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot = %v, \nwant %v", got, tt.want)
			}
		})
	}
}

func Test_spotifyService_UpSertActivity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSpotifyOutbond := NewMockSpotifyOutbond(mockCtrl)
	mockSpotifyRepo := NewMockSpotifyRepository(mockCtrl)

	isLikedTrue := true
	isLikedFalse := false

	type args struct {
		userID  uint
		request spotify.TrackActivityRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn func(args args)
	}{
		// TODO: Add test cases.
		{
			name: "success_create",
			args: args{
				userID: uint(1),
				request: spotify.TrackActivityRequest{
					SpotifyID: "SpotifyID",
					IsLiked: &isLikedTrue,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockSpotifyRepo.EXPECT().Get(gomock.Any(), uint(1), args.request.SpotifyID).
					Return(nil, gorm.ErrRecordNotFound)
				
				mockSpotifyRepo.EXPECT().Create(gomock.Any(), spotify.TrackActivity{
					UserID: args.userID,
					SpotifyID: args.request.SpotifyID,
					IsLiked: args.request.IsLiked,
				}).Return(nil)
			},
		},
		{
			name: "success_update",
			args: args{
				userID: uint(1),
				request: spotify.TrackActivityRequest{
					SpotifyID: "SpotifyID",
					IsLiked: &isLikedTrue,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockSpotifyRepo.EXPECT().Get(gomock.Any(), args.userID, args.request.SpotifyID).
					Return(&spotify.TrackActivity{
						IsLiked: &isLikedFalse,
					}, nil)
				
				mockSpotifyRepo.EXPECT().Update(gomock.Any(), spotify.TrackActivity{
					IsLiked: args.request.IsLiked,
				}).Return(nil)
			},
		},
		{
			name: "error",
			args: args{
				userID: uint(1),
				request: spotify.TrackActivityRequest{
					SpotifyID: "SpotifyID",
					IsLiked: &isLikedTrue,
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockSpotifyRepo.EXPECT().Get(gomock.Any(), uint(1), args.request.SpotifyID).
					Return(nil, nil)
				
				mockSpotifyRepo.EXPECT().Create(gomock.Any(), spotify.TrackActivity{
					UserID: args.userID,
					SpotifyID: args.request.SpotifyID,
					IsLiked: args.request.IsLiked,
				}).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			s := &spotifyService{
				spotifyOutbond: mockSpotifyOutbond,
				spotifyRepo: mockSpotifyRepo,
			}

			if err := s.UpSertActivity(context.Background(), tt.args.userID, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("spotifyService.UpSertActivity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
