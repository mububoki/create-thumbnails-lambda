package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/interface/object"
	"golang.org/x/xerrors"
)

type Handler struct {
	controller *object.Controller
}

func NewHandler(controller *object.Controller) *Handler {
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
