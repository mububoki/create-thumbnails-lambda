//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../test/mock/mock_$GOPACKAGE/$GOFILE

package port

import (
	"context"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/domain"
)

type ImageRepository interface {
	Save(context.Context, *domain.Image) error
	Find(ctx context.Context, name string, format domain.ImageFormat, isThumbnail bool) (*domain.Image, error)
}
