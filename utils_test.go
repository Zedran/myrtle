package main

import "testing"

/* Test case data structure for TestFormatNumber. */
type formatNumberCase struct {
	// Tested number
	Number           float64

	// Correct output
	CorString,

	// Actual output
	Output           string

	// FormatNumber function parameters
	leftPadding,
	precision        int
	isAngle,
	adjustPrecision  bool

	// Positive test completion indicator
	Passed           bool
}

/* Runs the test and compares its output with the expected string. */
func (tc *formatNumberCase) run() {
	tc.Output = FormatNumber(tc.Number, tc.leftPadding, tc.precision, tc.isAngle, tc.adjustPrecision)
	tc.Passed = (tc.Output == tc.CorString)
}

/* A constructor for formatNumberCase struct. */
func createCase(n float64, corOut string, pad, prec int, isAngle, adjustPrec bool) *formatNumberCase {
	return &formatNumberCase{n, corOut, "", pad, prec, isAngle, adjustPrec, false}
}

/* Test for FormatNumber function. Every output must be exactly the same the expected string. */
func TestFormatNumber(t *testing.T) {
	// A degree symbol in unicode
	const deg = "\u00b0"

	cases := []*formatNumberCase{
		// Elements related to radius, time and velocity
		createCase(-4.125e6, "-4.12M", 5, 3, false, true),
		createCase( 6.717e6, "6.717M", 5, 3, false, true),

		// Elements related to altitude ASL
		createCase( 45.90e3, "45.90k", 5, 1, false, true),
		createCase( 646.5e3, "646.5k", 5, 1, false, true),
		createCase(-45.90e3, "-45.9k", 5, 1, false, true),
		createCase(-646.5e3, " -646k", 5, 1, false, true),

		// Orbital eccentricity
		createCase(0.0000, "0.0000", 6, 4, false, false),
		createCase(0.0447, "0.0447", 6, 4, false, false),

		// Angles
		createCase( 48.78, " 48.78" + deg, 6, 2, true, false),
		createCase(210.30, "210.30" + deg, 6, 2, true, false),
		createCase(  2.77, "  2.77" + deg, 6, 2, true, false),
	}

	passed := true
	for i := range cases {
		cases[i].run()
		passed = passed && cases[i].Passed
	}

	if !passed {
		for i := range cases {
			if !cases[i].Passed {
				t.Logf("'%f' : '%s'\n", cases[i].Number, cases[i].Output)
			}
		}
		t.Fatal("Test failed for cases listed above.")
	}
}

/* Tests whether NormalizeFloat function works correctly. Assumed decimal point and the exponent of ten
 * notation ('e') must be inserted properly for the test to complete.
 */
func TestNormalizeFloat(t *testing.T) {
	cases := map[string]string{
		"-.00002182" : "-.00002182",
		"00000-0"    : ".00000e-0",
		"-11606-4"   : "-.11606e-4",
		"0006703"    : ".0006703",
		"-11606+4"   : "-.11606e+4",
	}

	failed := make(map[string]string)

	for k, v := range cases {
		out := NormalizeFloat(k)
		if out != v {
			failed[k] = out
		}
	}

	if len(failed) != 0 {
		t.Fatalf("Failed for: %#v", failed)
	}
}
