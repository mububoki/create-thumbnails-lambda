package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/xerrors"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/interface/controller"
)

type Handler struct {
	controller controller.ObjectController
}

func NewHandler(controller controller.ObjectController) *Handler {
	return &Handler{
		controller: controller,
	}
}

func (h *Handler) CreateThumbnail() {
	lambda.Start(h.handleLambdaS3Event)
}

func (h *Handler) handleLambdaS3Event(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {
		if err := h.controller.CreateThumbnail(ctx, record.S3.Object.Key, record.S3.Bucket.Name); err != nil {
			return xerrors.Errorf("failed to CreateThumbnail: %w", err)
		}
	}

	return nil
}
