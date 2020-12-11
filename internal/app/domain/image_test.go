package domain

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/mububoki/graffiti"
	gif2 "github.com/mububoki/graffiti/gif"
	jpeg2 "github.com/mububoki/graffiti/jpeg"
	png2 "github.com/mububoki/graffiti/png"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

func TestImage_CreateThumbnail(t *testing.T) {
	testCases := []struct {
		name        string
		img         *Image
		rate        float64
		expectedErr error
	}{
		{
			name: "OK",
			img: &Image{
				Name:        "test",
				Format:      ImageFormatJPEG,
				IsThumbnail: false,
				Image:       graffiti.RandomImage(image.Rect(0, 0, 1024, 780)),
			},
			rate: 0.5,
		},
		{
			name:        "NG: rate zero",
			rate:        0,
			expectedErr: xerrors.New("rate must be in (0, 1)"),
		},
		{
			name:        "NG: rate one",
			rate:        1,
			expectedErr: xerrors.New("rate must be in (0, 1)"),
		},
		{
			name: "NG: img is already thumbnail",
			rate: 0.5,
			img: &Image{
				IsThumbnail: true,
			},
			expectedErr: xerrors.New("image is already thumbnail"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tmb, err := tc.img.CreateThumbnail(tc.rate)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.img.Name, tmb.Name)
			assert.Equal(t, tc.img.Format, tc.img.Format)
			assert.True(t, tmb.IsThumbnail)
			assert.Equal(t, int(float64(tc.img.Image.Bounds().Dx())*tc.rate), tmb.Image.Bounds().Dx())
			assert.Equal(t, int(float64(tc.img.Image.Bounds().Dy())*tc.rate), tmb.Image.Bounds().Dy())
		})
	}
}

func TestImage_Encode(t *testing.T) {
	img := graffiti.RandomImage(image.Rect(0, 0, 10, 10))
	bJPEG := new(bytes.Buffer)
	require.NoError(t, jpeg.Encode(bJPEG, img, &jpeg.Options{Quality: 100}))
	bGIF := new(bytes.Buffer)
	require.NoError(t, gif.Encode(bGIF, img, nil))
	bPNG := new(bytes.Buffer)
	require.NoError(t, png.Encode(bPNG, img))

	testCases := []struct {
		name        string
		img         *Image
		expected    []byte
		expectedErr error
	}{
		{
			name: "OK: jpeg",
			img: &Image{
				Format: ImageFormatJPEG,
				Image:  img,
			},
			expected: bJPEG.Bytes(),
		},
		{
			name: "OK: gif",
			img: &Image{
				Format: ImageFormatGIF,
				Image:  img,
			},
			expected: bGIF.Bytes(),
		},
		{
			name: "OK: png",
			img: &Image{
				Format: ImageFormatPNG,
				Image:  img,
			},
			expected: bPNG.Bytes(),
		},
		{
			name: "NG: nil image",
			img: &Image{
				Image: nil,
			},
			expectedErr: xerrors.New("misspecified image"),
		},
		{
			name: "NG: too large image",
			img: &Image{
				Format: ImageFormatJPEG,
				Image:  image.NewRGBA(image.Rect(0, 0, 1<<16, 1<<16)),
			},
			expectedErr: xerrors.Errorf("failed to jpeg.Encode: %w", xerrors.New("jpeg: image is too large to encode")),
		},
		{
			name: "NG: too large image",
			img: &Image{
				Format: ImageFormatGIF,
				Image:  image.NewRGBA(image.Rect(0, 0, 1<<16, 1<<16)),
			},
			expectedErr: xerrors.Errorf("failed to gif.Encode: %w", xerrors.New("gif: image is too large to encode")),
		},
		{
			name: "NG: invalid format png",
			img: &Image{
				Format: ImageFormatPNG,
				Image:  image.NewRGBA(image.Rect(0, 0, 1<<32, 1<<32)),
			},
			expectedErr: xerrors.Errorf("failed to png.Encode: %w", xerrors.New("png: invalid format: invalid image size: 4294967296x4294967296")),
		},
		{
			name: "NG: initial format",
			img: &Image{
				Format: ImageFormat(0),
				Image:  img,
			},
			expectedErr: xerrors.New("not supported image format"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.img.Encode()
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestDecodeImage(t *testing.T) {
	name := "test"
	bJPEG := new(bytes.Buffer)
	require.NoError(t, jpeg2.EncodeRandom(bJPEG, image.Rect(0, 0, 10, 10), &jpeg.Options{Quality: 100}))
	bGIF := new(bytes.Buffer)
	require.NoError(t, gif2.EncodeRandom(bGIF, image.Rect(0, 0, 10, 10), nil))
	bPNG := new(bytes.Buffer)
	require.NoError(t, png2.EncodeRandom(bPNG, image.Rect(0, 0, 10, 10)))

	imgIMGJPEG, _, err := image.Decode(bytes.NewReader(bJPEG.Bytes()))
	require.NoError(t, err)
	imgIMGGIF, _, err := image.Decode(bytes.NewReader(bGIF.Bytes()))
	require.NoError(t, err)
	imgIMGPNG, _, err := image.Decode(bytes.NewReader(bPNG.Bytes()))
	require.NoError(t, err)

	imgJPEG := &Image{
		Name:   name,
		Format: ImageFormatJPEG,
		Image:  imgIMGJPEG,
	}
	imgGIF := &Image{
		Name:   name,
		Format: ImageFormatGIF,
		Image:  imgIMGGIF,
	}
	imgPNG := &Image{
		Name:   name,
		Format: ImageFormatPNG,
		Image:  imgIMGPNG,
	}

	testCases := []struct {
		name        string
		b           []byte
		expected    *Image
		expectedErr error
	}{
		{
			name:     "OK: jpeg",
			b:        bJPEG.Bytes(),
			expected: imgJPEG,
		},
		{
			name:     "OK: gif",
			b:        bGIF.Bytes(),
			expected: imgGIF,
		},
		{
			name:     "OK: png",
			b:        bPNG.Bytes(),
			expected: imgPNG,
		},
		{
			name:        "NG: failed to Decode",
			b:           []byte("invalid image"),
			expectedErr: xerrors.Errorf("failed to Decode: %w", xerrors.New("image: unknown format")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			img, err := DecodeImage(tc.b, name, false)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, img)
		})
	}
}
