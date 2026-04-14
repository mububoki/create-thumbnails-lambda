package domain

import (
	"bytes"
	"errors"
	"fmt"
	"image"

	"golang.org/x/image/draw"
)

const (
	MaxImageWidth    = 10000
	MaxImageHeight   = 10000
	MaxImageFileSize = 50 * 1024 * 1024 // 50MB
)

type Image struct {
	Name        string
	Format      ImageFormat
	IsThumbnail bool
	Image       image.Image
}

func (i *Image) CreateThumbnail(rate float64) (*Image, error) {
	if rate <= 0 || rate >= 1 {
		return nil, errors.New("rate must be in (0, 1)")
	}

	if i.IsThumbnail {
		return nil, errors.New("image is already thumbnail")
	}

	rectangle := i.Image.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, int(float64(rectangle.Dx())*rate), int(float64(rectangle.Dy())*rate)))
	draw.CatmullRom.Scale(dst, dst.Bounds(), i.Image, rectangle, draw.Over, nil)

	return &Image{
		Name:        i.Name,
		Format:      i.Format,
		IsThumbnail: true,
		Image:       dst,
	}, nil
}

func (i *Image) Encode() ([]byte, error) {
	if i.Image == nil {
		return nil, errors.New("misspecified image")
	}

	return i.Format.encode(i.Image)
}

func DecodeImage(src []byte, name string, isThumbnail bool) (*Image, error) {
	if len(src) > MaxImageFileSize {
		return nil, fmt.Errorf("file size %d exceeds maximum %d bytes", len(src), MaxImageFileSize)
	}

	img, f, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("failed to Decode: %w", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() > MaxImageWidth || bounds.Dy() > MaxImageHeight {
		return nil, fmt.Errorf("image dimensions %dx%d exceed maximum %dx%d", bounds.Dx(), bounds.Dy(), MaxImageWidth, MaxImageHeight)
	}

	var format ImageFormat
	if err := format.UnmarshalText([]byte(f)); err != nil {
		return nil, fmt.Errorf("failed to UnmarshalText: %w", err)
	}

	return &Image{
		Name:        name,
		Format:      format,
		IsThumbnail: isThumbnail,
		Image:       img,
	}, nil
}
