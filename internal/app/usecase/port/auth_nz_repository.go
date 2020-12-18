package port

import (
	"context"
)

type AuthNZRepository interface {
	CreateRole(ctx context.Context, roleName string, actions []string) error
}
