package tracks

import (
	"context"
	"fmt"
	"testing"

	"github.com/YugaAI/MusicCatalog/internal/models/trackactivities"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_service_UpsertTrackActivity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTrackActivityRepo := NewMockTracksActivityRepository(mockCtrl)

	isLikeTrue := true
	isLikeFalse := false

	type args struct {
		userID  uint
		request trackactivities.TrackActivitiesRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "Success : create",
			args: args{
				userID: 1,
				request: trackactivities.TrackActivitiesRequest{
					SpotifyID: "spotifyID",
					IsLiked:   isLikeTrue,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockTrackActivityRepo.EXPECT().Get(gomock.Any(), args.userID, args.request.SpotifyID).Return(nil, gorm.ErrRecordNotFound)

				mockTrackActivityRepo.EXPECT().Create(gomock.Any(), trackactivities.TrackActivities{
					UserID:    args.userID,
					SpotifyID: args.request.SpotifyID,
					IsLiked:   args.request.IsLiked,
					CreatedBy: fmt.Sprintf("%d", args.userID),
					UpdatedBy: fmt.Sprintf("%d", args.userID),
				}).Return(nil)
			},
		},
		{
			name: "Success : Update",
			args: args{
				userID: 1,
				request: trackactivities.TrackActivitiesRequest{
					SpotifyID: "spotifyID",
					IsLiked:   isLikeTrue,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockTrackActivityRepo.EXPECT().Get(gomock.Any(), args.userID, args.request.SpotifyID).Return(&trackactivities.TrackActivities{
					IsLiked: isLikeFalse,
				}, nil)

				mockTrackActivityRepo.EXPECT().Update(gomock.Any(), trackactivities.TrackActivities{
					IsLiked: args.request.IsLiked,
				}).Return(nil)
			},
		},
		{
			name: "Error",
			args: args{
				userID: 1,
				request: trackactivities.TrackActivitiesRequest{
					SpotifyID: "spotifyID",
					IsLiked:   isLikeTrue,
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockTrackActivityRepo.EXPECT().Get(gomock.Any(), args.userID, args.request.SpotifyID).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		tt.mockFn(tt.args)
		s := &service{
			trackAktivityRepo: mockTrackActivityRepo,
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := s.UpsertTrackActivity(context.Background(), tt.args.userID, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("service.UpsertTrackActivity() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
