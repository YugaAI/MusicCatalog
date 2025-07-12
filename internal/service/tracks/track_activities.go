package tracks

import (
	"context"
	"fmt"

	"github.com/YugaAI/MusicCatalog/internal/models/trackactivities"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (s *service) UpsertTrackActivity(ctx context.Context, userID uint, request trackactivities.TrackActivitiesRequest) error {
	activity, err := s.trackAktivityRepo.Get(ctx, userID, request.SpotifyID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Msg("failed to get track activity")
		return err
	}
	if err == gorm.ErrRecordNotFound || activity == nil {
		//create user activity
		err = s.trackAktivityRepo.Create(ctx, trackactivities.TrackActivities{
			UserID:    userID,
			SpotifyID: request.SpotifyID,
			IsLiked:   request.IsLiked,
			CreatedBy: fmt.Sprintf("%d", userID),
			UpdatedBy: fmt.Sprintf("%d", userID),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to create track activity")
			return err
		}
		return err
	}
	activity.IsLiked = request.IsLiked
	err = s.trackAktivityRepo.Update(ctx, *activity)
	if err != nil {
		log.Error().Err(err).Msg("failed to update track activity")
		return err
	}
	return nil
}
