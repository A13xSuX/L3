package handlers

import (
	"errors"
	"l3/ImageProcessor/repo"
	"net/http"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
)

func GetImageMeta(imagesRepo *repo.ImagesRepo) ginext.HandlerFunc { //потом заменим на postgreSQL
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
		c.JSON(http.StatusOK, ginext.H{
			"id":             img.ID.String(),
			"status":         img.Status,
			"original_path":  img.OriginalPath,
			"processed_path": img.ProcessedPath,
			"thumb_path":     img.ThumbPath,
			"error":          img.Error,
			"created_at":     img.CreatedAt,
			"updated_at":     img.UpdatedAt,
		})
	}
}
