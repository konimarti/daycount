# Day Count Conventions in Golang

[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/konimarti/daycount/blob/main/LICENSE)
[![GoDoc](https://godoc.org/github.com/konimarti/observer?status.svg)](https://godoc.org/github.com/konimarti/daycount)
[![goreportcard](https://goreportcard.com/badge/github.com/konimarti/observer)](https://goreportcard.com/report/github.com/konimarti/dayount)

Implements different day counting conventions for finance applications.

The following counting conventions are implemented:

- 30E/360
- Act/360
- Act/Act
- Eurobond
- Bondbasis

This package implements two public functions `Days(date1, date2, convention)` and
`Fraction(date1, date2, date3, convention)`

- `Days` returns the number of days between two dates (date2 - date1) for the given counting convention
- `Fraction` returns the fraction of days between two dates ((date2-date1)/(date3-date1))
