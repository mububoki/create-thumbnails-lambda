package port

import (
	"context"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/domain"
)

type ImageRepository interface {
	Save(context.Context, *domain.Image) error
	Search(ctx context.Context, name, format string) (*domain.Image, error)
}
