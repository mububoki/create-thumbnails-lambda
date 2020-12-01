package s3

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/interface/gateway"
	"golang.org/x/xerrors"
)

var _ gateway.ObjectStorage = (*Handler)(nil)

type Handler struct {
	s3 s3.S3
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Save(ctx context.Context, object []byte, key, bucketName string) error {
	if _, err := h.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(object),
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}); err != nil {
		return xerrors.Errorf("failed to PutObjectWithContext: %w", err)
	}

	return nil
}

func (h *Handler) Find(ctx context.Context, key, bucketName string) ([]byte, error) {
	out, err := h.s3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, xerrors.Errorf("failed to GetObjectWithContext: %w", err)
	}

	body, err := ioutil.ReadAll(out.Body)
	if err != nil {
		return nil, xerrors.Errorf("failed to ReadAll: %w", err)
	}

	return body, nil
}
