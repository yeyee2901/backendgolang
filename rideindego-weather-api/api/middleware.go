package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-www-form-urlencoded")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func SimpleAuthorizationMiddleware(secretKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// skip this middleware for swagger
		if strings.HasPrefix(ctx.Request.URL.Path, "/swagger/") {
			ctx.Next()
			return
		}

		clientAuth := ctx.Request.Header.Get("Authorization")
		if clientAuth != "Bearer "+secretKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized key",
			})
			return
		}

		ctx.Next()
	}
}
