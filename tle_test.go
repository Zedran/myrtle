package main

import "testing"

// Tests ParseLines function, feeding it with data. Output must be equal to
// the specified values.
func TestParseLines(t *testing.T) {
	// Sample Match struct with TLE set downloaded from the API
	m := Match{
		Title: "ISS (ZARYA)",
		Line1: "1 25544U 98067A   22014.20078024 -.00001581  00000+0 -20061-4 0  9991",
		Line2: "2 25544  51.6452  19.1428 0006828  17.5887  10.3753 15.49476744321309",
	}

	tle := ParseMatch(&m)

	pass1 :=
		tle.L1.Number == 1 &&
			tle.L1.CatNum == "25544" &&
			tle.L1.Class == "Unclassified" &&
			tle.L1.IntlDesig.LaunchYear == 1998 &&
			tle.L1.IntlDesig.LaunchNum == 67 &&
			tle.L1.IntlDesig.LaunchComp == "A" &&
			tle.L1.Epoch.Year == 2022 &&
			tle.L1.Epoch.Day == 14.20078024 &&
			tle.L1.MnMDrvs.First == -0.00001581 &&
			tle.L1.MnMDrvs.Second == 0 &&
			tle.L1.BSTAR == -0.20061e-4 &&
			tle.L1.EphemerisType == 0 &&
			tle.L1.ElSetNum == 999 &&
			tle.L1.Checksum == 1

	pass2 :=
		tle.L2.Number == 2 &&
			tle.L2.CatNum == "25544" &&
			tle.L2.Inc == 51.6452 &&
			tle.L2.LAN == 19.1428 &&
			tle.L2.Ecc == 0.0006828 &&
			tle.L2.AgP == 17.5887 &&
			tle.L2.MnA == 10.3753 &&
			tle.L2.MnM == 15.49476744 &&
			tle.L2.RevN == 32130 &&
			tle.L2.Checksum == 9

	var fatalMsg string

	if !pass1 {
		fatalMsg = "Parsing failed at Line 1"
	}

	if !pass2 {
		if len(fatalMsg) > 0 {
			fatalMsg += " and Line 2"
		} else {
			fatalMsg = "Parsing failed at Line 2"
		}
	}

	if !(pass1 && pass2) {
		t.Fatal(fatalMsg)
	}
}
