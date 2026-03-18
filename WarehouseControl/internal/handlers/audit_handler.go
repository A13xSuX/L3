package handlers

import (
	"l3/WarehouseControl/internal/service"
	"net/http"
	"strconv"

	"github.com/wb-go/wbf/ginext"
)

type AuditHandler struct {
	auditService *service.AuditService
}

func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
	}
}

func (h *AuditHandler) GetByItemID(c *ginext.Context) {
	itemIDStr := c.Param("id")
	if itemIDStr == "" {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "id is empty",
		})
		return
	}
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid id",
		})
		return
	}
	if itemID <= 0 {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid id",
		})
		return
	}
	audit, err := h.auditService.GetByItemID(c.Request.Context(), int64(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "internal server",
		})
		return
	}
	c.JSON(http.StatusOK, audit)
}
