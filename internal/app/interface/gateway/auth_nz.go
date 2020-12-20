//go:generate mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=../../test/mock/mock$GOPACKAGE/$GOFILE

package gateway

import (
	"context"
)

type AuthNZ interface {
	CreateRole(ctx context.Context, roleName string, serviceName string, actions []string) error
}
