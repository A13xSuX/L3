package handlers

import (
	"l3/WarehouseControl/internal/middleware"
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Me(c *ginext.Context) {
	user, ok := middleware.GetCurrentUser(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, ginext.H{
			"error": "unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"user_id":  user.UserID,
		"username": user.Username,
		"role":     user.Role,
	})
}
