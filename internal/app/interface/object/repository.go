package object

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/domain"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/interface/gateway"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/usecase/port"
)

var _ port.ImageRepository = (*Repository)(nil)

type Repository struct {
	storage             gateway.ObjectStorage
	bucketNameOriginal  string
	bucketNameThumbnail string
}

func NewRepository(storage gateway.ObjectStorage, bucketNameOriginal string, bucketNameThumbnail string) *Repository {
	return &Repository{
		storage:             storage,
		bucketNameOriginal:  bucketNameOriginal,
		bucketNameThumbnail: bucketNameThumbnail,
	}
}

func (repo *Repository) Save(ctx context.Context, img *domain.Image) error {
	object, err := img.Encode()
	if err != nil {
		return xerrors.Errorf("failed to Encode: %w", err)
	}

	return repo.storage.Save(ctx, object, keyImage(img.Name, img.Format), repo.bucketNameImage(img.IsThumbnail))
}

func (repo *Repository) Find(ctx context.Context, name string, format domain.ImageFormat, isThumbnail bool) (*domain.Image, error) {
	bytesIMG, err := repo.storage.Find(ctx, keyImage(name, format), repo.bucketNameImage(isThumbnail))
	if err != nil {
		return nil, xerrors.Errorf("failed to Find: %w", err)
	}

	return domain.DecodeImage(bytesIMG, name, isThumbnail)
}

func keyImage(name string, format domain.ImageFormat) string {
	return name + "." + format.String()
}

func (repo *Repository) bucketNameImage(isThumbnail bool) string {
	if isThumbnail {
		return repo.bucketNameThumbnail
	}

	return repo.bucketNameOriginal
}
