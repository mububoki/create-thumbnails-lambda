package image

import (
	"context"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/domain"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/usecase/port"
	"golang.org/x/xerrors"
)

const (
	rate = 1 / 2
)

type Interactor struct {
	repo port.ImageRepository
}

func NewInteractor(repo port.ImageRepository) *Interactor {
	return &Interactor{
		repo: repo,
	}
}

func (i *Interactor) CreateThumbnail(ctx context.Context, name string, format domain.ImageFormat) error {
	src, err := i.repo.Find(ctx, name, format, false)
	if err != nil {
		return xerrors.Errorf("failed to Search: %w", err)
	}

	dst, err := src.CreateThumbnail(rate)
	if err != nil {
		return xerrors.Errorf("failed to CreateThumbnail: %w", err)
	}

	if err := i.repo.Save(ctx, dst); err != nil {
		return xerrors.Errorf("failed to Save: %w", err)
	}

	return nil
}
