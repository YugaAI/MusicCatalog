package tracks

import (
	"context"

	"github.com/YugaAI/MusicCatalog/internal/middleware"
	"github.com/YugaAI/MusicCatalog/internal/models/spotify"
	"github.com/YugaAI/MusicCatalog/internal/models/trackactivities"
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=hendler.go -destination=hendler_mock.go -package=tracks
type service interface {
	Search(ctx context.Context, query string, pageSize, pageIndex int, userID uint) (*spotify.SearchResponse, error)
	UpsertTrackActivity(ctx context.Context, userID uint, request trackactivities.TrackActivitiesRequest) error
}
type Handler struct {
	*gin.Engine
	service service
}

func NewHandler(api *gin.Engine, service service) *Handler {
	return &Handler{
		api,
		service,
	}
}

func (h *Handler) RegisterRoutes() {
	route := h.Group("/tracks")
	route.Use(middleware.AuthMiddlewere())
	route.GET("/search", h.Search)
	route.POST("/activity", h.UpsertTrackActivity)
}
