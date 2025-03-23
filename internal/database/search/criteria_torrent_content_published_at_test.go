package search

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// This file contains tests for the torrent published date filter criteria

func TestTorrentContentPublishedAtCriteria(t *testing.T) {
	// Test basic creation of criteria
	criteria := TorrentContentPublishedAtCriteria("7d")
	assert.NotNil(t, criteria, "Criteria should not be nil")
	
	// Also test empty string handling
	emptyTimeFrame := TorrentContentPublishedAtCriteria("")
	assert.NotNil(t, emptyTimeFrame)
}

// TestParseTimeFrameWithFixedTime tests parseTimeFrame with a fixed time to verify time-based calculations
func TestParseTimeFrameWithFixedTime(t *testing.T) {
	// Define a fixed time for consistent testing: January 15, 2023 10:30:00 UTC
	fixedTime := time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)
	
	// Store the original timeNow function
	originalTimeNow := timeNow
	
	// Replace with our fixed time function
	timeNow = func() time.Time {
		return fixedTime
	}
	
	// Restore original function when we're done
	defer func() { timeNow = originalTimeNow }()
	
	// Test cases for different time frames
	testCases := []struct {
		name            string
		timeFrame       string
		expectedErr     bool
		expectedGteTime time.Time
		expectedLteTime time.Time
	}{
		{
			name:            "7 days relative time",
			timeFrame:       "7d",
			expectedGteTime: fixedTime.AddDate(0, 0, -7),
			expectedLteTime: fixedTime,
		},
		{
			name:            "24 hours relative time",
			timeFrame:       "24h",
			expectedGteTime: fixedTime.Add(-24 * time.Hour),
			expectedLteTime: fixedTime,
		},
		{
			name:            "today special time",
			timeFrame:       "today",
			expectedGteTime: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			expectedLteTime: fixedTime,
		},
		{
			name:            "yesterday special time",
			timeFrame:       "yesterday",
			expectedGteTime: time.Date(2023, 1, 14, 0, 0, 0, 0, time.UTC),
			expectedLteTime: time.Date(2023, 1, 14, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:            "this month special time",
			timeFrame:       "this month",
			expectedGteTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedLteTime: fixedTime,
		},
		{
			name:            "this year special time",
			timeFrame:       "this year",
			expectedGteTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedLteTime: fixedTime,
		},
		{
			name:            "date range",
			timeFrame:       "2023-01-01 to 2023-01-31",
			expectedGteTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedLteTime: time.Date(2023, 1, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:            "single date",
			timeFrame:       "2023-01-01",
			expectedGteTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedLteTime: time.Date(2023, 1, 1, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:        "empty time frame",
			timeFrame:   "",
			expectedErr: false, // Should not produce an error, just zero times
		},
		{
			name:        "invalid time frame",
			timeFrame:   "invalid_format",
			expectedErr: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startTime, endTime, err := parseTimeFrame(tc.timeFrame)
			
			if tc.expectedErr {
				assert.Error(t, err, "Expected error for timeFrame: %s", tc.timeFrame)
				return
			}
			
			assert.NoError(t, err, "Unexpected error for timeFrame: %s", tc.timeFrame)
			
			if tc.timeFrame == "" {
				// Empty time frame returns zero times
				assert.True(t, startTime.IsZero(), "Start time should be zero for empty time frame")
				assert.True(t, endTime.IsZero(), "End time should be zero for empty time frame")
				return
			}
			
			// Round to seconds for consistent comparison
			startTime = startTime.Truncate(time.Second)
			endTime = endTime.Truncate(time.Second)
			expectedGteTime := tc.expectedGteTime.Truncate(time.Second)
			expectedLteTime := tc.expectedLteTime.Truncate(time.Second)
			
			assert.Equal(t, expectedGteTime, startTime,
				"Start time does not match expected for timeFrame: %s", tc.timeFrame)
			assert.Equal(t, expectedLteTime, endTime,
				"End time does not match expected for timeFrame: %s", tc.timeFrame)
		})
	}
}

func TestParseTimeFrame(t *testing.T) {
	// Define a fixed time for consistent testing
	fixedTime := time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)
	
	// Store the original timeNow function
	originalTimeNow := timeNow
	
	// Replace with our fixed time function
	timeNow = func() time.Time {
		return fixedTime
	}
	
	// Restore original function when we're done
	defer func() { timeNow = originalTimeNow }()
	
	// Define test cases
	tests := []struct {
		name           string
		timeFrame      string
		expectNonEmpty bool
		expectError    bool
	}{
		{
			name:           "empty string",
			timeFrame:      "",
			expectNonEmpty: false,
		},
		{
			name:           "invalid format",
			timeFrame:      "invalid_format",
			expectNonEmpty: false,
			expectError:    true,
		},
		{
			name:           "relative time - 7 days",
			timeFrame:      "7d",
			expectNonEmpty: true,
		},
		{
			name:           "relative time - 24 hours",
			timeFrame:      "24h",
			expectNonEmpty: true,
		},
		{
			name:           "relative time - 2 weeks",
			timeFrame:      "2w",
			expectNonEmpty: true,
		},
		{
			name:           "relative time - 3 months",
			timeFrame:      "3M",
			expectNonEmpty: true,
		},
		{
			name:           "relative time - 1 year",
			timeFrame:      "1y",
			expectNonEmpty: true,
		},
		{
			name:           "date range",
			timeFrame:      "2023-01-01 to 2023-01-31",
			expectNonEmpty: true,
		},
		{
			name:           "date range with slashes",
			timeFrame:      "2023/01/01 to 2023/01/31",
			expectNonEmpty: true,
		},
		{
			name:           "single date",
			timeFrame:      "2023-01-01",
			expectNonEmpty: true,
		},
		{
			name:           "special time - today",
			timeFrame:      "today",
			expectNonEmpty: true,
		},
		{
			name:           "special time - yesterday",
			timeFrame:      "yesterday",
			expectNonEmpty: true,
		},
		{
			name:           "special time - this week",
			timeFrame:      "this week",
			expectNonEmpty: true,
		},
		{
			name:           "special time - last week",
			timeFrame:      "last week",
			expectNonEmpty: true,
		},
		{
			name:           "special time - this month",
			timeFrame:      "this month",
			expectNonEmpty: true,
		},
		{
			name:           "special time - last month",
			timeFrame:      "last month",
			expectNonEmpty: true,
		},
		{
			name:           "special time - this year",
			timeFrame:      "this year",
			expectNonEmpty: true,
		},
		{
			name:           "special time - last year",
			timeFrame:      "last year",
			expectNonEmpty: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime, endTime, err := parseTimeFrame(tt.timeFrame)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			assert.NoError(t, err)
			
			if !tt.expectNonEmpty {
				assert.True(t, startTime.IsZero())
				assert.True(t, endTime.IsZero())
				return
			}
			
			// Verify non-zero times
			assert.False(t, startTime.IsZero())
			assert.False(t, endTime.IsZero())
			assert.True(t, endTime.After(startTime) || endTime.Equal(startTime))
			
			// For special time frames, do additional validation
			switch tt.timeFrame {
			case "today":
				assert.Equal(t, fixedTime.Year(), startTime.Year())
				assert.Equal(t, fixedTime.Month(), startTime.Month())
				assert.Equal(t, fixedTime.Day(), startTime.Day())
				assert.Equal(t, 0, startTime.Hour())
				assert.Equal(t, 0, startTime.Minute())
			case "yesterday":
				yesterday := fixedTime.AddDate(0, 0, -1)
				assert.Equal(t, yesterday.Year(), startTime.Year())
				assert.Equal(t, yesterday.Month(), startTime.Month())
				assert.Equal(t, yesterday.Day(), startTime.Day())
				assert.Equal(t, 0, startTime.Hour())
				assert.Equal(t, 0, startTime.Minute())
				assert.Equal(t, 23, endTime.Hour())
				assert.Equal(t, 59, endTime.Minute())
			case "this month":
				assert.Equal(t, fixedTime.Year(), startTime.Year())
				assert.Equal(t, fixedTime.Month(), startTime.Month())
				assert.Equal(t, 1, startTime.Day())
			case "this year":
				assert.Equal(t, fixedTime.Year(), startTime.Year())
				assert.Equal(t, time.January, startTime.Month())
				assert.Equal(t, 1, startTime.Day())
			}
		})
	}
}

func TestParseRelativeTime(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		hasError bool
	}{
		{"1s", time.Second, false},
		{"60s", 60 * time.Second, false},
		{"5m", 5 * time.Minute, false},
		{"24h", 24 * time.Hour, false},
		{"7d", 7 * 24 * time.Hour, false},
		{"2w", 2 * 7 * 24 * time.Hour, false},
		{"3M", 3 * 30 * 24 * time.Hour, false},
		{"1y", 365 * 24 * time.Hour, false},
		{"invalid", 0, true},
		{"10x", 0, true}, // Invalid unit
		{"-5d", 0, true}, // Negative value will fail regex
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			duration, err := parseRelativeTime(tt.input)
			
			if tt.hasError {
				assert.Error(t, err)
				return
			}
			
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, duration)
		})
	}
}

func TestParseDateString(t *testing.T) {
	tests := []struct {
		input     string
		year      int
		month     time.Month
		day       int
		hasError  bool
	}{
		{"2023-01-15", 2023, time.January, 15, false},
		{"2023/01/15", 2023, time.January, 15, false},
		{"01/15/2023", 2023, time.January, 15, false},
		{"15-Jan-2023", 2023, time.January, 15, false},
		{"Jan 15, 2023", 2023, time.January, 15, false},
		{"2023-01-15T12:30:45Z", 2023, time.January, 15, false},
		{"Invalid date", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			date, err := parseDateString(tt.input)
			
			if tt.hasError {
				assert.Error(t, err)
				return
			}
			
			assert.NoError(t, err)
			assert.Equal(t, tt.year, date.Year())
			assert.Equal(t, tt.month, date.Month())
			assert.Equal(t, tt.day, date.Day())
		})
	}
}