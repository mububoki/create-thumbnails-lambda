package s3

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/test/mock/mock_s3"
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

	someErr := xerrors.New("some error")

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
			putObjectErr: someErr,
			expectedErr:  xerrors.Errorf("failed to PutObjectWithContext: %w", someErr),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAPI.EXPECT().PutObjectWithContext(ctx, &s3.PutObjectInput{
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

	someErr := xerrors.New("some error")

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
				Body: ioutil.NopCloser(bytes.NewReader(object)),
			},
			expected: object,
		},
		{
			name:         "NG: failed to GetObjectWithContext",
			getObjectErr: someErr,
			expectedErr:  xerrors.Errorf("failed to GetObjectWithContext: %w", someErr),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAPI.EXPECT().GetObjectWithContext(ctx, &s3.GetObjectInput{
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
