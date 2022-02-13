package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

/* Atof handles parsing of the wild formats of TLE floats. Any error is printed to the log file. */
func Atof(s string, normalize bool) float64 {
	s = strings.Trim(s, " ")

	if normalize {
		s = NormalizeFloat(s)
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		Log(err)
	}

	return f
}

/* A more compact version of Atoi. Any error is printed to the log file. */
func Atoi(s string) int {
	i, err := strconv.Atoi(strings.Trim(s, " "))
	if err != nil {
		Log(err)
	}

	return i
}

/* Returns true if sequence is inside slice s. */
func Contains(s []string, seq string) bool {
	for i := range s {
		if s[i] == seq {
			return true
		}
	}
	return false
}

/* Looks for part of a string inside a slice index. Returns index or -1 if nothing is found. */
func ContainsPart(s []string, seq string) int {
	for i := range s {
		if strings.Contains(s[i], seq) {
			return i
		}
	}

	return -1
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

/* Formats the number n. leftPadding and precision are parameters for sprintf ensuring a proper placement 
 * of the number. isAngle causes the function to treat the number as an angle - a degree sign is appended 
 * and the reduction is omitted, since the angles never reach 1e3. adjustPrecision is a parameter that 
 * corrects the length of the fractional part to ensure the proper alignment of all values. It is needed 
 * to trim the number to a specific width, regardless of whether the negation sign is present or the 
 * integral part's digit count. It should be set to false for eccentricity and angular values, 
 * since they have a set precision and are never negative. 
 */
func FormatNumber(n float64, leftPadding, precision int, isAngle, adjustPrecision bool) string {
	const (
		// A degree symbol in unicode
		deg   string  = "\u00b0"

		div   float64 = 1e3

		templ string  = "%%%d.%df%%s"
	)

	// Prefixes indicating the order of magnitude
	var pfx [9]string = [9]string{"", "k", "M", "G", "T", "P", "E", "Z", "Y"}
	
	if isAngle {
		format := fmt.Sprintf(templ, leftPadding, precision)
		return fmt.Sprintf(format, n, deg)
	}

	var sign float64
	if n < 0 {
		n = math.Abs(n)
		sign = -1
	} else {
		sign = 1
	}

	var i int
	for i = 0; i < len(pfx) && n > div; i++ {
		n /= div
	}

	// Adjusts precision according to the number of digits in a number
	if adjustPrecision {
		switch {
		case n < 10:
			precision = 3
		case n < 100:
			precision = 2
		default:
			precision = 1
		}
		
		// The precision is reduced if the minus symbol occupies one place
		if sign < 0 {
			precision--
		}
	}

	format := fmt.Sprintf(templ, leftPadding, precision)

	return fmt.Sprintf(format, n * sign, pfx[i])
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

/* Rewrites the passed slice, omitting duplicate values. */
func RemoveDuplicates(s []string) []string {
	clean := make([]string, 0, len(s))

	for i := range s {
		if !Contains(clean, s[i]) {
			clean = append(clean, s[i])
		}
	}

	return clean
}

/* Converts radians to degrees. */
func Deg(rad float64) float64 {
	return rad * 180 / math.Pi
}

/* Converts degrees to radians. */
func Rad(deg float64) float64 {
	return deg * math.Pi / 180
}

/* Converts Unix Time to Julian Day Number. */
func UnixToJDN(unixSeconds int64) float64 {
	return float64(unixSeconds) / 86400 + 2440587.5
}

/* Converts Julian Day Number to Modified Julian Date. */
func JDNToMJD(jdn float64) float64 {
	return jdn - 2400000.5
}
