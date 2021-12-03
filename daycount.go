// Package daycount implements the most useful day counting conventions in finance
package daycount

import (
	"time"
)

// data is an internal data structure to simplify the function calls
type data struct {
	Date1       time.Time
	Date2       time.Time
	Date3       time.Time
	Compounding int
}

// Default day count convention
var Default string = "30E360"

// conventions is a map strcuture that contains the information
// to calculate the days between two dates and converts it into
// a day count fraction.
// https://www.isda.org/2008/12/22/30-360-day-count-conventions
var conventions = map[string]struct {
	Numerator   func(d data) float64
	Denominator func(d data) float64
}{
	// ISDA
	"30E360": {
		Numerator: func(d data) float64 {

			day1, day2 := d.Date1.Day(), d.Date2.Day()
			if day1 == 31 || isLastDayofFeb(d.Date1) {
				day1 = 30
			}
			// FIXME: if date2 is last day of Feb, we should ensure that date2 is not termination date
			if day2 == 31 || isLastDayofFeb(d.Date2) {
				day2 = 30
			}
			return days30360(d.Date1, d.Date2, day1, day2)
		},
		Denominator: func(d data) float64 { return 360.0 },
	},
	"EUROBOND": {
		Numerator: func(d data) float64 {
			day1, day2 := d.Date1.Day(), d.Date2.Day()
			if day1 == 31 {
				day1 = 30
			}
			if day2 == 31 {
				day2 = 30
			}
			return days30360(d.Date1, d.Date2, day1, day2)

		},
		Denominator: func(d data) float64 { return 360.0 },
	},
	// "30U360": {
	// 	Numerator: func(d data) float64 {
	// 		day1, day2 := d.Date1.Day(), d.Date2.Day()
	// 		if day2 == 31 && day1 >= 30 {
	// 			day2 = 30
	// 		}
	// 		if day1 == 31 {
	// 			day1 = 30
	// 		}
	// 		return days30360(d.Date1, d.Date2, day1, day2)
	// 	},
	// 	Denominator: func(d data) float64 { return 360.0 },
	// },
	"BONDBASIS": {
		Numerator: func(d data) float64 {
			day1, day2 := d.Date1.Day(), d.Date2.Day()
			if day1 == 31 {
				day1 = 30
			}
			if day2 == 31 && day1 >= 30 {
				day2 = 30
			}
			return days30360(d.Date1, d.Date2, day1, day2)

		},
		Denominator: func(d data) float64 { return 360.0 },
	},
	"ACT360": {
		Numerator: func(d data) float64 {
			return d.Date2.Sub(d.Date1).Hours() / 24.0
		},
		Denominator: func(d data) float64 {
			return 360.0
		},
	},
	"ACTACT": {
		Numerator: func(d data) float64 {
			return d.Date2.Sub(d.Date1).Hours() / 24.0
		},
		Denominator: func(d data) float64 {
			return float64(d.Compounding) * d.Date3.Sub(d.Date1).Hours() / 24.0
		},
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
func Fraction(date1, date2, date3 time.Time, compounding int, basis string) float64 {

	// create data struct
	d := data{
		date1,
		date2,
		date3,
		compounding,
	}

	// use default if basis is empty
	if basis == "" {
		basis = Default
	}

	// look for convention
	conv, ok := conventions[basis]
	if !ok {
		return 0.0
	}

	// calculate day count fraction
	return conv.Numerator(d) / conv.Denominator(d)
}

// Days counts the dates between two dates
func Days(date1, date2 time.Time, basis string) float64 {

	// create data struct
	d := data{
		date1,
		date2,
		time.Time{},
		1,
	}

	// use default if basis is empty
	if basis == "" {
		basis = Default
	}

	// look for convention
	conv, ok := conventions[basis]
	if !ok {
		return 0.0
	}

	// calculate days
	return conv.Numerator(d)
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
