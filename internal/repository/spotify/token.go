package spotify

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (o *outbound) GetTokenDetails() (string, string, error) {
	if o.AccessToken == "" || time.Now().After(o.ExpiredAt) {
		//call spotify API to get a new token
		err := o.generateToken()
		if err != nil {
			return "", "", err
		}
	}
	return o.AccessToken, o.TokenType, nil
}

func (o *outbound) generateToken() error {
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", o.cfg.SpotifyConfig.ClientID)
	formData.Set("client_secret", o.cfg.SpotifyConfig.ClientSecret)

	encodeURL := formData.Encode()

	req, err := http.NewRequest(http.MethodPost, `https://accounts.spotify.com/api/token`, strings.NewReader(encodeURL))
	if err != nil {
		log.Error().Err(err).Msg("failed to create request for spotify token")
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get response from spotify token endpoint")
		return err
	}

	defer resp.Body.Close()

	var response SpotifyTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode spotify token response")
		return err
	}

	o.AccessToken = response.AccessToken
	o.TokenType = response.TokenType
	o.ExpiredAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	return nil
}
