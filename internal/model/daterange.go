package model

import (
	"errors"
	"strings"
	"time"
)

type DateRange interface {
	Start() Date
	End() Date
	StartTime() time.Time
	EndTime() time.Time
}

func NewDateRangeFromYear(year Year) DateRange {
	return dateRangeYear{
		year: year,
	}
}

func NewDateRangeFromMonthAndYear(month time.Month, year Year) DateRange {
	return dateRangeMonthAndYear{
		month: month,
		year:  year,
	}
}

func NewDateRangeFromDates(start, end Date) DateRange {
	return dateRange{
		start: start,
		end:   end,
	}
}

func NewNilDateRange() DateRange {
	return dateRange{}
}

func NewDateRangeFromString(str string) (DateRange, error) {
	str = strings.ToLower(strings.TrimSpace(str))
	if str == "" {
		return NewNilDateRange(), nil
	}
	parts := strings.Split(str, " to ")
	switch len(parts) {
	case 1:
		date, err := NewDateFromIsoString(str)
		if err == nil {
			return date, nil
		}
		t, err := time.Parse("2006-01", str)
		if err == nil {
			return NewDateRangeFromMonthAndYear(t.Month(), Year(t.Year())), nil
		}
		t, err = time.Parse("2006", str)
		if err == nil {
			return NewDateRangeFromYear(Year(t.Year())), nil
		}
	case 2:
		start, err := NewDateRangeFromString(parts[0])
		if err == nil {
			end, err := NewDateRangeFromString(parts[1])
			if err == nil {
				return NewDateRangeFromDates(start.Start(), end.End()), nil
			}
		}
	}
	return nil, errors.New("invalid date range")
}

type dateRange struct {
	start, end Date
}

func (d dateRange) Start() Date {
	if d.start.IsNil() {
		return NewDateFromParts(1, time.January, 1)
	}
	return d.start
}

func (d dateRange) End() Date {
	if d.end.IsNil() {
		return NewDateFromParts(9999, time.December, 31)
	}
	return d.end
}

func (d dateRange) StartTime() time.Time {
	return d.Start().Time()
}

func (d dateRange) EndTime() time.Time {
	return d.End().EndOfDayTime()
}

type dateRangeYear struct {
	year Year
}

func (d dateRangeYear) Start() Date {
	return Date{
		Year:  d.year,
		Month: time.January,
		Day:   1,
	}
}

func (d dateRangeYear) End() Date {
	return Date{
		Year:  d.year,
		Month: time.December,
		Day:   31,
	}
}

func (d dateRangeYear) StartTime() time.Time {
	return d.Start().Time()
}

func (d dateRangeYear) EndTime() time.Time {
	return d.End().EndOfDayTime()
}

type dateRangeMonthAndYear struct {
	month time.Month
	year  Year
}

func (d dateRangeMonthAndYear) Start() Date {
	return Date{
		Year:  d.year,
		Month: d.month,
		Day:   1,
	}
}

func numDaysInMonth(year Year, month time.Month) uint8 {
	switch month {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return 31
	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		if year%4 == 0 {
			return 29
		} else {
			return 28
		}
	}
	return 0
}

func (d dateRangeMonthAndYear) End() Date {
	return Date{
		Year:  d.year,
		Month: d.month,
		Day:   numDaysInMonth(d.year, d.month),
	}
}

func (d dateRangeMonthAndYear) StartTime() time.Time {
	return d.Start().Time()
}

func (d dateRangeMonthAndYear) EndTime() time.Time {
	return d.End().EndOfDayTime()
}

func (d Date) Start() Date {
	return d
}

func (d Date) End() Date {
	return d
}

func (d Date) StartTime() time.Time {
	return d.Time()
}

func (d Date) EndTime() time.Time {
	return d.EndOfDayTime()
}
