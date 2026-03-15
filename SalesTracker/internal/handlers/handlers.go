package handlers

import (
	"l3/SalesTracker/internal/models"
	"l3/SalesTracker/internal/repository"
	"net/http"
	"time"

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

func (h *SaleHandler) Analytics(c *ginext.Context) {
	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "from or to is empty",
		})
		return
	}
	category := c.Query("category")
	title := c.Query("title")
	fromTime, err := time.Parse("2006-01-02", from)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}
	toTime, err := time.Parse("2006-01-02", to) //2006-01-02T15:04:05Z07:00
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}
	if fromTime.After(toTime) {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "from must be before to",
		})
		return
	}
	analyticsFilters := &models.AnalyticsFilter{
		From:     fromTime,
		To:       toTime,
		Category: category,
		Title:    title,
	}
	resp, err := h.salesRepo.Analytics(c.Request.Context(), analyticsFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
