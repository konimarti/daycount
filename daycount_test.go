package daycount_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/konimarti/daycount"
)

func TestDayCountFraction(t *testing.T) {

	// define tests
	testData := []struct {
		Date1       time.Time
		Date2       time.Time
		Date3       time.Time
		Compounding int
		Conventions []string
		Expected    []float64
	}{
		//
		//  awk '{print "{Date1: time.Date(20"$1",0,0,0,0,time.UTC), Date2: time.Date(20"$2",0,0,0,0,time.UTC), Date3: time.Time{}, Compounding:1, Conventions: []string{\"BONDBASIS\",\"EUROBOND\",\"30E360\",}, Expected: []")")}"}' out
		//
		{Date1: time.Date(2007, 01, 15, 0, 0, 0, 0, time.UTC), Date2: time.Date(2007, 01, 30, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.0416667, 0.0416667, 0.0416667}},
		{Date1: time.Date(2007, 01, 15, 0, 0, 0, 0, time.UTC), Date2: time.Date(2007, 02, 15, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.0833333, 0.0833333, 0.0833333}},
		{Date1: time.Date(2007, 01, 15, 0, 0, 0, 0, time.UTC), Date2: time.Date(2007, 07, 15, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.5, 0.5, 0.5}},
		{Date1: time.Date(2007, 9, 30, 0, 0, 0, 0, time.UTC), Date2: time.Date(2008, 03, 31, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.5, 0.5, 0.5}},
		{Date1: time.Date(2007, 9, 30, 0, 0, 0, 0, time.UTC), Date2: time.Date(2007, 10, 31, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.0833333, 0.0833333, 0.0833333}},
		{Date1: time.Date(2007, 9, 30, 0, 0, 0, 0, time.UTC), Date2: time.Date(2008, 9, 30, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{1, 1, 1}},
		{Date1: time.Date(2007, 01, 15, 0, 0, 0, 0, time.UTC), Date2: time.Date(2007, 01, 31, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.0444444, 0.0416667, 0.0416667}},
		{Date1: time.Date(2007, 01, 31, 0, 0, 0, 0, time.UTC), Date2: time.Date(2007, 02, 28, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.0777778, 0.0777778, 0.0833333}},
		{Date1: time.Date(2007, 02, 28, 0, 0, 0, 0, time.UTC), Date2: time.Date(2007, 03, 31, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.0916667, 0.0888889, 0.0833333}},
		{Date1: time.Date(2006, 8, 31, 0, 0, 0, 0, time.UTC), Date2: time.Date(2007, 02, 28, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.494444, 0.494444, 0.5}},
		{Date1: time.Date(2007, 02, 28, 0, 0, 0, 0, time.UTC), Date2: time.Date(2007, 8, 31, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.508333, 0.505556, 0.5}},
		{Date1: time.Date(2007, 02, 14, 0, 0, 0, 0, time.UTC), Date2: time.Date(2007, 02, 28, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.0388889, 0.0388889, 0.0444444}},
		{Date1: time.Date(2007, 02, 26, 0, 0, 0, 0, time.UTC), Date2: time.Date(2008, 02, 29, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{1.00833, 1.00833, 1.01111}},
		// does not work because of missing termination date:
		//{Date1: time.Date(2008, 02, 29, 0, 0, 0, 0, time.UTC), Date2: time.Date(2009, 02, 28, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.997222, 0.997222, 0.994444}},
		{Date1: time.Date(2008, 02, 29, 0, 0, 0, 0, time.UTC), Date2: time.Date(2008, 03, 30, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.0861111, 0.0861111, 0.0833333}},
		{Date1: time.Date(2008, 02, 29, 0, 0, 0, 0, time.UTC), Date2: time.Date(2008, 03, 31, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.0888889, 0.0861111, 0.0833333}},
		{Date1: time.Date(2007, 02, 28, 0, 0, 0, 0, time.UTC), Date2: time.Date(2007, 03, 05, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.0194444, 0.0194444, 0.0138889}},
		{Date1: time.Date(2007, 10, 31, 0, 0, 0, 0, time.UTC), Date2: time.Date(2007, 11, 28, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.0777778, 0.0777778, 0.0777778}},
		{Date1: time.Date(2007, 8, 31, 0, 0, 0, 0, time.UTC), Date2: time.Date(2008, 02, 29, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.497222, 0.497222, 0.5}},
		{Date1: time.Date(2008, 02, 29, 0, 0, 0, 0, time.UTC), Date2: time.Date(2008, 8, 31, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.505556, 0.502778, 0.5}},
		// does not work because of missing termination date:
		//{Date1: time.Date(2008, 8, 31, 0, 0, 0, 0, time.UTC), Date2: time.Date(2009, 02, 28, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.494444, 0.494444, 0.494444}},
		{Date1: time.Date(2009, 02, 28, 0, 0, 0, 0, time.UTC), Date2: time.Date(2009, 8, 31, 0, 0, 0, 0, time.UTC), Date3: time.Time{}, Compounding: 1, Conventions: []string{"BONDBASIS", "EUROBOND", "30E360"}, Expected: []float64{0.508333, 0.505556, 0.5}},
	}

	// run tests
	tolerance := 0.00001
	for nr, test := range testData {
		for i, conv := range test.Conventions {
			//frac, err := daycount.Fraction(test.Date1, test.Date2, test.Date1.AddDate(1, 0, 0), conv)
			days, err := daycount.Days(test.Date1, test.Date2, conv)
			frac := days / 360.0
			if math.Abs(frac-test.Expected[i]) > tolerance || err != nil {
				t.Errorf("test %d for %s failed, got: %f, want: %f\n", nr, conv, frac, test.Expected[i])
				fmt.Println(test)
			}
		}
	}
}

func TestDayCountFraction_MissingConvention(t *testing.T) {
	date1 := time.Date(2008, 02, 29, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2008, 8, 31, 0, 0, 0, 0, time.UTC)
	date3 := time.Time{}
	_, err := daycount.Fraction(date1, date2, date3, "THISISNOTIMPLEMENTED")
	if err == nil {
		t.Errorf("day count fraction should return an error when basis is not implemented")
	}
	_, err = daycount.Days(date1, date2, "THISISNOTIMPLEMENTED")
	if err == nil {
		t.Errorf("day count fraction should return an error when basis is not implemented")
	}
}
