package trackactivities

import "gorm.io/gorm"

type (
	TrackActivities struct {
		gorm.Model
		UserID    uint   `gorm:"not null"`
		SpotifyID string `gorm:"not null"`
		IsLiked   bool   `gorm:"not null"`
		CreatedBy string `gorm:"not null"`
		UpdatedBy string `gorm:"not null"`
	}
)

type(
	TrackActivitiesRequest struct {
		SpotifyID string `json:"spotify_id"`
		IsLiked   bool   `json:"isLiked"` //true=like, false=islike, null=neutral
	}
)