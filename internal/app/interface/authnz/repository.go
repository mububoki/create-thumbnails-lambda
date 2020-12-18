package authnz

import (
	"context"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/interface/gateway"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/usecase/port"
)

var _ port.AuthNZRepository = (*Repository)(nil)

type Repository struct {
	authNZ    gateway.AuthNZ
	informant gateway.ServiceInformant
}

func NewRepository(authNZ gateway.AuthNZ, informant gateway.ServiceInformant) *Repository {
	return &Repository{
		authNZ:    authNZ,
		informant: informant,
	}
}

func (r *Repository) CreateRole(ctx context.Context, roleName string, actions []string) error {
	return r.authNZ.CreateRole(ctx, roleName, r.informant.ServiceName(), actions)
}
