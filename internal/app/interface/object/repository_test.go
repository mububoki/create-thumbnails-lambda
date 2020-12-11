package object

import (
	"bytes"
	"context"
	"image"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mububoki/graffiti"
	"github.com/mububoki/graffiti/gif"
	"github.com/mububoki/graffiti/jpeg"
	"github.com/mububoki/graffiti/png"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/domain"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/infrastructure/env"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/test/mock/mock_gateway"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/test/testutil"
)

func TestRepository_Save(t *testing.T) {
	ctx := context.Background()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockRepository := mock_gateway.NewMockObjectStorage(mockController)
	repository := NewRepository(mockRepository, env.Object.BucketNameOriginal, env.Object.BucketNameThumbnail)

	validIMG := graffiti.RandomImage(image.Rect(0, 0, 10, 10))

	testCases := []struct {
		name        string
		image       *domain.Image
		saveErr     error
		expectedErr error
	}{
		{
			name: "OK: JPEG and IsNotThumbnail",
			image: &domain.Image{
				Name:   "test",
				Format: domain.ImageFormatJPEG,
				Image:  validIMG,
			},
		},
		{
			name: "OK: JPEG and IsThumbnail",
			image: &domain.Image{
				Name:        "test",
				Format:      domain.ImageFormatJPEG,
				IsThumbnail: true,
				Image:       validIMG,
			},
		},
		{
			name: "OK: GIF and IsNotThumbnail",
			image: &domain.Image{
				Name:   "test",
				Format: domain.ImageFormatGIF,
				Image:  validIMG,
			},
		},
		{
			name: "OK: GIF and IsThumbnail",
			image: &domain.Image{
				Name:        "test",
				Format:      domain.ImageFormatGIF,
				IsThumbnail: true,
				Image:       validIMG,
			},
		},
		{
			name: "OK: PNG and IsNotThumbnail",
			image: &domain.Image{
				Name:   "test",
				Format: domain.ImageFormatPNG,
				Image:  validIMG,
			},
		},
		{
			name: "OK: PNG and IsThumbnail",
			image: &domain.Image{
				Name:        "test",
				Format:      domain.ImageFormatPNG,
				IsThumbnail: true,
				Image:       validIMG,
			},
		},
		{
			name:        "NG: nil Image",
			image:       &domain.Image{},
			expectedErr: xerrors.Errorf("failed to Encode: %w", xerrors.New("misspecified image")),
		},
		{
			name: "NG: failed to Save",
			image: &domain.Image{
				Name:        "test",
				Format:      domain.ImageFormatJPEG,
				IsThumbnail: true,
				Image:       validIMG,
			},
			saveErr:     testutil.ErrSome,
			expectedErr: testutil.ErrSome,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			func() {
				object, err := tc.image.Encode()
				if err != nil {
					return
				}

				bucketName := repository.bucketNameOriginal
				if tc.image.IsThumbnail {
					bucketName = repository.bucketNameThumbnail
				}

				mockRepository.EXPECT().Save(ctx, object, tc.image.Name+"."+tc.image.Format.String(), bucketName).Return(tc.saveErr)
			}()

			err := repository.Save(ctx, tc.image)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestRepository_Find(t *testing.T) {
	ctx := context.Background()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockRepository := mock_gateway.NewMockObjectStorage(mockController)
	repository := NewRepository(mockRepository, env.Object.BucketNameOriginal, env.Object.BucketNameThumbnail)

	bJPEG := new(bytes.Buffer)
	require.NoError(t, jpeg.EncodeRandom(bJPEG, image.Rect(0, 0, 10, 10), nil))
	bGIF := new(bytes.Buffer)
	require.NoError(t, gif.EncodeRandom(bGIF, image.Rect(0, 0, 10, 10), nil))
	bPNG := new(bytes.Buffer)
	require.NoError(t, png.EncodeRandom(bPNG, image.Rect(0, 0, 10, 10)))

	testCases := []struct {
		name        string
		imageName   string
		format      domain.ImageFormat
		isThumbnail bool
		bytesIMG    []byte
		findErr     error
		expectedErr error
	}{
		{
			name:      "OK: jpeg and isNotThumbnail",
			imageName: "test",
			format:    domain.ImageFormatJPEG,
			bytesIMG:  bJPEG.Bytes(),
		},
		{
			name:        "OK: gif and isThumbnail",
			imageName:   "test",
			format:      domain.ImageFormatGIF,
			isThumbnail: true,
			bytesIMG:    bGIF.Bytes(),
		},
		{
			name:        "OK: png and isThumbnail",
			imageName:   "test",
			format:      domain.ImageFormatPNG,
			isThumbnail: true,
			bytesIMG:    bPNG.Bytes(),
		},
		{
			name:        "NG: failed to repository.Find",
			imageName:   "test",
			format:      domain.ImageFormatJPEG,
			bytesIMG:    bJPEG.Bytes(),
			findErr:     testutil.ErrSome,
			expectedErr: xerrors.Errorf("failed to Find: %w", testutil.ErrSome),
		},
		{
			name:        "NG: failed to DecodeImage",
			imageName:   "test",
			format:      domain.ImageFormatJPEG,
			bytesIMG:    []byte("invalid image"),
			expectedErr: xerrors.Errorf("failed to Decode: %w", xerrors.New("image: unknown format")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var expected *domain.Image
			func() {
				bucketName := repository.bucketNameOriginal
				if tc.isThumbnail {
					bucketName = repository.bucketNameThumbnail
				}

				mockRepository.EXPECT().Find(ctx, tc.imageName+"."+tc.format.String(), bucketName).Return(tc.bytesIMG, tc.findErr)
				if tc.findErr != nil {
					return
				}

				expected, _ = domain.DecodeImage(tc.bytesIMG, tc.imageName, tc.isThumbnail)
			}()

			actual, err := repository.Find(ctx, tc.imageName, tc.format, tc.isThumbnail)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}
