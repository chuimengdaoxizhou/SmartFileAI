package middleware

import (
	"RAG/services"
	"RAG/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			utils.ErrorResponse(c, 401, "Authorization token required")
			c.Abort()
			return
		}

		userID, err := services.ValidateJWT(token)
		if err != nil {
			utils.ErrorResponse(c, 401, "Invalid token")
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
