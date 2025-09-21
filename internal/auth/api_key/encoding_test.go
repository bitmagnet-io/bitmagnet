package api_key_test

import (
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/api_key"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestEncoding(t *testing.T) {
	t.Parallel()

	secret := api_key.NewSecret()

	require.NoError(t, bcrypt.CompareHashAndPassword(secret.Hash, secret.Secret))

	apiKeyID := 12345

	keyData := &api_key.KeyData{
		ID:     apiKeyID,
		Secret: secret.Secret,
	}

	encoded := keyData.Encode()
	assert.Len(t, encoded, 22)

	keyData2 := &api_key.KeyData{}
	require.NoError(t, keyData2.Decode(encoded))

	assert.Equal(t, keyData, keyData2)
}
