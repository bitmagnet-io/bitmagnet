package health

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithPeriodicCheckConfig(t *testing.T) {
	t.Parallel()

	// Arrange
	expectedName := "test"
	cfg := checkerConfig{checks: map[string]*Check{}}
	interval := 5 * time.Second
	initialDelay := 1 * time.Minute
	check := Check{Name: expectedName, updateInterval: interval, initialDelay: initialDelay}

	// Act
	WithPeriodicCheck(interval, initialDelay, check)(&cfg)

	// Assert
	assert.Len(t, cfg.checks, 1)
	assert.True(t, reflect.DeepEqual(check, *cfg.checks[expectedName]))
}

func TestWithCheckConfig(t *testing.T) {
	t.Parallel()

	// Arrange
	expectedName := "test"
	cfg := checkerConfig{checks: map[string]*Check{}}
	check := Check{Name: "test"}

	// Act
	WithCheck(check)(&cfg)

	// Assert
	require.Len(t, cfg.checks, 1)
	assert.True(t, reflect.DeepEqual(&check, cfg.checks[expectedName]))
}

func TestWithCacheDurationConfig(t *testing.T) {
	t.Parallel()

	// Arrange
	cfg := checkerConfig{}
	duration := 5 * time.Hour

	// Act
	WithCacheDuration(duration)(&cfg)

	// Assert
	assert.Equal(t, duration, cfg.cacheTTL)
}

func TestWithDisabledCacheConfig(t *testing.T) {
	t.Parallel()

	// Arrange
	cfg := checkerConfig{}

	// Act
	WithDisabledCache()(&cfg)

	// Assert
	assert.Equal(t, 0*time.Second, cfg.cacheTTL)
}

func TestWithTimeoutStartConfig(t *testing.T) {
	t.Parallel()

	// Arrange
	cfg := checkerConfig{}

	// Act
	WithTimeout(5 * time.Hour)(&cfg)

	// Assert
	assert.Equal(t, 5*time.Hour, cfg.timeout)
}

func TestWithDisabledDetailsConfig(t *testing.T) {
	t.Parallel()

	// Arrange
	cfg := checkerConfig{}

	// Act
	WithDisabledDetails()(&cfg)

	// Assert
	assert.True(t, cfg.detailsDisabled)
}

func TestWithMiddlewareConfig(t *testing.T) {
	t.Parallel()

	// Arrange
	cfg := HandlerConfig{}
	mw := func(MiddlewareFunc) MiddlewareFunc {
		return func(*http.Request) CheckerResult {
			return CheckerResult{nil, StatusUp, nil}
		}
	}

	// Act
	WithMiddleware(mw)(&cfg)

	// Assert
	assert.Len(t, cfg.middleware, 1)
}

func TestWithInterceptorConfig(t *testing.T) {
	t.Parallel()

	// Arrange
	cfg := checkerConfig{}
	interceptor := func(InterceptorFunc) InterceptorFunc {
		return func(context.Context, string, CheckState) CheckState {
			return CheckState{}
		}
	}

	// Act
	WithInterceptors(interceptor)(&cfg)

	// Assert
	assert.Len(t, cfg.interceptors, 1)
}

func TestWithResultWriterConfig(t *testing.T) {
	t.Parallel()

	// Arrange
	cfg := HandlerConfig{}
	w := resultWriterMock{}

	// Act
	WithResultWriter(&w)(&cfg)

	// Assert
	assert.Equal(t, &w, cfg.resultWriter)
}

func TestWithStatusChangeListenerConfig(t *testing.T) {
	t.Parallel()

	// Arrange
	cfg := checkerConfig{}

	// Act
	// Use of non standard AvailabilityStatus codes.
	WithStatusListener(func(context.Context, CheckerState) {})(&cfg)

	// Assert
	// Not possible in Go to compare functions.
	assert.NotNil(t, cfg.statusChangeListener)
}

func TestNewWithDefaults(t *testing.T) {
	t.Parallel()

	// Arrange
	configApplied := false
	opt := func(*checkerConfig) { configApplied = true }

	// Act
	checker := NewChecker(opt)

	// Assert
	ckr := checker.(*defaultChecker)
	assert.Equal(t, 1*time.Second, ckr.cfg.cacheTTL)
	assert.Equal(t, 10*time.Second, ckr.cfg.timeout)
	assert.True(t, configApplied)
}

func TestNewCheckerWithDefaults(t *testing.T) {
	t.Parallel()

	// Arrange
	configApplied := false
	opt := func(*checkerConfig) { configApplied = true }

	// Act
	checker := NewChecker(opt)

	// Assert
	ckr := checker.(*defaultChecker)
	assert.Equal(t, 1*time.Second, ckr.cfg.cacheTTL)
	assert.Equal(t, 10*time.Second, ckr.cfg.timeout)
	assert.True(t, configApplied)
}

func TestCheckerAutostartConfig(t *testing.T) {
	t.Parallel()

	// Arrange + Act
	c := NewChecker()

	// Assert
	assert.True(t, c.IsStarted())
}

func TestCheckerAutostartDisabledConfig(t *testing.T) {
	t.Parallel()

	// Arrange
	c := NewChecker(WithDisabledAutostart())

	// Assert
	assert.False(t, c.IsStarted())
}
