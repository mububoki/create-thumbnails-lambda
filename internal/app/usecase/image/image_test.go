package image

import (
	"context"
	"image"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mububoki/graffiti"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/domain"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/infrastructure/env"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/test/mock/mockport"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/test/testutil"
)

func TestInteractor_CreateThumbnail(t *testing.T) {
	ctx := context.Background()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockRepo := mockport.NewMockImageRepository(mockController)
	interactor := NewInteractor(mockRepo, env.Image.Rate)

	img := graffiti.RandomImage(image.Rect(0, 0, 10, 10))

	testCases := []struct {
		name        string
		imageName   string
		format      domain.ImageFormat
		img         image.Image
		findErr     error
		saveErr     error
		expectedErr error
	}{
		{
			name:      "OK: jpeg",
			imageName: "test",
			format:    domain.ImageFormatJPEG,
			img:       img,
		},
		{
			name:      "OK: gif",
			imageName: "test",
			format:    domain.ImageFormatGIF,
			img:       img,
		},
		{
			name:      "OK: png",
			imageName: "test",
			format:    domain.ImageFormatPNG,
			img:       img,
		},
		{
			name:        "NG: failed to Find",
			findErr:     testutil.ErrSome,
			expectedErr: xerrors.Errorf("failed to Find: %w", testutil.ErrSome),
		},
		{
			name:        "NG: failed to Save",
			img:         img,
			saveErr:     testutil.ErrSome,
			expectedErr: xerrors.Errorf("failed to Save: %w", testutil.ErrSome),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			func() {
				src := &domain.Image{
					Name:        tc.imageName,
					Format:      tc.format,
					IsThumbnail: false,
					Image:       tc.img,
				}

				mockRepo.EXPECT().Find(ctx, tc.imageName, tc.format, false).Return(src, tc.findErr)
				if tc.findErr != nil {
					return
				}

				dst, err := src.CreateThumbnail(interactor.rate)
				require.NoError(t, err)

				mockRepo.EXPECT().Save(ctx, dst).Return(tc.saveErr)
				if tc.saveErr != nil {
					return
				}
			}()

			err := interactor.CreateThumbnail(ctx, tc.imageName, tc.format)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}
			assert.NoError(t, err)
		})
	}
}
