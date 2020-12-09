//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../test/mock/mock_$GOPACKAGE/$GOFILE

package gateway

import (
	"context"
)

type ObjectStorage interface {
	Save(ctx context.Context, object []byte, key, bucketName string) error
	Find(ctx context.Context, key, bucketName string) ([]byte, error)
}
