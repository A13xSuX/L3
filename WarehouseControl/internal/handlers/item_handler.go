package handlers

import (
	"errors"
	"l3/WarehouseControl/internal/customErrs"
	"l3/WarehouseControl/internal/models"
	"l3/WarehouseControl/internal/service"
	"net/http"
	"strconv"

	"github.com/wb-go/wbf/ginext"
)

type ItemHandler struct {
	itemService *service.ItemService
}

func NewItemHandler(itemService *service.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

func (h *ItemHandler) Create(c *ginext.Context) {
	var itemReq models.CreateItemRequest
	err := c.ShouldBindJSON(&itemReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}
	itemResponse, err := h.itemService.Create(c.Request.Context(), &itemReq)
	if err != nil {
		if errors.Is(err, customErrs.ErrTitleEmpty) {
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": "title is empty",
			})
			return
		}
		if errors.Is(err, customErrs.ErrSkuEmpty) {
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": "sku is empty",
			})
			return
		}
		if errors.Is(err, customErrs.ErrQuantityNotPositive) {
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": "quantity cannot be negative",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "internal error",
		})
		return
	}
	c.JSON(http.StatusCreated, itemResponse)
}

func (h *ItemHandler) GetAll(c *ginext.Context) {
	itemsResp, err := h.itemService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "internal error",
		})
		return
	}
	c.JSON(http.StatusOK, itemsResp)
}

func (h *ItemHandler) Update(c *ginext.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "id is empty",
		})
		return
	}
	var newItem models.UpdateItemRequest
	err := c.ShouldBindJSON(&newItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid id",
		})
		return
	}
	err = h.itemService.Update(c.Request.Context(), int64(id), &newItem)
	if err != nil {
		if errors.Is(err, customErrs.ErrTitleEmpty) {
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": "title is empty",
			})
			return
		}
		if errors.Is(err, customErrs.ErrSkuEmpty) {
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": "sku is empty",
			})
			return
		}
		if errors.Is(err, customErrs.ErrQuantityNotPositive) {
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": "quantity cannot be negative",
			})
			return
		}
		if errors.Is(err, customErrs.ErrNotFoundID) {
			c.JSON(http.StatusNotFound, ginext.H{
				"error": "id not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "internal error",
		})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ItemHandler) Delete(c *ginext.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "id is empty",
		})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid id",
		})
		return
	}
	err = h.itemService.Delete(c.Request.Context(), int64(id))
	if err != nil {
		if errors.Is(err, customErrs.ErrNotFoundID) {
			c.JSON(http.StatusNotFound, ginext.H{
				"error": "id not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "internal server",
		})
		return
	}
	c.Status(http.StatusNoContent)
}
