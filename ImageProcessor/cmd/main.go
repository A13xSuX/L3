package main

import (
	"context"
	"errors"
	"l3/ImageProcessor/internal/appcfg"
	"l3/ImageProcessor/repo"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	cfg, err := appcfg.NewAppConfig()
	if err != nil {
		zlog.Logger.Error().Err(err)
		return
	}

	opts := &dbpg.Options{MaxOpenConns: cfg.PostgresConfig.MaxOpenConns,
		MaxIdleConns:    cfg.PostgresConfig.MaxIdleConns,
		ConnMaxLifetime: cfg.PostgresConfig.ConnMaxLifeTime,
	}
	db, err := dbpg.New(cfg.PostgresConfig.MasterDSN, cfg.PostgresConfig.SlaveDSN, opts)
	if err != nil {
		zlog.Logger.Error().Err(err)
		return
	}
	imagesRepo := repo.NewImagesRepo(db)

	templatesGlob := filepath.Join("..", "web", "*")
	//baseDir := filepath.Join(".", "out")
	zlog.InitConsole()
	_ = zlog.SetLevel(cfg.LoggerConfig.LogLevel)
	r := ginext.New("release")
	r.LoadHTMLGlob(templatesGlob)
	r.Use(ginext.Logger(), ginext.Recovery()) //logs of gin
	r.GET("/healthz", func() ginext.HandlerFunc {
		return func(c *ginext.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var res string
			err := db.QueryRowContext(ctx, "SELECT 1").Scan(&res)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ginext.H{
					"status": "error",
					"error":  err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, ginext.H{"status": "ok"})
		}
	}())
	r.MaxMultipartMemory = 8 << 20
	r.POST("/upload", uploadMultipleFile(imagesRepo))
	r.GET("/", template())
	r.GET("/image/:id", getImage(imagesRepo))
	r.DELETE("/image/:id", deleteImage(imagesRepo))

	//create fs
	//err = os.MkdirAll(baseDir, os.ModePerm)
	//if err != nil {
	//	zlog.Logger.Error().Err(err).Msg("Filesystem is not created")
	//	return
	//}
	err = os.MkdirAll("./data/original", os.ModePerm)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Filesystem is not created")
		return
	}

	if err := r.Run(cfg.ServerConfig.Addr); err != nil {
		zlog.Logger.Error().Msg("Server is down")
		return
	}
}
func template() ginext.HandlerFunc {
	return func(c *ginext.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	}
}
func uploadMultipleFile(imagesRepo *repo.ImagesRepo) ginext.HandlerFunc {
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

			ids = append(ids, id.String())
		}
		c.JSON(http.StatusCreated, ginext.H{
			"ids": ids,
		})
	}
}

func getImage(imagesRepo *repo.ImagesRepo) ginext.HandlerFunc { //потом заменим на postgreSQL
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

func deleteImage(imagesRepo *repo.ImagesRepo) ginext.HandlerFunc {
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
