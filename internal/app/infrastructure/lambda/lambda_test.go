package lambda

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/test/mock/mock_controller"
)

func TestHandler_handleLambdaS3Events(t *testing.T) {
	ctx := context.Background()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockObjectController := mock_controller.NewMockObjectController(mockController)
	handler := NewHandler(mockObjectController)

	someErr := xerrors.New("some error")

	event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{
						Name: "name0",
					},
					Object: events.S3Object{
						Key: "key0",
					},
				},
			},
			{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{
						Name: "name1",
					},
					Object: events.S3Object{
						Key: "key1",
					},
				},
			},
		},
	}

	testCases := []struct {
		name                string
		event               events.S3Event
		createThumbnailErrs []error
		expectedErr         error
	}{
		{
			name:                "OK",
			event:               event,
			createThumbnailErrs: []error{nil, nil},
		},
		{
			name:                "NG",
			event:               event,
			createThumbnailErrs: []error{nil, someErr},
			expectedErr:         xerrors.Errorf("failed to CreateThumbnail: %w", someErr),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Len(t, tc.createThumbnailErrs, len(tc.event.Records))
			for i, record := range tc.event.Records {
				mockObjectController.EXPECT().CreateThumbnail(ctx, record.S3.Object.Key, record.S3.Bucket.Name).Return(tc.createThumbnailErrs[i])
			}

			err := handler.handleLambdaS3Event(ctx, tc.event)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
