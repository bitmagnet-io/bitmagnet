package model

import (
	"database/sql/driver"
	"errors"
	"io"
	"time"
)

type Date struct {
	Year  Year
	Month time.Month
	Day   uint8
}

func NewDateFromTime(t time.Time) Date {
	return Date{
		Year:  Year(t.Year()),
		Month: t.Month(),
		Day:   uint8(t.Day()),
	}
}

func NewDateFromParts(year Year, month time.Month, day uint8) Date {
	return Date{
		Year:  year,
		Month: month,
		Day:   day,
	}
}

func NewDateFromIsoString(str string) (Date, error) {
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return Date{}, err
	}

	return NewDateFromTime(t), nil
}

func (d Date) Time() time.Time {
	nsec := 0
	// avoid returning a zero value
	if d.Year == 1 && d.Month == time.January && d.Day == 1 {
		nsec = 1
	}

	return time.Date(int(d.Year), d.Month, int(d.Day), 0, 0, 0, nsec, time.UTC)
}

func (d Date) EndOfDayTime() time.Time {
	return time.Date(int(d.Year), d.Month, int(d.Day), 24, 0, 0, 0, time.UTC)
}

func (d Date) IsNil() bool {
	return d == Date{}
}

func (d Date) IsValid() bool {
	return d.Year >= 1000 &&
		d.Year <= 9999 &&
		d.Month >= 1 &&
		d.Month <= 12 &&
		d.Day >= 1 &&
		d.Day <= numDaysInMonth(d.Year, d.Month)
}

func (d *Date) Scan(value interface{}) error {
	t, tOk := value.(time.Time)
	if tOk {
		if !t.IsZero() {
			*d = NewDateFromTime(t)
		}

		return nil
	}

	return nil
}

func (d Date) Value() (driver.Value, error) {
	if d.IsNil() {
		//nolint:nilnil
		return nil, nil
	}

	return d.Time(), nil
}

func (d Date) IsoDateString() string {
	return d.Time().Format("2006-01-02")
}

func (d Date) YearString() string {
	return d.Time().Format("2006")
}

func (d *Date) UnmarshalGQL(v interface{}) error {
	str, strOk := v.(string)
	if !strOk {
		return errors.New("must be a string")
	}

	if len(str) > 0 {
		nD, dErr := NewDateFromIsoString(str)
		if dErr != nil {
			return dErr
		}

		*d = nD
	}

	return nil
}

func (d Date) MarshalGQL(w io.Writer) {
	if d.IsNil() {
		_, _ = w.Write([]byte("null"))
		return
	}

	_, _ = w.Write([]byte(`"` + d.IsoDateString() + `"`))
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
		}

		return 28
	}

	return 0
}
