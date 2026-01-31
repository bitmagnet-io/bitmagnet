package health

import (
	"encoding/json"
	"errors"
	"time"
)

type (
	// CheckerState represents the current state of the Checker.
	CheckerState struct {
		// Status is the aggregated system health status.
		Status AvailabilityStatus
		// CheckState contains the state of all checks.
		CheckState map[string]CheckState
	}

	// CheckState represents the current state of a component check.
	CheckState struct {
		// LastCheckedAt holds the time of when the check was last executed.
		LastCheckedAt time.Time
		// LastCheckedAt holds the last time of when the check did not return an error.
		LastSuccessAt time.Time
		// LastFailureAt holds the last time of when the check did return an error.
		LastFailureAt time.Time
		// FirstCheckStartedAt holds the time of when the first check was started.
		FirstCheckStartedAt time.Time
		// ContiguousFails holds the number of how often the check failed in a row.
		ContiguousFails uint
		// Result holds the error of the last check (nil if successful).
		Result error
		// The current availability status of the check.
		Status AvailabilityStatus
	}

	// CheckerResult holds the aggregated system availability status and
	// detailed information about the individual checks.
	CheckerResult struct {
		// Info contains additional information about this health result.
		Info map[string]interface{} `json:"info,omitempty"`
		// Status is the aggregated system availability status.
		Status AvailabilityStatus `json:"status"`
		// Details contains health information for all checked components.
		Details map[string]CheckResult `json:"details,omitempty"`
	}

	// CheckResult holds a components health information.
	// Attention: This type is converted from/to JSON using a custom
	// marshalling/unmarshalling function (see type jsonCheckResult).
	// This is required because some fields are not converted automatically
	// by the standard json.Marshal/json.Unmarshal functions
	// (such as the error interface). The JSON tags you see here, are
	// just there for the readers' convenience.
	CheckResult struct {
		// Status is the availability status of a component.
		Status AvailabilityStatus `json:"status"`
		// Timestamp holds the time when the check was executed.
		Timestamp time.Time `json:"timestamp,omitempty"`
		// Error contains the check error message, if the check failed.
		Error error `json:"error,omitempty"`
	}
)

const (
	// StatusUnknown holds the information that the availability
	// status is not known, because not all checks were executed yet.
	StatusUnknown AvailabilityStatus = "unknown"
	// StatusUp holds the information that the system or a component
	// is up and running.
	StatusUp AvailabilityStatus = "up"
	// StatusDown holds the information that the system or a component
	// down and not available.
	StatusDown AvailabilityStatus = "down"
	// StatusInactive holds the information that a component
	// is not currently active.
	StatusInactive AvailabilityStatus = "inactive"
)

// MarshalJSON provides a custom marshaller for the CheckResult type.
func (cr CheckResult) MarshalJSON() ([]byte, error) {
	errorMsg := ""
	if cr.Error != nil {
		errorMsg = cr.Error.Error()
	}

	return json.Marshal(&jsonCheckResult{
		Status:    string(cr.Status),
		Timestamp: cr.Timestamp,
		Error:     errorMsg,
	})
}

func (cr *CheckResult) UnmarshalJSON(data []byte) error {
	var result jsonCheckResult
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	cr.Status = AvailabilityStatus(result.Status)
	cr.Timestamp = result.Timestamp

	if result.Error != "" {
		cr.Error = errors.New(result.Error)
	}

	return nil
}

func (s AvailabilityStatus) criticality() int {
	switch s {
	case StatusDown:
		return 2
	case StatusUnknown:
		return 1
	default:
		return 0
	}
}
