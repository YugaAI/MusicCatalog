package trackactivities

import (
	"context"

	"github.com/YugaAI/MusicCatalog/internal/models/trackactivities"
)

func (r *Repository) Create(ctx context.Context, model trackactivities.TrackActivities) error {
	return r.db.Create(&model).Error
}

func (r *Repository) Update(ctx context.Context, model trackactivities.TrackActivities) error {
	return r.db.Save(&model).Error
}

func (r *Repository) Get(ctx context.Context, userID uint, spotifyID string) (*trackactivities.TrackActivities, error) {
	trackActivity := trackactivities.TrackActivities{}
	res := r.db.Where("user_id = ? AND spotify_id = ?", userID, spotifyID).First(&trackActivity)
	if res.Error != nil {
		return nil, res.Error
	}
	return &trackActivity, nil
}

// digunakan untuk mengintegasikan dengan trcks di bagian search
func (r *Repository) GetBulkSpotifyID(ctx context.Context, userID uint, spotifyIDs []string) (map[string]trackactivities.TrackActivities, error) {
	trackActivities := make([]trackactivities.TrackActivities, 0)
	res := r.db.Where("user_id = ? AND spotify_id IN ?", userID, spotifyIDs).Find(&trackActivities)
	if res.Error != nil {
		return nil, res.Error
	}

	//
	result:= make(map[string]trackactivities.TrackActivities,0)
	for _, trackActivity := range trackActivities {
		result[trackActivity.SpotifyID] = trackActivity
	}
	return result, nil
}
