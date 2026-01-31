package auth_test

import (
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	t.Parallel()

	seenStrings := make(map[string]struct{})

	length := 10

	for range 10 {
		length++

		for range 100 {
			str := auth.GenerateRandomString(length)

			assert.Len(t, str, length)

			_, seen := seenStrings[str]

			assert.False(t, seen)

			seenStrings[str] = struct{}{}
		}
	}
}
