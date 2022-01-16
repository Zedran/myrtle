package main

import (
	"log"
	"strconv"
	"strings"
)

/* Atof handles parsing of the wild formats of TLE floats. Any error is printed to log output. */
func Atof(s string, normalize bool) float64 {
	s = strings.Trim(s, " ")

	if normalize {
		s = NormalizeFloat(s)
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println(err)
	}

	return f
}

/* A more compact version of Atoi. Any error is printed to log output. */
func Atoi(s string) int {
	i, err := strconv.Atoi(strings.Trim(s, " "))
	if err != nil {
		log.Println(err)
	}

	return i
}

/* Expands NORAD classification abbreviations. */
func ExpandClass(c string) string {
	switch c {
	case "C":
		return "Classified"
	case "S":
		return "Secret"
	case "U":
		return "Unclassified"
	default:
		return "Unknown"
	}
}

/* Normalizes floating point numbers contained in TLE, as the decimal point and the exponent
 * of ten notation are often assumed. This function prepares a string for parsing into float.
 */
func NormalizeFloat(s string) string {
	// Add decimal point
	if !strings.Contains(s, ".") {
		if strings.HasPrefix(s, "-") {
			s = s[:1] + "." + s[1:]
		} else {
			s = "." + s
		}
	}

	// Add the exponent of ten notation
	liMinus := strings.LastIndex(s, "-")
	liPlus  := strings.LastIndex(s, "+")

	var li int

	if liMinus > 0 {
		li = liMinus
	} else if liPlus > 0 {
		li = liPlus
	} else {
		return s
	}

	return s[:li] + "e" + s[li:]
}

/* Since the year in international designator is represented as a two-digit number, this function
 * is needed to make sure a proper century is indicated. 
 * 57-99 == 1957-1999 
 * 00-56 == 2000-2056
 */
func NormalizeYear(y int) int {
	if y > 56 {
		return 1900 + y
	}
	return 2000 + y
}
