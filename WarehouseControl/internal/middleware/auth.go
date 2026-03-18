package middleware

import (
	"l3/WarehouseControl/internal/auth"
	"l3/WarehouseControl/internal/models"
	"net/http"
	"strings"

	"github.com/wb-go/wbf/ginext"
)

type AuthMiddleware struct {
	jwt *auth.JWT
}

func NewAuthMiddleware(jwt *auth.JWT) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt,
	}
}

func (m *AuthMiddleware) Auth() ginext.HandlerFunc {
	return func(c *ginext.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ginext.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, ginext.H{
				"error": "invalid authorization header",
			})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		claims, err := m.jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ginext.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		currentUser := &models.CurrentUser{
			UserID:   claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		ctx := SetCurrentUser(c.Request.Context(), currentUser)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
