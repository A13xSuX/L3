package handlers

import (
	"errors"
	"l3/ImageProcessor/repo"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
)

func DeleteImage(imagesRepo *repo.ImagesRepo) ginext.HandlerFunc {
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
		paths := []string{img.OriginalPath}
		if img.ProcessedPath != nil {
			paths = append(paths, *img.ProcessedPath)
		}
		if img.ThumbPath != nil {
			paths = append(paths, *img.ThumbPath)
		}
		for _, p := range paths {
			if p == "" {
				continue
			}
			if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
				c.JSON(http.StatusInternalServerError, ginext.H{
					"status":  "error",
					"message": "failed to remove file: " + err.Error(),
				})
				return
			}
		}
		if err := imagesRepo.Delete(c.Request.Context(), id); err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				// запись могла исчезнуть между Get и Delete
				c.JSON(http.StatusNotFound, ginext.H{"status": "error", "message": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, ginext.H{"status": "error", "message": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
