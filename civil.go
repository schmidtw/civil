// SPDX-FileCopyrightText: 2016 Google LLC
// SPDX-License-Identifier: Apache-2.0

// Package civil implements types for civil time, a time-zone-independent
// representation of time that follows the rules of the proleptic Gregorian
// calendar with exactly 24-hour days.
package civil

import (
	"fmt"
	"time"
)

// A Date represents a date (year, month, day).
//
// This type does not include location information, and therefore does not
// describe a unique 24-hour timespan.
type Date struct {
	Year  int        // Year (e.g., 2014).
	Month time.Month // Month of the year (January = 1, ...).
	Day   int        // Day of the month, starting at 1.
}

// DateOf returns the Date in which a time occurs in that time's location.
func DateOf(t time.Time) Date {
	var d Date
	d.Year, d.Month, d.Day = t.Date()
	return d
}

// ParseDate parses a string in RFC3339 full-date format and returns the date value it represents.
func ParseDate(s string) (Date, error) {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return Date{}, err
	}
	return DateOf(t), nil
}

var daysInMonth = map[time.Month]int{
	time.January:   31,
	time.February:  29, // max days
	time.March:     31,
	time.April:     30,
	time.May:       31,
	time.June:      30,
	time.July:      31,
	time.August:    31,
	time.September: 30,
	time.October:   31,
	time.November:  30,
	time.December:  31,
}

// ParseUint32Date parses a uint32(YYYYMMDD) format date and returns the date
// value it represents.
//
// A very simple check to see if the date is somewhat real is done.  If you want
// to validate it completely, call IsValid().
func ParseUint32Date(v uint32) (Date, error) {
	d := Date{
		Year:  int(v / 10000),
		Month: time.Month((v % 10000) / 100),
		Day:   int(v % 100),
	}

	days, ok := daysInMonth[d.Month]
	if !ok {
		return Date{}, fmt.Errorf("invalid date %d: month %d out of range", v, d.Month)
	}
	if d.Day < 1 || d.Day > days {
		return Date{}, fmt.Errorf("invalid date %d: day %d out of range for %s %d", v, d.Day, d.Month, d.Year)
	}

	return d, nil
}

// String returns the date in RFC3339 full-date format.
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// IsValid reports whether the date is valid.
func (d Date) IsValid() bool {
	return DateOf(d.In(time.UTC)) == d
}

// In returns the time corresponding to time 00:00:00 of the date in the location.
//
// In is always consistent with time.Date, even when time.Date returns a time
// on a different day. For example, if loc is America/Indiana/Vincennes, then both
//
//	time.Date(1955, time.May, 1, 0, 0, 0, 0, loc)
//
// and
//
//	civil.Date{Year: 1955, Month: time.May, Day: 1}.In(loc)
//
// return 23:00:00 on April 30, 1955.
//
// In panics if loc is nil.
func (d Date) In(loc *time.Location) time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, loc)
}

// AddDays returns the date that is n days in the future.
// n can also be negative to go into the past.
func (d Date) AddDays(n int) Date {
	return DateOf(d.In(time.UTC).AddDate(0, 0, n))
}

// AddMonths returns the date that is n months in the future.
// n can also be negative to go into the past.
func (d Date) AddMonths(n int) Date {
	return DateOf(d.In(time.UTC).AddDate(0, n, 0))
}

// AddYears returns the date that is n years in the future.
// n can also be negative to go into the past.
func (d Date) AddYears(n int) Date {
	return DateOf(d.In(time.UTC).AddDate(n, 0, 0))
}

// DaysSince returns the signed number of days between the date and s, not including the end day.
// This is the inverse operation to AddDays.
func (d Date) DaysSince(s Date) (days int) {
	// We convert to Unix time so we do not have to worry about leap seconds:
	// Unix time increases by exactly 86400 seconds per day.
	deltaUnix := d.In(time.UTC).Unix() - s.In(time.UTC).Unix()
	return int(deltaUnix / 86400)
}

// Before reports whether d occurs before d2.
func (d Date) Before(d2 Date) bool {
	if d.Year != d2.Year {
		return d.Year < d2.Year
	}
	if d.Month != d2.Month {
		return d.Month < d2.Month
	}
	return d.Day < d2.Day
}

// After reports whether d occurs after d2.
func (d Date) After(d2 Date) bool {
	return d2.Before(d)
}

// Compare compares d and d2. If d is before d2, it returns -1;
// if d is after d2, it returns +1; otherwise it returns 0.
func (d Date) Compare(d2 Date) int {
	if d.Before(d2) {
		return -1
	} else if d.After(d2) {
		return +1
	}
	return 0
}

// IsZero reports whether date fields are set to their default value.
func (d Date) IsZero() bool {
	return (d.Year == 0) && (int(d.Month) == 0) && (d.Day == 0)
}

// Weekday returns the day of the week for the date.
func (d Date) Weekday() time.Weekday {
	return d.In(time.UTC).Weekday()
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of d.String().
func (d Date) MarshalText() ([]byte, error) {
	if d.IsZero() {
		return nil, nil
	}

	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be a string in a format accepted by ParseDate.
func (d *Date) UnmarshalText(data []byte) error {
	if len(data) == 0 || string(data) == "0000-00-00" {
		*d = Date{}
		return nil
	}

	var err error
	*d, err = ParseDate(string(data))
	return err
}
