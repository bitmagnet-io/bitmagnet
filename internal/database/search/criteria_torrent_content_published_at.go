package search

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"gorm.io/gen/field"
)

// timeNow is a replaceable function for time.Now, making testing easier
var timeNow = time.Now

// TorrentContentPublishedAtCriteria returns a criteria that filters torrents by published_at timestamp
func TorrentContentPublishedAtCriteria(timeFrame string) query.Criteria {
	return query.DaoCriteria{
		Conditions: func(ctx query.DbContext) ([]field.Expr, error) {
			if timeFrame == "" {
				return nil, nil
			}
			
			startTime, endTime, err := parseTimeFrame(timeFrame)
			if err != nil {
				return nil, err
			}
			
			return []field.Expr{
				ctx.Query().TorrentContent.PublishedAt.Gte(startTime),
				ctx.Query().TorrentContent.PublishedAt.Lte(endTime),
			}, nil
		},
	}
}

// ParseTimeFrame parses a time frame string into start and end times
func parseTimeFrame(timeFrame string) (time.Time, time.Time, error) {
	timeFrame = strings.TrimSpace(timeFrame)
	
	// Default end time is now
	endTime := timeNow().UTC()
	var startTime time.Time
	
	// Empty string means no time filter
	if timeFrame == "" {
		return time.Time{}, time.Time{}, nil
	}
	
	// Handle relative time expressions (e.g., "3h", "7d")
	if relativeMatch, _ := regexp.MatchString(`^\d+[smhdwMy]$`, timeFrame); relativeMatch {
		duration, err := parseRelativeTime(timeFrame)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		startTime = endTime.Add(-duration)
		return startTime, endTime, nil
	}
	
	// Handle special expressions
	switch timeFrame {
	case "today":
		startTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, time.UTC)
		return startTime, endTime, nil
		
	case "yesterday":
		yesterday := endTime.AddDate(0, 0, -1)
		startTime = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
		endTime = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, time.UTC)
		return startTime, endTime, nil
		
	case "this week":
		// Calculate days since start of week (Monday)
		daysSinceMonday := int(endTime.Weekday())
		if daysSinceMonday == 0 { // Sunday
			daysSinceMonday = 6
		} else {
			daysSinceMonday--
		}
		startTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day()-daysSinceMonday, 0, 0, 0, 0, time.UTC)
		return startTime, endTime, nil
		
	case "last week":
		// Calculate days since start of week (Monday)
		daysSinceMonday := int(endTime.Weekday())
		if daysSinceMonday == 0 { // Sunday
			daysSinceMonday = 6
		} else {
			daysSinceMonday--
		}
		// Start of this week
		thisWeekStart := time.Date(endTime.Year(), endTime.Month(), endTime.Day()-daysSinceMonday, 0, 0, 0, 0, time.UTC)
		// Start of last week is 7 days before start of this week
		startTime = thisWeekStart.AddDate(0, 0, -7)
		// End of last week is 1 second before start of this week
		endTime = thisWeekStart.Add(-time.Second)
		return startTime, endTime, nil
		
	case "this month":
		startTime = time.Date(endTime.Year(), endTime.Month(), 1, 0, 0, 0, 0, time.UTC)
		return startTime, endTime, nil
		
	case "last month":
		// Start of this month
		thisMonthStart := time.Date(endTime.Year(), endTime.Month(), 1, 0, 0, 0, 0, time.UTC)
		// Start of last month
		startTime = thisMonthStart.AddDate(0, -1, 0)
		// End of last month is 1 second before start of this month
		endTime = thisMonthStart.Add(-time.Second)
		return startTime, endTime, nil
		
	case "this year":
		startTime = time.Date(endTime.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		return startTime, endTime, nil
		
	case "last year":
		// Start of this year
		thisYearStart := time.Date(endTime.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		// Start of last year
		startTime = thisYearStart.AddDate(-1, 0, 0)
		// End of last year is 1 second before start of this year
		endTime = thisYearStart.Add(-time.Second)
		return startTime, endTime, nil
	}
	
	// Try to parse as absolute date range (e.g., "2023-01-01 to 2023-01-31")
	if strings.Contains(timeFrame, " to ") {
		parts := strings.Split(timeFrame, " to ")
		if len(parts) != 2 {
			return time.Time{}, time.Time{}, errors.New("invalid date range format. Expected 'start to end'")
		}
		
		var err error
		startTime, err = parseDateString(strings.TrimSpace(parts[0]))
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		
		endTime, err = parseDateString(strings.TrimSpace(parts[1]))
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		
		// If end time doesn't have a time component, set it to end of day
		if endTime.Hour() == 0 && endTime.Minute() == 0 && endTime.Second() == 0 {
			endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, endTime.Location())
		}
		
		return startTime, endTime, nil
	}
	
	// Try to parse as a single date (e.g., "2023-01-01")
	parsedDate, err := parseDateString(timeFrame)
	if err == nil {
		startTime = parsedDate
		endTime = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, parsedDate.Location())
		return startTime, endTime, nil
	}
	
	return time.Time{}, time.Time{}, errors.New("could not parse time frame")
}

// parseRelativeTime parses a relative time string (e.g., "3h", "7d") into a time.Duration
func parseRelativeTime(relTime string) (time.Duration, error) {
	// Extract the number and unit
	re := regexp.MustCompile(`^(\d+)([smhdwMy])$`)
	matches := re.FindStringSubmatch(relTime)
	if len(matches) != 3 {
		return 0, errors.New("invalid relative time format. Expected format: '3h', '7d', etc.")
	}
	
	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}
	
	unit := matches[2]
	
	// Convert to duration
	switch unit {
	case "s": // seconds
		return time.Duration(value) * time.Second, nil
	case "m": // minutes
		return time.Duration(value) * time.Minute, nil
	case "h": // hours
		return time.Duration(value) * time.Hour, nil
	case "d": // days
		return time.Duration(value) * 24 * time.Hour, nil
	case "w": // weeks
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	case "M": // months (approximate)
		return time.Duration(value) * 30 * 24 * time.Hour, nil
	case "y": // years (approximate)
		return time.Duration(value) * 365 * 24 * time.Hour, nil
	default:
		return 0, errors.New("unknown time unit. Valid units: s, m, h, d, w, M, y")
	}
}

// parseDateString attempts to parse a date string in various formats
func parseDateString(dateStr string) (time.Time, error) {
	// Try standard formats
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"01/02/2006",
		"2-Jan-2006",
		"Jan 2, 2006",
	}
	
	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, errors.New("could not parse date string")
}