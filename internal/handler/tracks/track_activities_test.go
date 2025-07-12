package tracks

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YugaAI/MusicCatalog/internal/models/trackactivities"
	"github.com/YugaAI/MusicCatalog/pkg/jwt"
	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_UpsertTrackActivity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := NewMockservice(mockCtrl)

	isLikeTrue := true
	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
	}{
		{
			name: "Success",
			mockFn: func() {
				mockService.EXPECT().UpsertTrackActivity(gomock.Any(), uint(1), trackactivities.TrackActivitiesRequest{
					SpotifyID: "spotifyID",
					IsLiked:   isLikeTrue,
				}).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name: "Error",
			mockFn: func() {
				mockService.EXPECT().UpsertTrackActivity(gomock.Any(), uint(1), trackactivities.TrackActivitiesRequest{
					SpotifyID: "spotifyID",
					IsLiked:   isLikeTrue,
				}).Return(assert.AnError)
			},
			expectedStatusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()

			h := &Handler{
				Engine:  api,
				service: mockService,
			}
			w := httptest.NewRecorder()
			h.RegisterRoutes()

			endpoint := "/tracks/activity"

			payload := trackactivities.TrackActivitiesRequest{
				SpotifyID: "spotifyID",
				IsLiked:   isLikeTrue,
			}

			payloadBytes, err := json.Marshal(payload)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, endpoint, io.NopCloser(bytes.NewBuffer(payloadBytes)))
			assert.NoError(t, err)

			token, err := jwt.CreateToken(1, "username", "")
			assert.NoError(t, err)
			req.Header.Set("Authorization", token)

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
