package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/YugaAI/MusicCatalog/internal/configs"
	"github.com/YugaAI/MusicCatalog/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddlewere() gin.HandlerFunc {
	secretKey := configs.Get().Service.SecretKey
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")

		header = strings.TrimSpace(header)
		if header == "" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("missing token"))
			return
		}

		userID, username, err := jwt.ValidateToken(header, secretKey)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set("userID", userID)
		c.Set("username", username)
		c.Next()
	}
}

func AuthRefreshMiddlewere() gin.HandlerFunc {
	secretKey := configs.Get().Service.SecretKey
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")

		header = strings.TrimSpace(header)
		if header == "" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("missing token"))
			return
		}

		userID, username, err := jwt.ValidateToken(header, secretKey)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set("userID", userID)
		c.Set("username", username)
		c.Next()
	}
}
