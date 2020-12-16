package gateway

import (
	"context"
)

type AuthNZ interface {
	CreateRole(ctx context.Context, roleName string, serviceName string, actions []string) error
}
