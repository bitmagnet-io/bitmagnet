package identity

import (
	"context"
)

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (Identity, bool, error)
}
