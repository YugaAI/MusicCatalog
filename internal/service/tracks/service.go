package tracks

import (
	"context"

	"github.com/YugaAI/MusicCatalog/internal/repository/spotify"
)
//go:generate mockgen -source=service.go -destination=service_mock.go -package=tracks
type SpotifyOutbound interface {
	Search(ctx context.Context, query string, limit, offsite int) (*spotify.SpotifySearchResponse, error)
}

type service struct {
	spotifyOutbound SpotifyOutbound
}

func NewService(spotifyOutbound SpotifyOutbound) *service {
	return &service{spotifyOutbound: spotifyOutbound}
}
