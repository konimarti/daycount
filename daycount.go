// Package daycount implements the most useful day counting conventions in finance
package daycount

import (
	"fmt"
	"time"
)

// Default day count convention
var Default string = "30E360"

type DateDiffFunc func(date1, date2 time.Time) float64

// conventions is a map strcuture that contains the information
// to calculate the days between two dates and converts it into
// a day count fraction.
// https://www.isda.org/2008/12/22/30-360-day-count-conventions
var conventions = map[string]struct {
	Numerator   DateDiffFunc
	Denominator DateDiffFunc
}{
	// ISDA
	"30E360": {
		Numerator:   days30e360,
		Denominator: days30e360,
	},
	"EUROBOND": {
		Numerator:   eurobond,
		Denominator: eurobond,
	},
	"BONDBASIS": {
		Numerator:   bondbasis,
		Denominator: bondbasis,
	},
	"ACT360": {
		Numerator:   act,
		Denominator: days30e360,
	},
	"ACTACT": {
		Numerator:   act,
		Denominator: act,
	},
}

// Implemented returns a slice of strings of the implemented day count conventions
func Implemented() []string {
	list := []string{}
	for conv := range conventions {
		list = append(list, conv)
	}
	return list
}

// Fraction returns the fraction of coupon that has been accrued between date1 and date2
// date1: last coupon payment, starting date for interest accrual
// date2: date through which interest rate is being accrued (settlement dates for bonds)
// date3: next coupon payment
// compounding: compounding frequency per year
func Fraction(date1, date2, date3 time.Time, basis string) (float64, error) {

	// use default if basis is empty
	if basis == "" {
		basis = Default
	}

	// look for convention
	conv, ok := conventions[basis]
	if !ok {
		return 0.0, fmt.Errorf("day count convention %s not implemented", basis)
	}

	// calculate day count fraction
	return conv.Numerator(date1, date2) / conv.Denominator(date1, date3), nil
}

// Days counts the dates between two dates
func Days(date1, date2 time.Time, basis string) (float64, error) {

	// use default if basis is empty
	if basis == "" {
		basis = Default
	}

	// look for convention
	conv, ok := conventions[basis]
	if !ok {
		return 0.0, fmt.Errorf("day count convention %s not implemented", basis)
	}

	// calculate days
	return conv.Numerator(date1, date2), nil
}

// days30360 is the helper function to calculate the days between two dates for the 30/360 methods
func days30360(d1, d2 time.Time, day1, day2 int) float64 {
	return 360.0*float64(d2.Year()-d1.Year()) + 30.0*float64(d2.Month()-d1.Month()) + float64(day2-day1)
}

// isLastDayOfFeb checks if time is the last day of February
func isLastDayofFeb(d time.Time) bool {
	if d.Month() == 2 {
		if d.YearDay() == time.Date(d.Year(), 3, 0, 0, 0, 0, 0, d.Location()).YearDay() {
			return true
		}
	}
	return false
}

func days30e360(date1, date2 time.Time) float64 {
	day1, day2 := date1.Day(), date2.Day()
	if day1 == 31 || isLastDayofFeb(date1) {
		day1 = 30
	}
	// FIXME: if date2 is last day of Feb, we should ensure that date2 is not termination date
	if day2 == 31 || isLastDayofFeb(date2) {
		day2 = 30
	}
	return days30360(date1, date2, day1, day2)
}

func eurobond(date1, date2 time.Time) float64 {
	day1, day2 := date1.Day(), date2.Day()
	if day1 == 31 {
		day1 = 30
	}
	if day2 == 31 {
		day2 = 30
	}
	return days30360(date1, date2, day1, day2)
}
func bondbasis(date1, date2 time.Time) float64 {
	day1, day2 := date1.Day(), date2.Day()
	if day1 == 31 {
		day1 = 30
	}
	if day2 == 31 && day1 >= 30 {
		day2 = 30
	}
	return days30360(date1, date2, day1, day2)
}
func act(date1, date2 time.Time) float64 {
	return date2.Sub(date1).Hours() / 24.0
}
