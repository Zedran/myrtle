package main

import "testing"

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