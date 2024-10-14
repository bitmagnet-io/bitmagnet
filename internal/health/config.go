package health

import (
	"context"
	"time"
)

type (
	// Check allows to configure health checks.
	Check struct {
		// The Name must be unique among all checks. Name is a required attribute.
		Name string // Required

		IsActive func() bool

		// Check is the check function that will be executed to check availability.
		// This function must return an error if the checked service is considered
		// not available. Check is a required attribute.
		Check func(ctx context.Context) error // Required

		// Timeout will override the global timeout value, if it is smaller than
		// the global timeout (see WithTimeout).
		Timeout time.Duration // Optional

		// MaxTimeInError will set a duration for how long a service must be
		// in an error state until it is considered down/unavailable.
		MaxTimeInError time.Duration // Optional

		// MaxContiguousFails will set a maximum number of contiguous
		// check fails until the service is considered down/unavailable.
		MaxContiguousFails uint // Optional

		// StatusListener allows to set a listener that will be called
		// whenever the AvailabilityStatus (e.g. from "up" to "down").
		StatusListener func(ctx context.Context, name string, state CheckState) // Optional

		// Interceptors holds a list of Interceptor instances that will be executed one after another in the
		// order as they appear in the list.
		Interceptors []Interceptor

		// DisablePanicRecovery disables automatic recovery from panics. If left in its default value (false),
		// panics will be automatically converted into errors instead.
		DisablePanicRecovery bool

		updateInterval time.Duration
		initialDelay   time.Duration
	}

	// CheckerOption is a configuration option for a Checker.
	CheckerOption func(config *checkerConfig)

	// HandlerOption is a configuration option for a Handler (see NewHandler).
	HandlerOption func(*HandlerConfig)
)

// NewChecker creates a new Checker. The provided options will be
// used to modify its configuration. If the Checker was not yet started
// (see Checker.IsStarted), it will be started automatically
// (see Checker.Start). You can disable this autostart by
// adding the WithDisabledAutostart configuration option.
func NewChecker(options ...CheckerOption) Checker {
	cfg := checkerConfig{
		cacheTTL:     1 * time.Second,
		timeout:      10 * time.Second,
		checks:       map[string]*Check{},
		interceptors: []Interceptor{},
	}

	for _, opt := range options {
		opt(&cfg)
	}

	return newChecker(cfg)
}

// WithDisabledDetails disables all data in the JSON response body. The AvailabilityStatus will be the only
// content. Example: { "status":"down" }. Enabled by default.
func WithDisabledDetails() CheckerOption {
	return func(cfg *checkerConfig) {
		cfg.detailsDisabled = true
	}
}

// WithTimeout defines a timeout duration for all checks. You can override
// this timeout by using the timeout value in the Check configuration.
// Default value is 10 seconds.
func WithTimeout(timeout time.Duration) CheckerOption {
	return func(cfg *checkerConfig) {
		cfg.timeout = timeout
	}
}

// WithStatusListener registers a listener function that will be called whenever the overall/aggregated system health
// status changes (e.g. from "up" to "down"). Attention: Because this listener is also executed for synchronous
// (i.e, request-based) health checks, it should not block processing.
func WithStatusListener(listener func(ctx context.Context, state CheckerState)) CheckerOption {
	return func(cfg *checkerConfig) {
		cfg.statusChangeListener = listener
	}
}

// WithMiddleware configures a middleware that will be used by the handler
// to pro- and post-process HTTP requests and health checks.
// Refer to the documentation of type Middleware for more information.
func WithMiddleware(middleware ...Middleware) HandlerOption {
	return func(cfg *HandlerConfig) {
		cfg.middleware = append(cfg.middleware, middleware...)
	}
}

// WithStatusCodeUp sets an HTTP status code that will be used for responses
// where the system is considered to be available ("up").
// Default is HTTP status code 200 (OK).
func WithStatusCodeUp(httpStatus int) HandlerOption {
	return func(cfg *HandlerConfig) {
		cfg.statusCodeUp = httpStatus
	}
}

// WithStatusCodeDown sets an HTTP status code that will be used for responses
// where the system is considered to be unavailable ("down").
// Default is HTTP status code 503 (Service Unavailable).
func WithStatusCodeDown(httpStatus int) HandlerOption {
	return func(cfg *HandlerConfig) {
		cfg.statusCodeDown = httpStatus
	}
}

// WithResultWriter is responsible for writing a health check result (see CheckerResult)
// into an HTTP response. By default, JSONResultWriter will be used.
func WithResultWriter(writer ResultWriter) HandlerOption {
	return func(cfg *HandlerConfig) {
		cfg.resultWriter = writer
	}
}

// WithDisabledAutostart disables automatic startup of a Checker instance.
func WithDisabledAutostart() CheckerOption {
	return func(cfg *checkerConfig) {
		cfg.autostartDisabled = true
	}
}

// WithDisabledCache disabled the check cache. This is not recommended in most cases.
// This will effectively lead to a health endpoint that initiates a new health check for each incoming HTTP request.
// This may have an impact on the systems that are being checked (especially if health checks are expensive).
// Caching also mitigates "denial of service" attacks. Caching is enabled by default.
func WithDisabledCache() CheckerOption {
	return WithCacheDuration(0)
}

// WithCacheDuration sets the duration for how long the aggregated health check result will be
// cached. By default, the cache TTL (i.e, the duration for how long responses will be cached) is set to 1 second.
// Caching will prevent that each incoming HTTP request triggers a new health check. A duration of 0 will
// effectively disable the cache and has the same effect as WithDisabledCache.
func WithCacheDuration(duration time.Duration) CheckerOption {
	return func(cfg *checkerConfig) {
		cfg.cacheTTL = duration
	}
}

// WithCheck adds a new health check that contributes to the overall service availability status.
// This check will be triggered each time Checker.Check is called (i.e., for each HTTP request).
// If health checks are expensive, or you expect a higher amount of requests on the health endpoint,
// consider using WithPeriodicCheck instead.
func WithCheck(check Check) CheckerOption {
	return func(cfg *checkerConfig) {
		cfg.checks[check.Name] = &check
	}
}

// WithPeriodicCheck adds a new health check that contributes to the overall service availability status.
// The health check will be performed on a fixed schedule and will not be executed for each HTTP request
// (as in contrast to WithCheck). This allows to process a much higher number of HTTP requests without
// actually calling the checked services too often or to execute long-running checks.
// This way Checker.Check (and the health endpoint) always returns the last result of the periodic check.
func WithPeriodicCheck(refreshPeriod time.Duration, initialDelay time.Duration, check Check) CheckerOption {
	return func(cfg *checkerConfig) {
		check.updateInterval = refreshPeriod
		check.initialDelay = initialDelay
		cfg.checks[check.Name] = &check
	}
}

// WithInterceptors adds a list of interceptors that will be applied to every check function. Interceptors
// may intercept the function call and do some pre- and post-processing, having the check state and check function
// result at hand. The interceptors will be executed in the order they are passed to this function.
func WithInterceptors(interceptors ...Interceptor) CheckerOption {
	return func(cfg *checkerConfig) {
		cfg.interceptors = interceptors
	}
}

// WithInfo sets values that will be available in every health check result. For example, you can use this option
// if you want to set information about your system that will be returned in every health check result, such as
// version number, Git SHA, build date, etc. These values will be available in CheckerResult.Info. If you use the
// default HTTP handler of this library (see NewHandler) or convert the CheckerResult to JSON on your own,
// these values will be available in the "info" field.
func WithInfo(values map[string]interface{}) CheckerOption {
	return func(cfg *checkerConfig) {
		cfg.info = values
	}
}
