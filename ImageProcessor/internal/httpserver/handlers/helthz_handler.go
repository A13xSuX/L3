package handlers

import (
	"context"
	"l3/ImageProcessor/repo"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
)

func Healthz(imagesRepo *repo.ImagesRepo) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := imagesRepo.Healthz(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ginext.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, ginext.H{"status": "ok"})
	}
}
