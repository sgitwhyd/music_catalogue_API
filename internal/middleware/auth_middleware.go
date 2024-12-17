package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/music-catalogue/internal/configs"
	"github.com/sgitwhyd/music-catalogue/pkg/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	secretKey := configs.Get().SecretJWT

	return func(ctx *gin.Context) {
		header :=  ctx.Request.Header.Get("Authorization")
		header = strings.TrimSpace(header)
		if header == "" {
			log.Error().Msg("Unauthorize request")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token not provided",
			})
			return 
		}

		userID, username, err := jwt.ValidateToken(header, secretKey)
		if err != nil {
			log.Error().Msg("Token Invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("username", username)
		ctx.Next()
	}
}

func AuthRefreshMiddleware() gin.HandlerFunc {
	secretKey := configs.Get().SecretJWT

	return func(ctx *gin.Context) {
		header :=  ctx.Request.Header.Get("Authorization")
		header = strings.TrimSpace(header)
		if header == "" {
			log.Error().Msg("Unauthorize request")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token not provided",
			})
			return 
		}

		userID, username, err := jwt.ValidateTokenWithoutExpiry(header, secretKey)
		if err != nil {
			log.Error().Msg("Token Invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("username", username)
		ctx.Next()
	}
}