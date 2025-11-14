package middleware

import (
	"net/http"
	"strings"

	"github.com/dmehra2102/StudySync/pkg/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtAuth *auth.JWTAuth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})

			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})

			ctx.Abort()
			return
		}

		token := parts[1]
		claims, err := jwtAuth.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", claims.UserID)
		ctx.Set("userEmail", claims.Email)
		ctx.Set("userRole", claims.Role)
		ctx.Next()
	}
}
