// SPDX-FileCopyrightText: 2016 Google LLC
// SPDX-License-Identifier: Apache-2.0

package civil

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDates(t *testing.T) {
	for _, test := range []struct {
		date     Date
		loc      *time.Location
		wantStr  string
		wantTime time.Time
	}{
		{
			date:     Date{2014, 7, 29},
			loc:      time.Local,
			wantStr:  "2014-07-29",
			wantTime: time.Date(2014, time.July, 29, 0, 0, 0, 0, time.Local),
		},
		{
			date:     DateOf(time.Date(2014, 8, 20, 15, 8, 43, 1, time.Local)),
			loc:      time.UTC,
			wantStr:  "2014-08-20",
			wantTime: time.Date(2014, 8, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			date:     DateOf(time.Date(999, time.January, 26, 0, 0, 0, 0, time.Local)),
			loc:      time.UTC,
			wantStr:  "0999-01-26",
			wantTime: time.Date(999, 1, 26, 0, 0, 0, 0, time.UTC),
		},
	} {
		if got := test.date.String(); got != test.wantStr {
			t.Errorf("%#v.String() = %q, want %q", test.date, got, test.wantStr)
		}
		if got := test.date.In(test.loc); !got.Equal(test.wantTime) {
			t.Errorf("%#v.In(%v) = %v, want %v", test.date, test.loc, got, test.wantTime)
		}
	}
}

func TestDateIsValid(t *testing.T) {
	for _, test := range []struct {
		date Date
		want bool
	}{
		{Date{2014, 7, 29}, true},
		{Date{2000, 2, 29}, true},
		{Date{10000, 12, 31}, true},
		{Date{1, 1, 1}, true},
		{Date{0, 1, 1}, true},  // year zero is OK
		{Date{-1, 1, 1}, true}, // negative year is OK
		{Date{1, 0, 1}, false},
		{Date{1, 1, 0}, false},
		{Date{2016, 1, 32}, false},
		{Date{2016, 13, 1}, false},
		{Date{1, -1, 1}, false},
		{Date{1, 1, -1}, false},
	} {
		got := test.date.IsValid()
		if got != test.want {
			t.Errorf("%#v: got %t, want %t", test.date, got, test.want)
		}
	}
}

func TestParseDate(t *testing.T) {
	for _, test := range []struct {
		str  string
		want Date // if empty, expect an error
	}{
		{"2016-01-02", Date{2016, 1, 2}},
		{"2016-12-31", Date{2016, 12, 31}},
		{"0003-02-04", Date{3, 2, 4}},
		{"999-01-26", Date{}},
		{"", Date{}},
		{"2016-01-02x", Date{}},
	} {
		got, err := ParseDate(test.str)
		if got != test.want {
			t.Errorf("ParseDate(%q) = %+v, want %+v", test.str, got, test.want)
		}
		if err != nil && test.want != (Date{}) {
			t.Errorf("Unexpected error %v from ParseDate(%q)", err, test.str)
		}
	}
}

func TestParseUint32Date(t *testing.T) {
	for _, test := range []struct {
		desc    string
		input   uint32
		want    Date
		wantErr bool
	}{
		{
			desc:  "typical date",
			input: 20240315,
			want:  Date{2024, 3, 15},
		},
		{
			desc:  "first day of year",
			input: 20230101,
			want:  Date{2023, 1, 1},
		},
		{
			desc:  "last day of year",
			input: 20231231,
			want:  Date{2023, 12, 31},
		},
		{
			desc:  "zero value",
			input: 0,
			want:  Date{0, 0, 0},
			// month=0, day=0 — both fail validation
			wantErr: true,
		},
		{
			desc:    "month too high",
			input:   20231301,
			wantErr: true,
		},
		{
			desc:    "month zero",
			input:   20230015,
			wantErr: true,
		},
		{
			desc:    "day zero",
			input:   20230100,
			wantErr: true,
		},
		{
			desc:    "day too high",
			input:   20230132,
			wantErr: true,
		},
		{
			desc:  "leap day",
			input: 20000229,
			want:  Date{2000, 2, 29},
		},
		{
			desc:  "year before 1000",
			input: 9990101,
			want:  Date{999, 1, 1},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			got, err := ParseUint32Date(test.input)
			if test.wantErr {
				if err == nil {
					t.Errorf("ParseUint32Date(%d) = %#v, want error", test.input, got)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error %v from ParseUint32Date(%d)", err, test.input)
			}
			if got != test.want {
				t.Errorf("ParseUint32Date(%d) = %#v, want %#v", test.input, got, test.want)
			}
		})
	}
}

func TestDateArithmetic(t *testing.T) {
	for _, test := range []struct {
		desc  string
		start Date
		end   Date
		days  int
	}{
		{
			desc:  "zero days noop",
			start: Date{2014, 5, 9},
			end:   Date{2014, 5, 9},
			days:  0,
		},
		{
			desc:  "crossing a year boundary",
			start: Date{2014, 12, 31},
			end:   Date{2015, 1, 1},
			days:  1,
		},
		{
			desc:  "negative number of days",
			start: Date{2015, 1, 1},
			end:   Date{2014, 12, 31},
			days:  -1,
		},
		{
			desc:  "full leap year",
			start: Date{2004, 1, 1},
			end:   Date{2005, 1, 1},
			days:  366,
		},
		{
			desc:  "full non-leap year",
			start: Date{2001, 1, 1},
			end:   Date{2002, 1, 1},
			days:  365,
		},
		{
			desc:  "crossing a leap second",
			start: Date{1972, 6, 30},
			end:   Date{1972, 7, 1},
			days:  1,
		},
		{
			desc:  "dates before the unix epoch",
			start: Date{101, 1, 1},
			end:   Date{102, 1, 1},
			days:  365,
		},
	} {
		if got := test.start.AddDays(test.days); got != test.end {
			t.Errorf("[%s] %#v.AddDays(%v) = %#v, want %#v", test.desc, test.start, test.days, got, test.end)
		}
		if got := test.end.DaysSince(test.start); got != test.days {
			t.Errorf("[%s] %#v.Sub(%#v) = %v, want %v", test.desc, test.end, test.start, got, test.days)
		}
	}
}

func TestDateArithmeticMonths(t *testing.T) {
	for _, test := range []struct {
		desc   string
		start  Date
		end    Date
		months int
	}{
		{
			desc:   "zero months noop",
			start:  Date{2024, 12, 16},
			end:    Date{2024, 12, 16},
			months: 0,
		},
		{
			desc:   "crossing a year boundary",
			start:  Date{2014, 8, 31},
			end:    Date{2015, 1, 31},
			months: 5,
		},
		{
			desc:   "negative number of months",
			start:  Date{2015, 1, 1},
			end:    Date{2014, 12, 1},
			months: -1,
		},
		{
			desc:   "full leap year",
			start:  Date{2008, 1, 1},
			end:    Date{2009, 1, 1},
			months: 12,
		},
		{
			desc:   "full non-leap year",
			start:  Date{1997, 1, 1},
			end:    Date{1998, 1, 1},
			months: 12,
		},
		{
			desc:   "crossing a leap second",
			start:  Date{1972, 6, 30},
			end:    Date{1972, 7, 30},
			months: 1,
		},
		{
			desc:   "dates before the unix epoch",
			start:  Date{101, 1, 1},
			end:    Date{101, 6, 1},
			months: 5,
		},
	} {
		if got := test.start.AddMonths(test.months); got.Compare(test.end) != 0 {
			t.Errorf("[%s] %#v.AddMonths(%v) = %#v, want %#v", test.desc, test.start, test.months, got, test.end)
		}
	}
}

func TestDateArithmeticYears(t *testing.T) {
	for _, test := range []struct {
		desc  string
		start Date
		end   Date
		years int
	}{
		{
			desc:  "zero years noop",
			start: Date{2024, 6, 19},
			end:   Date{2024, 6, 19},
			years: 0,
		},
		{
			desc:  "positive number of years",
			start: Date{2012, 4, 29},
			end:   Date{2014, 4, 29},
			years: 2,
		},
		{
			desc:  "negative number of years",
			start: Date{2027, 1, 1},
			end:   Date{2024, 1, 1},
			years: -3,
		},
		{
			desc:  "crossing a leap second",
			start: Date{1972, 6, 30},
			end:   Date{1973, 6, 30},
			years: 1,
		},
		{
			desc:  "dates before the unix epoch",
			start: Date{99, 1, 1},
			end:   Date{102, 1, 1},
			years: 3,
		},
	} {
		if got := test.start.AddYears(test.years); got.Compare(test.end) != 0 {
			t.Errorf("[%s] %#v.AddDays(%v) = %#v, want %#v", test.desc, test.start, test.years, got, test.end)
		}
	}
}

func TestDateBefore(t *testing.T) {
	for _, test := range []struct {
		d1, d2 Date
		want   bool
	}{
		{Date{2016, 12, 31}, Date{2017, 1, 1}, true},
		{Date{2016, 1, 1}, Date{2016, 1, 1}, false},
		{Date{2016, 12, 30}, Date{2016, 12, 31}, true},
	} {
		if got := test.d1.Before(test.d2); got != test.want {
			t.Errorf("%v.Before(%v): got %t, want %t", test.d1, test.d2, got, test.want)
		}
	}
}

func TestDateAfter(t *testing.T) {
	for _, test := range []struct {
		d1, d2 Date
		want   bool
	}{
		{Date{2016, 12, 31}, Date{2017, 1, 1}, false},
		{Date{2016, 1, 1}, Date{2016, 1, 1}, false},
		{Date{2016, 12, 30}, Date{2016, 12, 31}, false},
	} {
		if got := test.d1.After(test.d2); got != test.want {
			t.Errorf("%v.After(%v): got %t, want %t", test.d1, test.d2, got, test.want)
		}
	}
}

func TestDateCompare(t *testing.T) {
	for _, test := range []struct {
		d1, d2 Date
		want   int
	}{
		{Date{2016, 12, 31}, Date{2017, 1, 1}, -1},
		{Date{2016, 1, 1}, Date{2016, 1, 1}, 0},
		{Date{2016, 12, 31}, Date{2016, 12, 30}, +1},
	} {
		if got := test.d1.Compare(test.d2); got != test.want {
			t.Errorf("%v.Compare(%v): got %d, want %d", test.d1, test.d2, got, test.want)
		}
	}
}

func TestDateIsZero(t *testing.T) {
	for _, test := range []struct {
		date Date
		want bool
	}{
		{Date{2000, 2, 29}, false},
		{Date{10000, 12, 31}, false},
		{Date{-1, 0, 0}, false},
		{Date{0, 0, 0}, true},
		{Date{}, true},
	} {
		got := test.date.IsZero()
		if got != test.want {
			t.Errorf("%#v: got %t, want %t", test.date, got, test.want)
		}
	}
}

func TestDateWeekday(t *testing.T) {
	for _, test := range []struct {
		date Date
		want time.Weekday
	}{
		{Date{2024, 12, 21}, time.Saturday},
		{Date{1900, 1, 1}, time.Monday},
		{Date{2482, 3, 17}, time.Tuesday},
	} {
		got := test.date.Weekday()
		if got != test.want {
			t.Errorf("%#v: got %v, want %v", test.date, got, test.want)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	for _, test := range []struct {
		value interface{}
		want  string
	}{
		{Date{1987, 4, 15}, `"1987-04-15"`},
		{Date{}, `""`}, // zero value marshals to empty string
	} {
		bgot, err := json.Marshal(test.value)
		if err != nil {
			t.Fatal(err)
		}
		if got := string(bgot); got != test.want {
			t.Errorf("%#v: got %s, want %s", test.value, got, test.want)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var d Date
	for _, test := range []struct {
		data string
		ptr  *Date
		want *Date
	}{
		{`"1987-04-15"`, &d, &Date{1987, 4, 15}},
		{`"1987-04-\u0031\u0035"`, &d, &Date{1987, 4, 15}},
		{`"0000-00-00"`, &d, &Date{}}, // zero-value sentinel round-trips to zero
		{`""`, &d, &Date{}},           // empty string round-trips to zero
	} {
		if err := json.Unmarshal([]byte(test.data), test.ptr); err != nil {
			t.Fatalf("%s: %v", test.data, err)
		}
		if *test.want != *test.ptr {
			t.Errorf("%s: got %#v, want %#v", test.data, *test.ptr, *test.want)
		}
	}

	for _, bad := range []string{"", `"bad"`, `"1987-04-15x"`,
		`19870415`,     // a JSON number
		`11987-04-15x`, // not a JSON string

	} {
		if json.Unmarshal([]byte(bad), &d) == nil {
			t.Errorf("%q, Date: got nil, want error", bad)
		}
	}
}
