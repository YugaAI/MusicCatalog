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

type SpotifySearchResponse struct {
	Tracks SpotifyTracks `json:"tracks"` // Assuming spotifyTracks is a struct that holds track details
}

type SpotifyTracks struct {
	Href     string               `json:"href"`
	Limit    int                  `json:"limit"`
	Next     *string              `json:"next"`
	Offsite  int                  `json:"offset"`
	Previous *string              `json:"previous"`
	Total    int                  `json:"total"`
	Items    []SpotifyTrackObject `json:"items"` // Assuming spotifyTrack is a struct that holds individual track details
}

type SpotifyTrackObject struct {
	Album    SpotifyAlbumObject    `json:"album"`    // Assuming spotifyAlbum is a struct that holds album details
	Artists  []SpotifyArtistObject `json:"artists"`  // Assuming spotifyArtist is a struct that holds artist details
	Explicit bool                  `json:"explicit"` // Assuming Explicit is a boolean indicating if the track is explicit
	Href     string                `json:"href"`
	ID       string                `json:"id"`   // ID of the track
	Name     string                `json:"name"` // Name of the track
}

type SpotifyAlbumObject struct {
	AlbumType   string               `json:"album_type"`   // Assuming AlbumType is a string representing the type of album
	TotalTracks int                  `json:"total_tracks"` // Assuming TotalTracks is an integer representing the total number of tracks in the album
	Images      []SpotifyAlbumImages `json:"images"`       // Assuming Images is a struct that holds image details
	Name        string               `json:"name"`         // Name of the album
}

type SpotifyAlbumImages struct {
	URL string `json:"url"` // URL of the image
}

type SpotifyArtistObject struct {
	Href string `json:"href"`
	Name string `json:"name"` // Name of the artist
}

func (o *outbound) Search(ctx context.Context, Query string, Limit, Offsite int) (*SpotifySearchResponse, error) {

	params := url.Values{}
	params.Set("q", Query)
	params.Set("type", "track")
	params.Set("limit", strconv.Itoa(Limit))
	params.Set("offset", strconv.Itoa(Offsite))

	basePath := `https://api.spotify.com/v1/search`
	urlPath := fmt.Sprintf("%s?%s", basePath, params.Encode())
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to create search request for spotify token")
		return nil, err
	}

	accessToken, tokenType, err := o.GetTokenDetails()
	if err != nil {
		log.Error().Err(err).Msg("failed to get token details for spotify search")
		return nil, err
	}
	bearerToken := fmt.Sprintf("%s %s", tokenType, accessToken)
	// Set the Authorization header with the Bearer token
	req.Header.Set("Authorization", bearerToken)

	resp, err := o.client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get response from spotify search endpoint")
		return nil, err
	}

	defer resp.Body.Close()

	var response SpotifySearchResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode search spotify token response")
		return nil, err
	}
	return &response, nil
}
