package domain

import (
	"bytes"
	"image"

	"golang.org/x/image/draw"
	"golang.org/x/xerrors"
)

type Image struct {
	Name        string
	Format      ImageFormat
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
	return i.Format.Encode(i.Image)
}

func DecodeImage(src []byte, name string, isThumbnail bool) (*Image, error) {
	img, f, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return nil, xerrors.Errorf("faield to Decode: %w", err)
	}

	var format ImageFormat
	if err := format.UnmarshalText([]byte(f)); err != nil {
		return nil, xerrors.Errorf("failed to UnmarshalText: %w", err)
	}

	return &Image{
		Name:        name,
		Format:      format,
		IsThumbnail: isThumbnail,
		Image:       img,
	}, nil
}
