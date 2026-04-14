package domain

import (
	"bytes"
	"errors"
	"fmt"
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
			expectedErr: errors.New("rate must be in (0, 1)"),
		},
		{
			name:        "NG: rate one",
			rate:        1,
			expectedErr: errors.New("rate must be in (0, 1)"),
		},
		{
			name: "NG: img is already thumbnail",
			rate: 0.5,
			img: &Image{
				IsThumbnail: true,
			},
			expectedErr: errors.New("image is already thumbnail"),
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
			expectedErr: errors.New("misspecified image"),
		},
		{
			name: "NG: initial format",
			img: &Image{
				Format: ImageFormat(0),
				Image:  img,
			},
			expectedErr: errors.New("not supported image format"),
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
			expectedErr: fmt.Errorf("failed to Decode: %w", errors.New("image: unknown format")),
		},
		{
			name:        "NG: file size exceeds maximum",
			b:           make([]byte, MaxImageFileSize+1),
			expectedErr: fmt.Errorf("file size %d exceeds maximum %d bytes", MaxImageFileSize+1, MaxImageFileSize),
		},
		{
			name: "NG: image dimensions exceed maximum",
			b: func() []byte {
				buf := new(bytes.Buffer)
				if err := png.Encode(buf, graffiti.SolidImage(image.Rect(0, 0, MaxImageWidth+1, 1), image.Black)); err != nil {
					t.Fatal(err)
				}
				return buf.Bytes()
			}(),
			expectedErr: fmt.Errorf("image dimensions %dx%d exceed maximum %dx%d", MaxImageWidth+1, 1, MaxImageWidth, MaxImageHeight),
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
