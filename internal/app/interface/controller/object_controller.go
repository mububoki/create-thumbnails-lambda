//go:generate mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=../../test/mock/mock$GOPACKAGE/$GOFILE

package controller

import (
	"context"
)

type ObjectController interface {
	CreateThumbnail(ctx context.Context, key string, bucketName string) error
}
