package handlers

import (
	"errors"
	"l3/WarehouseControl/internal/customErrs"
	"l3/WarehouseControl/internal/models"
	"l3/WarehouseControl/internal/service"
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

type LoginHandler struct {
	loginService *service.LoginService
}

func NewLoginHandler(loginService *service.LoginService) *LoginHandler {
	return &LoginHandler{
		loginService: loginService,
	}
}

func (h *LoginHandler) Login(c *ginext.Context) {
	var req models.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}
	resp, err := h.loginService.Login(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, customErrs.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, ginext.H{
				"error": "invalid credentials",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "internal server error",
		})
		return

	}
	c.JSON(http.StatusOK, resp)
}
