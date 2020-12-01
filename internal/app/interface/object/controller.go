package object

import (
	"context"
	"strings"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/domain"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/usecase/image"
	"golang.org/x/xerrors"
)

type Controller struct {
	interactor *image.Interactor
}

func NewController(interactor *image.Interactor) *Controller {
	return &Controller{
		interactor: interactor,
	}
}

func (c *Controller) CreateThumbnail(ctx context.Context, key string, bucketName string) error {
	if bucketName == bucketNameThumbnail {
		return xerrors.New("src bucket and dst bucket is the same")
	}

	name, format, err := extractNameAndFormat(key)
	if err != nil {
		return xerrors.Errorf("faield to extractNameAndFormat: %w", err)
	}

	return c.interactor.CreateThumbnail(ctx, name, format)
}

func extractNameAndFormat(key string) (string, domain.ImageFormat, error) {
	lastIDX := strings.LastIndex(key, ".")
	if lastIDX < 0 || lastIDX > len(key)-1 {
		return "", domain.ImageFormat(0), xerrors.New("misspecified separator")
	}

	var format domain.ImageFormat
	if err := format.UnmarshalText([]byte(key[lastIDX+1:])); err != nil {
		return "", domain.ImageFormat(0), xerrors.Errorf("failed to UnmarshalText: %w", err)
	}

	return key[:lastIDX], format, nil
}
