//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../test/mock/mock_$GOPACKAGE/$GOFILE

package gateway

import (
	"context"
)

type AuthNZ interface {
	CreateRole(ctx context.Context, roleName string, serviceName string, actions []string) error
}
