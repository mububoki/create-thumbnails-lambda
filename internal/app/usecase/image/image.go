package image

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/domain"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/usecase/port"
)

type Interactor struct {
	repo port.ImageRepository
	rate float64
}

func NewInteractor(repo port.ImageRepository, rate float64) *Interactor {
	return &Interactor{
		repo: repo,
		rate: rate,
	}
}

func (i *Interactor) CreateThumbnail(ctx context.Context, name string, format domain.ImageFormat) error {
	src, err := i.repo.Find(ctx, name, format, false)
	if err != nil {
		return xerrors.Errorf("failed to Search: %w", err)
	}

	dst, err := src.CreateThumbnail(i.rate)
	if err != nil {
		return xerrors.Errorf("failed to CreateThumbnail: %w", err)
	}

	if err := i.repo.Save(ctx, dst); err != nil {
		return xerrors.Errorf("failed to Save: %w", err)
	}

	return nil
}
