package tracks

import (
	"context"

	"github.com/YugaAI/MusicCatalog/internal/models/trackactivities"
	"github.com/YugaAI/MusicCatalog/internal/repository/spotify"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=tracks
type SpotifyOutbound interface {
	Search(ctx context.Context, query string, limit, offsite int) (*spotify.SpotifySearchResponse, error)
}

type TracksActivityRepository interface {
	Create(ctx context.Context, model trackactivities.TrackActivities) error
	Update(ctx context.Context, model trackactivities.TrackActivities) error
	Get(ctx context.Context, userID uint, spotifyID string) (*trackactivities.TrackActivities, error)
	GetBulkSpotifyID(ctx context.Context, userID uint, spotifyIDs []string) (map[string]trackactivities.TrackActivities, error)
}

type service struct {
	spotifyOutbound   SpotifyOutbound
	trackAktivityRepo TracksActivityRepository
}

func NewService(spotifyOutbound SpotifyOutbound, trackAktivityRepo TracksActivityRepository) *service {
	return &service{
		spotifyOutbound:   spotifyOutbound,
		trackAktivityRepo: trackAktivityRepo,
	}
}
