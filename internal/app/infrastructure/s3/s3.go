//go:generate mockgen -package=mock_$GOPACKAGE -destination=../../mock/mock_$GOPACKAGE/s3_api.go github.com/aws/aws-sdk-go/service/s3/s3iface S3API

package s3

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"golang.org/x/xerrors"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/interface/gateway"
)

var _ gateway.ObjectStorage = (*Handler)(nil)

type Handler struct {
	s3 s3iface.S3API
}

func NewHandler() (*Handler, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, xerrors.Errorf("failed to NewSession: %w", err)
	}

	return &Handler{
		s3: s3.New(sess),
	}, nil
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
