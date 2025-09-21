package api_key_test

import (
	"testing"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/api_key"
	api_key_mocks "github.com/bitmagnet-io/bitmagnet/internal/auth/api_key/mocks"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testHarness struct {
	repository *api_key_mocks.Repository
	service    api_key.Service
}

func newTestHarness(t *testing.T) testHarness {
	repository := api_key_mocks.NewRepository(t)

	return testHarness{
		repository: repository,
		service:    api_key.NewService(repository),
	}
}

func TestService(t *testing.T) {
	t.Parallel()

	h := newTestHarness(t)

	permissions := []rbac.ObjectAction{
		{
			Namespace: "test",
			Object:    "foo",
			Action:    "bar",
		},
	}

	h.repository.EXPECT().
		Create(
			t.Context(),
			1,
			"test",
			mock.IsType([]byte{}),
			permissions,
			mock.IsType(time.Now()),
		).
		Return(2, nil).
		Once()

	result, err := h.service.Create(t.Context(), api_key.CreateRequest{
		UserID:      1,
		Name:        "test",
		Permissions: permissions,
		Expiry:      time.Hour * 24,
	})

	require.NoError(t, err)

	assert.Equal(t, 2, result.ID)
	assert.Len(t, result.APIKey, 22)
	assert.WithinDuration(t, time.Now().Add(time.Hour*24), result.ExpiresAt, time.Second)
}
