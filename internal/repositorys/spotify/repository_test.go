package spotify

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sgitwhyd/music-catalogue/internal/configs"
	"github.com/sgitwhyd/music-catalogue/internal/models/spotify"
	"github.com/sgitwhyd/music-catalogue/pkg/httpclient"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_outbond_Search(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHttpClient := httpclient.NewMockHTTPClient(mockCtrl)

	next := "https://api.spotify.com/v1/search?query=bohemian+rhapsody&type=track&market=ID&locale=en-US%2Cen%3Bq%3D0.9&offset=10&limit=10"

	type args struct {
		query  string
		limit  int
		offset int
	}
	tests := []struct {
		name    string
		args    args
		want    *SpotifySearchResponse
		wantErr bool
		mockFn  func(args args)
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				query:  "bohemian rhapsody",
				limit:  10,
				offset: 0,
			},
			want: &SpotifySearchResponse{
				Tracks: SpotifyTrack{
					Href:   "https://api.spotify.com/v1/search?query=bohemian+rhapsody&type=track&market=ID&locale=en-US%2Cen%3Bq%3D0.9&offset=0&limit=10",
					Limit:  10,
					Next:   &next,
					Offset: 0,
					Total:  905,
					Items: []SpotifyTrackObject{
						{
							Album: SpotifyAlbumObject{
								AlbumType:   "album",
								TotalTracks: 22,
								Images: []SpotifyImagesObject{
									{
										URL: "https://i.scdn.co/image/ab67616d0000b273e8b066f70c206551210d902b",
									},
									{
										URL: "https://i.scdn.co/image/ab67616d00001e02e8b066f70c206551210d902b",
									},
									{
										URL: "https://i.scdn.co/image/ab67616d00004851e8b066f70c206551210d902b",
									},
								},
								Name: "Bohemian Rhapsody (The Original Soundtrack)",
							},
							Artists: []SpotifyArtisObject{
								{
									Href: "https://api.spotify.com/v1/artists/1dfeR4HaWDbWqFHLkxsg1d",
									Name: "Queen",
								},
							},
							Explicit: false,
							Href:     "https://api.spotify.com/v1/tracks/3z8h0TU7ReDPLIbEnYhWZb",
							ID:       "3z8h0TU7ReDPLIbEnYhWZb",
							Name:     "Bohemian Rhapsody",
						},
						{
							Album: SpotifyAlbumObject{
								AlbumType:   "album",
								TotalTracks: 12,
								Images: []SpotifyImagesObject{
									{
										URL: "https://i.scdn.co/image/ab67616d0000b273e319baafd16e84f0408af2a0",
									},
									{
										URL: "https://i.scdn.co/image/ab67616d00001e02e319baafd16e84f0408af2a0",
									},
									{
										URL: "https://i.scdn.co/image/ab67616d00004851e319baafd16e84f0408af2a0",
									},
								},
								Name: "A Night At The Opera (2011 Remaster)",
							},
							Artists: []SpotifyArtisObject{
								{
									Href: "https://api.spotify.com/v1/artists/1dfeR4HaWDbWqFHLkxsg1d",
									Name: "Queen",
								},
							},
							Explicit: false,
							Href:     "https://api.spotify.com/v1/tracks/4u7EnebtmKWzUH433cf5Qv",
							ID:       "4u7EnebtmKWzUH433cf5Qv",
							Name:     "Bohemian Rhapsody - Remastered 2011",
						},
					},
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				params := url.Values{}
				params.Set("q", args.query)
				params.Set("type", "track")
				params.Set("limit", strconv.Itoa(args.limit))
				params.Set("offset", strconv.Itoa(args.offset))

				basePath := `https://api.spotify.com/v1/search`
				urlPath := fmt.Sprintf("%s?%s", basePath, params.Encode())
				req, err := http.NewRequest(http.MethodGet, urlPath, nil)
				assert.NoError(t, err)

				req.Header.Set("Authorization", "Bearer accessToken")
				mockHttpClient.EXPECT().Do(req).Return(&http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(searchResponse)),
				}, nil)
			},
		},
		{
			name: "failed",
			args: args{
				query:  "bohemian rhapsody",
				limit:  10,
				offset: 0,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				params := url.Values{}
				params.Set("q", args.query)
				params.Set("type", "track")
				params.Set("limit", strconv.Itoa(args.limit))
				params.Set("offset", strconv.Itoa(args.offset))

				basePath := `https://api.spotify.com/v1/search`
				urlPath := fmt.Sprintf("%s?%s", basePath, params.Encode())
				req, err := http.NewRequest(http.MethodGet, urlPath, nil)
				assert.NoError(t, err)

				req.Header.Set("Authorization", "Bearer accessToken")
				mockHttpClient.EXPECT().Do(req).Return(&http.Response{
					StatusCode: 500,
					Body:       io.NopCloser(bytes.NewBufferString(`internal server error`)),
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			o := &outbond{
				cfg:         &configs.Config{},
				client:      mockHttpClient,
				AccessToken: "accessToken",
				TokenType:   "Bearer",
				ExpiredAt:   time.Now().Add(1 * time.Hour),
			}

			got, err := o.Search(context.Background(), tt.args.query, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("outbond.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("outbond.Search() = %v,\n want %v", got, tt.want)
			}
		})
	}
}

func Test_spotifyRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	isLiked := true
	now := time.Now()
	type args struct {
		model spotify.TrackActivity
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				model: spotify.TrackActivity{
					Model: gorm.Model{
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:    1,
					SpotifyID: "spotifyID",
					IsLiked:   &isLiked,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "track_activities" (.+) VALUES (.+)`).WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					args.model.UserID,
					args.model.SpotifyID,
					args.model.IsLiked,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
		},
		{
			name: "failed",
			args: args{
				model: spotify.TrackActivity{
					Model: gorm.Model{
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:    1,
					SpotifyID: "spotifyID",
					IsLiked:   &isLiked,
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "track_activities" (.+) VALUES (.+)`).WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					args.model.UserID,
					args.model.SpotifyID,
					args.model.IsLiked,
				).WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &spotifyRepository{
				db: gormDB,
			}

			if err := r.Create(context.Background(), tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("spotifyRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_spotifyRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	now := time.Now()
	isLiked := true

	type args struct {
		model spotify.TrackActivity
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				model: spotify.TrackActivity{
					Model: gorm.Model{
						ID:        123,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:    1,
					SpotifyID: "spotifyID",
					IsLiked:   &isLiked,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "track_activities" SET (.+) WHERE (.+)`).WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					args.model.UserID,
					args.model.SpotifyID,
					args.model.IsLiked,
					args.model.ID,
				).WillReturnResult(sqlmock.NewResult(123, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "failed",
			args: args{
				model: spotify.TrackActivity{
					Model: gorm.Model{
						ID:        123,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:    1,
					SpotifyID: "spotifyID",
					IsLiked:   &isLiked,
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "track_activities" SET (.+) WHERE (.+)`).WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					args.model.UserID,
					args.model.SpotifyID,
					args.model.IsLiked,
					args.model.ID,
				).WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			r := &spotifyRepository{
				db: gormDB,
			}

			if err := r.Update(context.Background(), tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("spotifyRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_spotifyRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	now := time.Now()
	isLiked := true

	type args struct {
		UserID    uint
		spotifyID string
	}
	tests := []struct {
		name    string
		args    args
		want    *spotify.TrackActivity
		wantErr bool
		mockFn  func(args args)
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				UserID:    1,
				spotifyID: "spotifyID",
			},
			want: &spotify.TrackActivity{
				UserID:    1,
				SpotifyID: "spotifyID",
				IsLiked:   &isLiked,
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "track_activities" WHERE .*`).WithArgs(
					args.UserID,
					args.spotifyID,
					1,
				).WillReturnRows(sqlmock.NewRows(
					[]string{"id", "created_at", "updated_at", "user_id", "spotify_id", "is_liked"},
				).AddRow(1, now, now, args.UserID, args.spotifyID, true))
			},
		},
		{
			name: "failed",
			args: args{
				UserID:    1,
				spotifyID: "spotifyID",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "track_activities" WHERE .*`).WithArgs(
					args.UserID,
					args.spotifyID,
					1,
				).WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &spotifyRepository{
				db: gormDB,
			}
			got, err := r.Get(context.Background(), tt.args.UserID, tt.args.spotifyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("spotifyRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("spotifyRepository.Get() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_spotifyRepository_GetBulkSpotifyIDs(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	now := time.Now()
	isLiked := true

	type args struct {
		UserID     uint
		spotifyIDs []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]spotify.TrackActivity
		wantErr bool
		mockFn	func (args args)
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				UserID: 1,
				spotifyIDs: []string{
					"spotifyID", 
				},
			},
			want: map[string]spotify.TrackActivity{
				"spotifyID": {
					UserID: 1,
					SpotifyID: "spotifyID",
					IsLiked: &isLiked,
					Model: gorm.Model{
						ID: 1,
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			wantErr: false,
			mockFn: func(args args) {

				mock.ExpectQuery(`SELECT \* FROM "track_activities" .+`).
					WithArgs(
						args.UserID,
						strings.Join(args.spotifyIDs, ","),
					).
					WillReturnRows(sqlmock.NewRows(
					[]string{"id", "created_at", "updated_at", "user_id", "spotify_id", "is_liked"},
				).AddRow(1, now, now, args.UserID, "spotifyID", true))

			},
		},
		{
			name: "failed",
			args: args{
				UserID: 1,
				spotifyIDs: []string{
					"spotifyID", 
				},
			},
			want: nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "track_activities" .+`).
					WithArgs(
						args.UserID,
						strings.Join(args.spotifyIDs, ","),
					).
					WillReturnError(assert.AnError)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			r := &spotifyRepository{
				db: gormDB,
			}

			got, err := r.GetBulkSpotifyIDs(context.Background(), tt.args.UserID, tt.args.spotifyIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("spotifyRepository.GetBulkSpotifyIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("spotifyRepository.GetBulkSpotifyIDs() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
