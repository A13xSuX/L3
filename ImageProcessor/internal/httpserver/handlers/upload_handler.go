package handlers

import (
	"l3/ImageProcessor/repo"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	kafkawbf "github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
)

// сохраняет на диск и запись в БД + шлет msg в kafka
func UploadMultipleFile(imagesRepo *repo.ImagesRepo, producer *kafkawbf.Producer, strategy retry.Strategy) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		files := form.File["images"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, ginext.H{"status": "error", "message": "no files in field 'images'"})
			return
		}
		ids := []string{}
		for _, file := range files {
			id := uuid.New()
			dst := "./data/original/" + id.String() + filepath.Ext(file.Filename)

			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(http.StatusInternalServerError, ginext.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}

			if err := imagesRepo.Create(c.Request.Context(), id, dst); err != nil {
				c.JSON(http.StatusInternalServerError, ginext.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}

			key := []byte("image")
			value := []byte(id.String())

			if err := producer.SendWithRetry(c.Request.Context(), strategy, key, value); err != nil {
				c.JSON(http.StatusInternalServerError, ginext.H{
					"status":  "error",
					"message": "failed to send to kafka: " + err.Error(),
				})
				return
			}
			ids = append(ids, id.String())
		}
		c.JSON(http.StatusCreated, ginext.H{
			"ids": ids,
		})
	}
}
