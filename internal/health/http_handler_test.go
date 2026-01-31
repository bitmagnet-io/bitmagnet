package health_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/health"
	health_mocks "github.com/bitmagnet-io/bitmagnet/internal/health/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func doTestHandler(t *testing.T,
	statusCodeUp, statusCodeDown int, expectedStatus health.CheckerResult, expectedStatusCode int,
) {
	t.Helper()

	// Arrange
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "https://localhost/foo", nil)

	ckr := health_mocks.NewChecker(t)
	ckr.On("IsStarted").Return(false)
	ckr.On("Start")
	ckr.On("Check", mock.Anything).Return(expectedStatus)

	handler := health.NewHandler(
		ckr,
		health.WithStatusCodeUp(statusCodeUp),
		health.WithStatusCodeDown(statusCodeDown),
	)

	// Act
	handler.ServeHTTP(response, request)

	// Assert
	ckr.AssertNumberOfCalls(t, "Check", 1)
	assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("content-type"))
	assert.Equal(t, expectedStatusCode, response.Code)

	bytes := response.Body.Bytes()
	expectedBytes, err := json.Marshal(expectedStatus)
	require.NoError(t, err)
	assert.Equal(t, string(expectedBytes), string(bytes))
}

func TestHandlerIfCheckFailThenRespondWithNotAvailable(t *testing.T) {
	t.Parallel()

	status := health.CheckerResult{
		Status: health.StatusUnknown,
		Details: map[string]health.CheckResult{
			"check1": {Status: health.StatusDown, Timestamp: time.Now(), Error: fmt.Errorf("hello")},
			"check2": {Status: health.StatusUp, Timestamp: time.Now(), Error: nil},
		},
	}

	doTestHandler(t, http.StatusNoContent, http.StatusTeapot, status, http.StatusTeapot)
}

func TestHandlerIfCheckSucceedsThenRespondWithAvailable(t *testing.T) {
	t.Parallel()

	status := health.CheckerResult{
		Status: health.StatusUp,
		Details: map[string]health.CheckResult{
			"check1": {Status: health.StatusUp, Timestamp: time.Now(), Error: nil},
		},
	}

	doTestHandler(t, http.StatusNoContent, http.StatusTeapot, status, http.StatusNoContent)
}

func TestHandlerIfAuthFailsThenReturnNoDetails(t *testing.T) {
	t.Parallel()

	status := health.CheckerResult{
		Status: health.StatusDown,
		Details: map[string]health.CheckResult{
			"check1": {
				Status:    health.StatusDown,
				Timestamp: time.Now(),
				Error:     fmt.Errorf("an error message"),
			},
		},
	}
	doTestHandler(t, http.StatusNoContent, http.StatusTeapot, status, http.StatusTeapot)
}

func TestWhenChecksEmptyThenHandlerResultContainNoChecksMap(t *testing.T) {
	t.Parallel()

	// Arrange
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Act
	health.NewHandler(health.NewChecker()).ServeHTTP(w, r)

	// Assert
	if w.Body.String() != "{\"status\":\"up\"}" {
		t.Errorf("response does not contain the expected result")
	}
}
