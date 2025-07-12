package tracks

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Search(c *gin.Context) {
	ctx := c.Request.Context()
	// Extract query parameter from the request context
	query := c.Query("query")
	// Extract pagination parameters from the request context
	pageSizeStr := c.DefaultQuery("pageSize", "10")  // Default page size is 10
	pageIndexStr := c.DefaultQuery("pageIndex", "1") // Default page index is 1

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10 // Default page size
	}
	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil {
		pageIndex = 1 // Default page index
	}

	userID := c.GetUint("userID")
	response, err := h.service.Search(ctx, query, pageSize, pageIndex, userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to search tracks"})
		return
	}
	c.JSON(200, response)
}
