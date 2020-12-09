package domain

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"

	"golang.org/x/xerrors"
)

type ImageFormat uint8

const (
	ImageFormatJPEG = iota + 1
	ImageFormatGIF
)

func (f ImageFormat) String() string {
	switch f {
	case ImageFormatJPEG:
		return "jpg"
	case ImageFormatGIF:
		return "gif"
	}

	return ""
}

func (f *ImageFormat) UnmarshalText(text []byte) error {
	switch string(text) {
	case "jpg", "jpeg":
		*f = ImageFormatJPEG
	case "gif":
		*f = ImageFormatGIF
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
	default:
		return nil, xerrors.New("not supported image format")
	}

	return b.Bytes(), nil
}
