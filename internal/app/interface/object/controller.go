package object

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/domain"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/interface/controller"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/usecase/interactor"
)

var _ controller.ObjectController = (*Controller)(nil)

type Controller struct {
	interactor          interactor.ImageInteractor
	bucketNameThumbnail string
}

func NewController(interactor interactor.ImageInteractor, bucketNameThumbnail string) *Controller {
	return &Controller{
		interactor:          interactor,
		bucketNameThumbnail: bucketNameThumbnail,
	}
}

func (c *Controller) CreateThumbnail(ctx context.Context, key string, bucketName string) error {
	if bucketName == c.bucketNameThumbnail {
		return errors.New("src bucket and dst bucket is the same")
	}

	name, format, err := extractNameAndFormat(key)
	if err != nil {
		return fmt.Errorf("failed to extractNameAndFormat: %w", err)
	}

	return c.interactor.CreateThumbnail(ctx, name, format)
}

func extractNameAndFormat(key string) (string, domain.ImageFormat, error) {
	lastIDX := strings.LastIndex(key, ".")
	if lastIDX < 0 || lastIDX > len(key)-1 {
		return "", domain.ImageFormat(0), errors.New("misspecified separator")
	}

	var format domain.ImageFormat
	if err := format.UnmarshalText([]byte(key[lastIDX+1:])); err != nil {
		return "", domain.ImageFormat(0), fmt.Errorf("failed to UnmarshalText: %w", err)
	}

	return key[:lastIDX], format, nil
}
