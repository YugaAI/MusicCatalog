package spotify

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/YugaAI/MusicCatalog/internal/configs"
	"github.com/YugaAI/MusicCatalog/pkg/httpclient"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_outbound_Search(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHTTPClient := httpclient.NewMockHTTPClient(mockCtrl)
	next := "https://api.spotify.com/v1/search?offset=11&limit=10&query=bohemian%20rhapsody&type=track&market=ID&locale=id,en-US;q%3D0.9,en;q%3D0.8"
	previous := "https://api.spotify.com/v1/search?offset=0&limit=10&query=bohemian%20rhapsody&type=track&market=ID&locale=id,en-US;q%3D0.9,en;q%3D0.8"
	type args struct {
		Query   string
		Limit   int
		Offsite int
	}
	tests := []struct {
		name    string
		args    args
		want    *SpotifySearchResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "successful search",
			args: args{
				Query:   "bohemian rhapsody",
				Limit:   10,
				Offsite: 1,
			},
			want: &SpotifySearchResponse{
				Tracks: SpotifyTracks{
					Href:     "https://api.spotify.com/v1/search?offset=1&limit=10&query=bohemian%20rhapsody&type=track&market=ID&locale=id,en-US;q%3D0.9,en;q%3D0.8",
					Limit:    10,
					Next:     &next,
					Offsite:  1,
					Previous: &previous,
					Total:    23,
					Items: []SpotifyTrackObject{
						{
							Album: SpotifyAlbumObject{
								AlbumType:   "album",
								TotalTracks: 22,
								Images: []SpotifyAlbumImages{
									{
										URL: `https://i.scdn.co/image/ab67616d0000b273e8b066f70c206551210d902b`,
									},
									{
										URL: `https://i.scdn.co/image/ab67616d00001e02e8b066f70c206551210d902b`,
									},
									{
										URL: `https://i.scdn.co/image/ab67616d00004851e8b066f70c206551210d902b`,
									},
								},
								Name: "Bohemian Rhapsody (The Original Soundtrack)",
							},
							Artists: []SpotifyArtistObject{
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
								AlbumType:   "compilation",
								TotalTracks: 17,
								Images: []SpotifyAlbumImages{
									{
										URL: `https://i.scdn.co/image/ab67616d0000b273bb19d0c22d5709c9d73c8263`,
									},
									{
										URL: `https://i.scdn.co/image/ab67616d00001e02bb19d0c22d5709c9d73c8263`,
									},
									{
										URL: `https://i.scdn.co/image/ab67616d00004851bb19d0c22d5709c9d73c8263`,
									},
								},
								Name: "Greatest Hits (Remastered)",
							},
							Artists: []SpotifyArtistObject{
								{
									Href: "https://api.spotify.com/v1/artists/1dfeR4HaWDbWqFHLkxsg1d",
									Name: "Queen",
								},
							},
							Explicit: false,
							Href:     "https://api.spotify.com/v1/tracks/2OBofMJx94NryV2SK8p8Zf",
							ID:       "2OBofMJx94NryV2SK8p8Zf",
							Name:     "Bohemian Rhapsody - Remastered 2011",
						},
					},
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				params := url.Values{}
				params.Set("q", args.Query)
				params.Set("type", "track")
				params.Set("limit", strconv.Itoa(args.Limit))
				params.Set("offset", strconv.Itoa(args.Offsite))
				basePath := `https://api.spotify.com/v1/search`
				urlPath := fmt.Sprintf("%s?%s", basePath, params.Encode())
				req, err := http.NewRequest(http.MethodGet, urlPath, nil)
				assert.NoError(t, err)

				req.Header.Set("Authorization", "Bearer accessToken")
				mockHTTPClient.EXPECT().Do(req).Return(&http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(searchResponse)),
				}, nil)
			},
		},
		{
			name: "failed search",
			args: args{
				Query:   "bohemian rhapsody",
				Limit:   10,
				Offsite: 1,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				params := url.Values{}
				params.Set("q", args.Query)
				params.Set("type", "track")
				params.Set("limit", strconv.Itoa(args.Limit))
				params.Set("offset", strconv.Itoa(args.Offsite))
				basePath := `https://api.spotify.com/v1/search`
				urlPath := fmt.Sprintf("%s?%s", basePath, params.Encode())
				req, err := http.NewRequest(http.MethodGet, urlPath, nil)
				assert.NoError(t, err)

				req.Header.Set("Authorization", "Bearer accessToken")
				mockHTTPClient.EXPECT().Do(req).Return(&http.Response{
					StatusCode: 500,
					Body:       io.NopCloser(bytes.NewBufferString(`Internal server ERROR`)),
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			o := &outbound{
				cfg:         &configs.Config{},
				client:      mockHTTPClient,
				AccessToken: "accessToken",
				TokenType:   "Bearer",
				ExpiredAt:   time.Now().Add(time.Hour),
			}
			got, err := o.Search(context.Background(), tt.args.Query, tt.args.Limit, tt.args.Offsite)
			if (err != nil) != tt.wantErr {
				t.Errorf("outbound.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got, "unexpected result from outbound.Search")
		})
	}
}
