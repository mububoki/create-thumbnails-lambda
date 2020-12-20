//go:generate mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=../../test/mock/mock$GOPACKAGE/$GOFILE

package port

import (
	"context"
)

type AuthNZRepository interface {
	CreateRole(ctx context.Context, roleName string, actions []string) error
}
