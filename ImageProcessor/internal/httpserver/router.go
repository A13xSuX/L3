package httpserver

import (
	"l3/ImageProcessor/internal/httpserver/handlers"
	"l3/ImageProcessor/repo"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
)

func NewRouter(imagesRepo *repo.ImagesRepo, producer *kafka.Producer, strategy retry.Strategy, templatesGlob string) *ginext.Engine {
	r := ginext.New("release")
	r.LoadHTMLGlob(templatesGlob)
	r.Use(ginext.Logger(), ginext.Recovery()) //logs of gin
	r.GET("/healthz", handlers.Healthz(imagesRepo))
	r.MaxMultipartMemory = 8 << 20
	r.POST("/upload", handlers.UploadMultipleFile(imagesRepo, producer, strategy))
	r.GET("/", handlers.Template())
	r.GET("/image/:id/meta", handlers.GetImageMeta(imagesRepo))
	r.GET("/image/:id", handlers.GetImage(imagesRepo))
	r.DELETE("/image/:id", handlers.DeleteImage(imagesRepo))
	return r
}
