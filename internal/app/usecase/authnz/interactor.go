package authnz

import (
	"context"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/usecase/port"
)

type Interactor struct {
	repo     port.AuthNZRepository
	roleName string
	actions  []string
}

func NewInteractor(repo port.AuthNZRepository, roleName string, actions []string) *Interactor {
	return &Interactor{
		repo:     repo,
		roleName: roleName,
		actions:  actions,
	}
}

func (i *Interactor) CreateRoleForCreatingThumbnails(ctx context.Context) error {
	return i.repo.CreateRole(ctx, i.roleName, i.actions)
}
