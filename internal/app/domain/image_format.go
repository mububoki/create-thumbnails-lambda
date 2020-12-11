package domain

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"golang.org/x/xerrors"
)

type ImageFormat uint8

const (
	ImageFormatJPEG = iota + 1
	ImageFormatGIF
	ImageFormatPNG
)

func (f ImageFormat) String() string {
	switch f {
	case ImageFormatJPEG:
		return "jpg"
	case ImageFormatGIF:
		return "gif"
	case ImageFormatPNG:
		return "png"
	}

	return ""
}

func (f *ImageFormat) UnmarshalText(text []byte) error {
	switch string(text) {
	case "jpg", "jpeg":
		*f = ImageFormatJPEG
	case "gif":
		*f = ImageFormatGIF
	case "png":
		*f = ImageFormatPNG
	default:
		return xerrors.New("invalid ImageFormat")
	}

	return nil
}

func (f ImageFormat) encode(img image.Image) ([]byte, error) {
	b := new(bytes.Buffer)

	switch f {
	case ImageFormatJPEG:
		if err := jpeg.Encode(b, img, &jpeg.Options{Quality: 100}); err != nil {
			return nil, xerrors.Errorf("failed to jpeg.Encode: %w", err)
		}
	case ImageFormatGIF:
		if err := gif.Encode(b, img, nil); err != nil {
			return nil, xerrors.Errorf("failed to gif.Encode: %w", err)
		}
	case ImageFormatPNG:
		if err := png.Encode(b, img); err != nil {
			return nil, xerrors.Errorf("failed to png.Encode: %w", err)
		}
	default:
		return nil, xerrors.New("not supported image format")
	}

	return b.Bytes(), nil
}
