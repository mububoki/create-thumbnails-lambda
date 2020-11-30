package domain

import (
	"image"

	"golang.org/x/image/draw"
	"golang.org/x/xerrors"
)

type Image struct {
	Name        string
	Format      string
	IsThumbnail bool
	Image       image.Image
}

func (i *Image) CreateThumbnail(rate float64) (*Image, error) {
	if rate < 0 || rate >= 1 {
		return nil, xerrors.New("rate must be in [0, 1)")
	}

	if i.IsThumbnail {
		return nil, xerrors.New("image is already thumbnail")
	}

	rectange := i.Image.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, int(float64(rectange.Dx())*rate), int(float64(rectange.Dy())*rate)))
	draw.CatmullRom.Scale(dst, dst.Bounds(), i.Image, rectange, draw.Over, nil)

	return &Image{
		Name:        i.Name,
		Format:      i.Format,
		IsThumbnail: true,
		Image:       dst,
	}, nil
}
