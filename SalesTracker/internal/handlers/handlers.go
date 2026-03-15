package handlers

import (
	"l3/SalesTracker/internal/models"
	"l3/SalesTracker/internal/repository"
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

type SaleHandler struct {
	salesRepo *repository.SalesRepo
}

func NewSaleHandler(salesRepo *repository.SalesRepo) *SaleHandler {
	return &SaleHandler{
		salesRepo: salesRepo,
	}
}

func (h *SaleHandler) Create(c *ginext.Context) {
	var sale models.Sale
	err := c.ShouldBindJSON(&sale)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}
	err = sale.Validation()
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}
	err = h.salesRepo.Create(c.Request.Context(), &sale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, sale)
}

func (h *SaleHandler) Update(c *ginext.Context) {
	var newSale models.Sale
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "id is empty",
		})
		return
	}
	err := c.ShouldBindJSON(&newSale)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}
	newSale.ID = id
	err = newSale.Validation()
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}
	err = h.salesRepo.Update(c.Request.Context(), &newSale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, newSale)
}

func (h *SaleHandler) Delete(c *ginext.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "id is empty",
		})
		return
	}
	err := h.salesRepo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ginext.H{
		"details": "item deleted",
	})
}

func (h *SaleHandler) GetByID(c *ginext.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "id is empty",
		})
		return
	}

	sale, err := h.salesRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": err.Error(),
		})
		return
	}
	if sale == nil {
		c.JSON(http.StatusNotFound, ginext.H{
			"error": "sale id is not found",
		})
		return
	}
	c.JSON(http.StatusOK, sale)
}

func (h *SaleHandler) GetAll(c *ginext.Context) {
	sales, err := h.salesRepo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, sales)
}
