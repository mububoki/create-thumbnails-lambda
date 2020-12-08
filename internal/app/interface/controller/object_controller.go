//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../test/mock/mock_$GOPACKAGE/$GOFILE

package controller

import "context"

type ObjectController interface {
	CreateThumbnail(ctx context.Context, key string, bucketName string) error
}
