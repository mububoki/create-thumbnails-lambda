package object

import (
	"context"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/domain"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/interface/gateway"
	"golang.org/x/xerrors"
)

const (
	bucketNameOriginal  = "original"
	bucketNameThumbnail = "thumbnail"
)

type ObjectRepository struct {
	storage gateway.ObjectStorage
}

func (repo *ObjectRepository) Save(ctx context.Context, img *domain.Image) error {
	object, err := img.Encode()
	if err != nil {
		return xerrors.Errorf("failed to Encode: %w", err)
	}

	return repo.storage.Save(ctx, object, keyImage(img.Name, img.Format), bucketNameImage(img.IsThumbnail))
}

func (repo *ObjectRepository) Find(ctx context.Context, name string, format domain.ImageFormat, isThumbnail bool) (*domain.Image, error) {
	bytesIMG, err := repo.storage.Find(ctx, keyImage(name, format), bucketNameImage(isThumbnail))
	if err != nil {
		return nil, xerrors.Errorf("failed to Find: %w", err)
	}

	var img *domain.Image
	if err := img.Decode(bytesIMG, name, isThumbnail); err != nil {
		return nil, xerrors.Errorf("failed to Decode: %w", err)
	}

	return img, nil
}

func keyImage(name string, format domain.ImageFormat) string {
	return name + "." + format.String()
}

func bucketNameImage(isThumbnail bool) string {
	if isThumbnail {
		return bucketNameThumbnail
	}

	return bucketNameOriginal
}
