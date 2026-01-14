package worker

import (
	"context"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"l3/ImageProcessor/repo"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/segmentio/kafka-go"
	kafkawbf "github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func ProcessImages(ctx context.Context, imagesRepo *repo.ImagesRepo, consumer *kafkawbf.Consumer, strategy retry.Strategy) {
	msgCh := make(chan kafka.Message)

	consumer.StartConsuming(ctx, msgCh, strategy)
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgCh:
			if !ok {
				return
			}
			idStr := string(msg.Value)
			id, err := uuid.Parse(idStr)
			if err != nil {
				zlog.Logger.Error().Str("id", idStr).Err(err).Msg("invalid uuid from kafka")
				continue
			}
			img, err := imagesRepo.Get(ctx, id)
			if err != nil {
				if errors.Is(err, repo.ErrNotFound) {
					zlog.Logger.Warn().Str("id", idStr).Msg("image not found for kafka message")
					continue
				}
				zlog.Logger.Error().Err(err).Msg("get image failed")
				continue
			}
			err = imagesRepo.MarkProcessing(ctx, img.ID)
			if err != nil {
				zlog.Logger.Error().Err(err)
				continue
			}

			srcImg, _, err := loadImage(img.OriginalPath)
			if err != nil {
				_ = imagesRepo.MarkFailed(ctx, img.ID, err.Error())
				continue
			}
			processedImg := resize.Resize(1024, 0, srcImg, resize.Lanczos3)
			thumbImg := resize.Thumbnail(256, 256, srcImg, resize.Lanczos3)

			wmText := "ImageProcessor"
			processedImg = addWatermark(processedImg, wmText)

			ext := filepath.Ext(img.OriginalPath)
			processedDir := "./data/processed"
			thumbsDir := "./data/thumbs"
			processedPath := filepath.Join(processedDir, img.ID.String()+ext)
			thumbPath := filepath.Join(thumbsDir, img.ID.String()+ext)

			if err := saveImage(processedPath, processedImg); err != nil {
				_ = imagesRepo.MarkFailed(ctx, img.ID, err.Error())
				continue
			}
			if err := saveImage(thumbPath, thumbImg); err != nil {
				_ = imagesRepo.MarkFailed(ctx, img.ID, err.Error())
				continue
			}
			if err := imagesRepo.MarkReady(ctx, img.ID, processedPath, thumbPath); err != nil {
				zlog.Logger.Error().Err(err)
			}
		}
	}
}

func loadImage(path string) (image.Image, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	img, format, err := image.Decode(f)
	if err != nil {
		return nil, "", err
	}
	return img, format, nil
}

func saveImage(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	opts := &jpeg.Options{Quality: 85}
	return jpeg.Encode(f, img, opts)
}

func addWatermark(src image.Image, text string) image.Image {
	rgba := image.NewRGBA(src.Bounds())
	draw.Draw(rgba, rgba.Bounds(), src, image.Point{}, draw.Src)

	bounds := rgba.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	boxW := w / 3
	boxH := h / 12
	if boxH < 30 {
		boxH = 30
	}

	margin := 20
	boxX1 := w - boxW - margin
	boxY1 := h - boxH - margin
	boxX2 := w - margin
	boxY2 := h - margin

	boxColor := color.RGBA{0, 0, 0, 140}
	draw.Draw(rgba, image.Rect(boxX1, boxY1, boxX2, boxY2), &image.Uniform{boxColor}, image.Point{}, draw.Src)

	textColor := image.NewUniform(color.RGBA{255, 255, 255, 240})

	face := basicfont.Face7x13

	textPxW := len(text) * 7
	textX := boxX1 + (boxW-textPxW)/2
	if textX < boxX1+5 {
		textX = boxX1 + 5
	}
	textY := boxY1 + boxH/2 + 5

	d := &font.Drawer{
		Dst:  rgba,
		Src:  textColor,
		Face: face,
		Dot:  fixed.P(textX, textY),
	}
	d.DrawString(text)

	return rgba
}
