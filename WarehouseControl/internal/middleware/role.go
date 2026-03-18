package middleware

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

func RequireRoles(roles ...string) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		currentUser, ok := GetCurrentUser(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, ginext.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
		for _, role := range roles {
			if currentUser.Role == role {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, ginext.H{
			"error": "forbidden",
		})
		c.Abort()
	}
}
