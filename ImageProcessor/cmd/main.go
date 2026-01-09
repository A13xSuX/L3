package main

import (
	"context"
	"fmt"
	"io"
	"l3/ImageProcessor/internal/appcfg"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	path := filepath.Join("..", "web", "*")
	baseDir := filepath.Join(".", "out")
	zlog.InitConsole()
	_ = zlog.SetLevel(cfg.LoggerConfig.LogLevel)
	r := ginext.New("release")
	r.LoadHTMLGlob(path)
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
	r.POST("/upload", uploadMultipleFile(baseDir))
	r.GET("/", template())
	r.GET("/image/:id", getImage(baseDir))
	r.DELETE("/image/:id", deleteImage(baseDir))

	//create fs
	err = os.MkdirAll(baseDir, os.ModePerm)
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
func uploadMultipleFile(baseDir string) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("Не удалось получить файлы с формы")
			return
		}
		files := form.File["images"]
		ids := []string{}
		for _, file := range files {
			fileExt := filepath.Ext(file.Filename)
			switch fileExt {
			case ".png":
				func() {
					originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), fileExt)

					now := time.Now()
					filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "_") + "" + fmt.Sprintf("%v", now.Unix()) + fileExt
					ids = append(ids, filename)
					out, err := os.Create(filepath.Join(baseDir, filename))
					if err != nil {
						zlog.Logger.Error().Msg("Не удалось создать файл")
						return
					}
					defer out.Close()

					readerFile, err := file.Open()
					if err != nil {
						c.JSON(http.StatusInternalServerError, ginext.H{"status": "error", "message": "open file is crashed"})
						return
					}
					defer readerFile.Close()
					_, err = io.Copy(out, readerFile)
					if err != nil {
						zlog.Logger.Error().Err(err).Msg("Не удалось сохранить файл")
						return
					}
					zlog.Logger.Info().Msg("Файлы успешно скачаны")
				}()

			default:
				c.JSON(http.StatusBadRequest, ginext.H{"status": "error", "message": "invalid format file, only PNG"})
				return
			}
		}
		c.JSON(http.StatusCreated, ginext.H{
			"ids": ids,
		})
	}
}

func getImage(baseDir string) ginext.HandlerFunc { //потом заменим на postgreSQL
	return func(c *ginext.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, ginext.H{"status": "error", "message": "invalid id"})
			return
		}
		if strings.Contains(idStr, "/") || strings.Contains(idStr, "\\") || strings.Contains(idStr, "..") {
			c.JSON(http.StatusBadRequest, ginext.H{"status": "error", "message": "invalid id"})
			return
		}
		pathOut := filepath.Join(baseDir, idStr)
		_, err := os.Stat(pathOut)
		if err != nil {
			if os.IsNotExist(err) {
				c.JSON(http.StatusNotFound, ginext.H{"status": "error", "message": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, ginext.H{"status": "error", "message": err.Error()})
			return
		}
		c.File(pathOut)
	}
}

func deleteImage(baseDir string) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, ginext.H{"status": "error", "message": "invalid id"})
			return
		}
		if strings.Contains(idStr, "/") || strings.Contains(idStr, "\\") || strings.Contains(idStr, "..") {
			c.JSON(http.StatusBadRequest, ginext.H{"status": "error", "message": "invalid id"})
			return
		}
		deletePath := filepath.Join(baseDir, idStr)
		err := os.Remove(deletePath)
		if err != nil {
			if os.IsNotExist(err) {
				c.JSON(http.StatusNotFound, ginext.H{"status": "error", "message": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, ginext.H{"status": "error", "error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
