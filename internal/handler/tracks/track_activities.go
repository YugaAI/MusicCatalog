package tracks

import (
	"github.com/YugaAI/MusicCatalog/internal/models/trackactivities"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UpsertTrackActivity(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetUint("userID")

	var request trackactivities.TrackActivitiesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.service.UpsertTrackActivity(ctx, userID, request); err != nil {
		c.JSON(500, gin.H{"error": "Failed to upsert track activity"})
		return
	}

	c.Status(200) // No Content
}
