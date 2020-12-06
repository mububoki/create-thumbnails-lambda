package interactor

import (
	"context"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/domain"
)

type ImageInteractor interface {
	CreateThumbnail(ctx context.Context, name string, format domain.ImageFormat) error
}
