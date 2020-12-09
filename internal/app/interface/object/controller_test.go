package object

import (
	"context"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/domain"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/infrastructure/env"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/test/mock/mock_interactor"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/test/testutil"
)

func TestController_CreateThumbnail(t *testing.T) {
	ctx := context.Background()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockInteractor := mock_interactor.NewMockImageInteractor(mockController)
	controller := NewController(mockInteractor, env.Object.BucketNameThumbnail)

	testCases := []struct {
		name               string
		key                string
		bucketName         string
		createThumbnailErr error
		expectedErr        error
	}{
		{
			name:       "OK: jpg",
			key:        keyImage("test", domain.ImageFormatJPEG),
			bucketName: env.Object.BucketNameOriginal,
		},
		{
			name:       "OK: jpeg",
			key:        "test.jpeg",
			bucketName: env.Object.BucketNameOriginal,
		},
		{
			name:       "OK: gif",
			key:        "test.gif",
			bucketName: env.Object.BucketNameOriginal,
		},
		{
			name:        "NG: src bucket == dst bucket",
			bucketName:  env.Object.BucketNameThumbnail,
			expectedErr: xerrors.New("src bucket and dst bucket is the same"),
		},
		{
			name:        "NG: invalid key format",
			key:         "test",
			bucketName:  env.Object.BucketNameOriginal,
			expectedErr: xerrors.Errorf("failed to extractNameAndFormat: %w", xerrors.New("misspecified separator")),
		},
		{
			name:        "NG: invalid image format",
			key:         "test.jpg.gz",
			bucketName:  env.Object.BucketNameOriginal,
			expectedErr: xerrors.Errorf("failed to extractNameAndFormat: %w", xerrors.Errorf("failed to UnmarshalText: %w", xerrors.New("invalid ImageFormat"))),
		},
		{
			name:               "NG: interactor.CreateThumbnailErr",
			key:                "test.jpg",
			bucketName:         env.Object.BucketNameOriginal,
			createThumbnailErr: testutil.ErrSome,
			expectedErr:        testutil.ErrSome,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.createThumbnailErr != nil || tc.expectedErr == nil {
				splitKey := strings.Split(tc.key, ".")
				require.Len(t, splitKey, 2)
				var format domain.ImageFormat
				require.NoError(t, format.UnmarshalText([]byte(splitKey[1])))
				mockInteractor.EXPECT().CreateThumbnail(ctx, splitKey[0], format).Return(tc.createThumbnailErr)
			}

			err := controller.CreateThumbnail(ctx, tc.key, tc.bucketName)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}
			assert.NoError(t, err)
		})
	}
}
