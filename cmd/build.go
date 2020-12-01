package main

import (
	"github.com/mububoki/create-thumbnails-lambda/internal/app/infrastructure/env"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/infrastructure/lambda"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/infrastructure/s3"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/interface/object"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/usecase/image"
)

func build() *lambda.Handler {
	s3Handler := s3.NewHandler()
	objectRepository := object.NewRepository(s3Handler, env.Object.BucketNameOriginal, env.Object.BucketNameThumbnail)
	imageInteractor := image.NewInteractor(objectRepository, env.Image.Rate)
	objectController := object.NewController(imageInteractor, env.Object.BucketNameThumbnail)

	return lambda.NewHandler(objectController)
}
