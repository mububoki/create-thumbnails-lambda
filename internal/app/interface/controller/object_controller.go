package controller

import "context"

type ObjectController interface {
	CreateThumbnail(ctx context.Context, key string, bucketName string) error
}
