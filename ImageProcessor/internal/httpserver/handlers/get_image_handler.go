package handlers

import (
	"errors"
	"l3/ImageProcessor/repo"
	"net/http"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
)

func GetImage(imagesRepo *repo.ImagesRepo) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"status": "error", "message": "invalid uuid"})
			return
		}
		img, err := imagesRepo.Get(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				c.JSON(http.StatusNotFound, ginext.H{"status": "error", "message": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, ginext.H{"status": "error", "message": err.Error()})
			return
		}
		if img.Status != "ready" || img.ProcessedPath == nil {
			c.JSON(http.StatusConflict, ginext.H{
				"status": "error",
				"msg":    "processed is not finished",
			})
			return
		}
		c.File(*img.ProcessedPath)
	}
}
