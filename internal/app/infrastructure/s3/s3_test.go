package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/test/mock/mock_s3"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/test/testutil"
)

func TestHandler_Save(t *testing.T) {
	ctx := context.Background()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockAPI := mock_s3.NewMockS3API(mockController)
	handler := new(Handler)
	handler.s3 = mockAPI

	object := []byte("object")
	bucketName := "bucketName"
	key := "key"

	testCases := []struct {
		name         string
		putObjectErr error
		expectedErr  error
	}{
		{
			name: "OK",
		},
		{
			name:         "NG",
			putObjectErr: testutil.ErrSome,
			expectedErr:  fmt.Errorf("failed to PutObject: %w", testutil.ErrSome),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAPI.EXPECT().PutObject(ctx, &s3.PutObjectInput{
				Body:   bytes.NewReader(object),
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
			}).Return(nil, tc.putObjectErr)

			err := handler.Save(ctx, object, key, bucketName)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestHandler_Find(t *testing.T) {
	ctx := context.Background()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockAPI := mock_s3.NewMockS3API(mockController)
	handler := new(Handler)
	handler.s3 = mockAPI

	object := []byte("object")
	bucketName := "bucketName"
	key := "key"

	testCases := []struct {
		name         string
		getObjectRet *s3.GetObjectOutput
		getObjectErr error
		expected     []byte
		expectedErr  error
	}{
		{
			name: "OK",
			getObjectRet: &s3.GetObjectOutput{
				Body: io.NopCloser(bytes.NewReader(object)),
			},
			expected: object,
		},
		{
			name:         "NG: failed to GetObject",
			getObjectErr: testutil.ErrSome,
			expectedErr:  fmt.Errorf("failed to GetObject: %w", testutil.ErrSome),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAPI.EXPECT().GetObject(ctx, &s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
			}).Return(tc.getObjectRet, tc.getObjectErr)

			actual, err := handler.Find(ctx, key, bucketName)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
