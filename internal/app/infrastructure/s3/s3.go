//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../test/mock/mock_$GOPACKAGE/s3_api.go

package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/interface/gateway"
)

var _ gateway.ObjectStorage = (*Handler)(nil)

type S3API interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type Handler struct {
	s3 S3API
}

func NewHandler() (*Handler, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to LoadDefaultConfig: %w", err)
	}

	return &Handler{
		s3: s3.NewFromConfig(cfg),
	}, nil
}

func (h *Handler) Save(ctx context.Context, object []byte, key, bucketName string) error {
	if _, err := h.s3.PutObject(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(object),
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}); err != nil {
		return fmt.Errorf("failed to PutObject: %w", err)
	}

	return nil
}

func (h *Handler) Find(ctx context.Context, key, bucketName string) ([]byte, error) {
	out, err := h.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to GetObject: %w", err)
	}

	body, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to ReadAll: %w", err)
	}

	return body, nil
}
