// SPDX-FileCopyrightText: 2024 Weston Schmidt
// SPDX-License-Identifier: Apache-2.0

package civil_test

import (
	"fmt"
	"time"

	"github.com/schmidtw/civil"
)

func Example() {
	// Create a date directly
	d := civil.Date{Year: 2024, Month: time.March, Day: 15}
	fmt.Println(d)

	// Parse from a string (RFC 3339 full-date)
	d, _ = civil.ParseDate("2024-03-15")
	fmt.Println(d)

	// Parse from a uint32 (YYYYMMDD)
	d, _ = civil.ParseUint32Date(20240315)
	fmt.Println(d)

	// Extract from a time.Time
	d = civil.DateOf(time.Date(2024, time.March, 15, 12, 0, 0, 0, time.UTC))
	fmt.Println(d)

	// Arithmetic
	tomorrow := d.AddDays(1)
	fmt.Println(tomorrow)

	nextMonth := d.AddMonths(1)
	fmt.Println(nextMonth)

	nextYear := d.AddYears(1)
	fmt.Println(nextYear)

	// Difference
	fmt.Println(tomorrow.DaysSince(d))

	// Comparison
	fmt.Println(d.Before(tomorrow))
	fmt.Println(d.After(tomorrow))
	fmt.Println(d.Compare(tomorrow))

	// Convert back to time.Time
	t := d.In(time.UTC)
	fmt.Println(t)

	// Validation and inspection
	fmt.Println(d.IsValid())
	fmt.Println(d.IsZero())
	fmt.Println(d.Weekday())

	// Output:
	// 2024-03-15
	// 2024-03-15
	// 2024-03-15
	// 2024-03-15
	// 2024-03-16
	// 2024-04-15
	// 2025-03-15
	// 1
	// true
	// false
	// -1
	// 2024-03-15 00:00:00 +0000 UTC
	// true
	// false
	// Friday
}
