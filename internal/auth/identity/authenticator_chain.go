package identity

import (
	"context"
	"fmt"
)

type authenticatorChain []Authenticator

func (a authenticatorChain) Authenticate(ctx context.Context, token string) (Identity, bool, error) {
	for _, p := range a {
		id, matched, err := p.Authenticate(ctx, token)
		if matched {
			if err != nil {
				return nil, true, fmt.Errorf("%w: %w: %w", Err, ErrAuthentication, err)
			}

			return id, true, nil
		}
	}

	return nil, false, fmt.Errorf("%w: %w", Err, ErrUnmatched)
}
