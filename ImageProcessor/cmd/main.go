package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

//прием изображения
//кладем в kafka
// в фоне обрабатываем файл(resize/watermark) параллельно

//можно создать папку(fs) туда будем класть фотки и оттуда можно будем нормально их доставать

//Пока все в одном файле

func main() {
	path := filepath.Join("..", "web", "*")
	baseDir := filepath.Join(".", "out")
	zlog.InitConsole()
	_ = zlog.SetLevel("debug")
	r := ginext.New("release")
	r.LoadHTMLGlob(path)
	r.Use(ginext.Logger(), ginext.Recovery()) //logs of gin
	r.POST("/upload", uploadMultipleFile(baseDir))
	r.GET("/", template())

	//create fs
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Filesystem is not created")
		return
	}

	if err := r.Run(":8080"); err != nil {
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
		filePaths := []string{}
		for _, file := range files {
			fileExt := filepath.Ext(file.Filename)
			switch fileExt {
			case ".png":
				originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), fileExt)
				now := time.Now()
				filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "_") + "" + fmt.Sprintf("%v", now.Unix()) + fileExt
				filePaths = append(filePaths, filename)
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
				_, err = io.Copy(out, readerFile)
				if err != nil {
					zlog.Logger.Error().Err(err).Msg("Не удалось сохранить файл")
					return
				}
				zlog.Logger.Info().Msg("Файлы успешно скачаны")

			default:
				c.JSON(http.StatusBadRequest, ginext.H{"status": "error", "message": "invalid format file, only PNG"})
				return
			}
		}
		c.JSON(http.StatusOK, ginext.H{
			"filepath": filePaths,
		})
	}
}
