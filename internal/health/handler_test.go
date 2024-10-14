package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type checkerMock struct {
	mock.Mock
}

func (ck *checkerMock) Start() {
	ck.Called()
}

func (ck *checkerMock) Stop() {
	ck.Called()
}

func (ck *checkerMock) Check(ctx context.Context) CheckerResult {
	return ck.Called(ctx).Get(0).(CheckerResult)
}

func (ck *checkerMock) GetRunningPeriodicCheckCount() int {
	return ck.Called().Get(0).(int)
}

func (ck *checkerMock) IsStarted() bool {
	return ck.Called().Get(0).(bool)
}

func (ck *checkerMock) StartedAt() time.Time {
	return ck.Called().Get(0).(time.Time)
}

type resultWriterMock struct {
	mock.Mock
}

func (ck *resultWriterMock) Write(result *CheckerResult, statusCode int, w http.ResponseWriter, r *http.Request) error {
	return ck.Called(result, statusCode, w, r).Get(0).(error)
}

//var testTimestamp = time.Now()

func doTestHandler(t *testing.T, statusCodeUp, statusCodeDown int, expectedStatus CheckerResult, expectedStatusCode int) {
	// Arrange
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://localhost/foo", nil)

	ckr := checkerMock{}
	ckr.On("IsStarted").Return(false)
	ckr.On("Start")
	ckr.On("Check", mock.Anything).Return(expectedStatus)

	handler := NewHandler(&ckr, WithStatusCodeUp(statusCodeUp), WithStatusCodeDown(statusCodeDown))

	// Act
	handler.ServeHTTP(response, request)

	// Assert
	ckr.Mock.AssertNumberOfCalls(t, "Check", 1)
	assert.Equal(t, response.Header().Get("content-type"), "application/json; charset=utf-8")
	assert.Equal(t, response.Result().StatusCode, expectedStatusCode)

	bytes := response.Body.Bytes()
	expectedBytes, err := json.Marshal(expectedStatus)
	assert.NoError(t, err)
	assert.Equal(t, string(expectedBytes), string(bytes))
}

func TestHandlerIfCheckFailThenRespondWithNotAvailable(t *testing.T) {
	status := CheckerResult{
		Status: StatusUnknown,
		Details: map[string]CheckResult{
			"check1": {Status: StatusDown, Timestamp: time.Now(), Error: fmt.Errorf("hello")},
			"check2": {Status: StatusUp, Timestamp: time.Now(), Error: nil},
		},
	}

	doTestHandler(t, http.StatusNoContent, http.StatusTeapot, status, http.StatusTeapot)
}

func TestHandlerIfCheckSucceedsThenRespondWithAvailable(t *testing.T) {
	status := CheckerResult{
		Status: StatusUp,
		Details: map[string]CheckResult{
			"check1": {Status: StatusUp, Timestamp: time.Now(), Error: nil},
		},
	}

	doTestHandler(t, http.StatusNoContent, http.StatusTeapot, status, http.StatusNoContent)
}

func TestHandlerIfAuthFailsThenReturnNoDetails(t *testing.T) {
	status := CheckerResult{
		Status: StatusDown,
		Details: map[string]CheckResult{
			"check1": {Status: StatusDown, Timestamp: time.Now(), Error: fmt.Errorf("an error message")},
		},
	}
	doTestHandler(t, http.StatusNoContent, http.StatusTeapot, status, http.StatusTeapot)
}

func TestWhenChecksEmptyThenHandlerResultContainNoChecksMap(t *testing.T) {
	// Arrange
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Act
	NewHandler(NewChecker()).ServeHTTP(w, r)

	// Assert
	if w.Body.String() != "{\"status\":\"up\"}" {
		t.Errorf("response does not contain the expected result")
	}

}
