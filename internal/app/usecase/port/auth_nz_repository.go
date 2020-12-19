//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../test/mock/mock_$GOPACKAGE/$GOFILE

package port

import (
	"context"
)

type AuthNZRepository interface {
	CreateRole(ctx context.Context, roleName string, actions []string) error
}
