package main

/* TLE struct holds values extracted from the original two-line string set contained in Match. */
type TLE struct {
	// Original Match struct that underwent parsing into TLE struct
	Match *Match

	// Line1
	L1 struct {
		// Line number
		Number int

		// Satellite catalog number
		CatNum string

		// Classification: Classified, Secret, Unclassified
		Class string

		// International Designator
		IntlDesig struct {
			LaunchYear int    // Launch year
			LaunchNum  int    // Number of launch made on that year
			LaunchComp string // Launch component in order
		}

		Epoch struct {
			Year int
			Day  float64
		}

		// Mean Motion Derivatives
		MnMDrvs struct {
			First,
			Second float64
		}

		// B* drag term
		BSTAR float64

		EphemerisType int

		// Element set number (number of TLE sets generated for this object)
		ElSetNum int

		Checksum int
	}

	// Line 2
	L2 struct {
		// Line number
		Number int

		// Satellite catalog number
		CatNum string

		Inc, // Inclination
		LAN, // Longitude of Ascending Node
		Ecc, // Eccentricity
		AgP, // Argument of Periapsis
		MnA, // Mean Anomaly
		MnM float64 // Mean Motion

		// Revolutions number
		RevN int

		Checksum int
	}
}

/* Converts the Match struct into the TLE struct, splitting all the values.*/
func ParseMatch(m *Match) *TLE {
	var tle TLE

	// --------------- Original --------------- //

	tle.Match = m

	// ---------------- Line 1 ---------------- //

	tle.L1.Number = Atoi(tle.Match.Line1[:1])
	tle.L1.CatNum = tle.Match.Line1[2:7]
	tle.L1.Class = ExpandClass(tle.Match.Line1[7:8])

	tle.L1.IntlDesig.LaunchYear = NormalizeYear(Atoi(tle.Match.Line1[9:11]))
	tle.L1.IntlDesig.LaunchNum = Atoi(tle.Match.Line1[11:14])
	tle.L1.IntlDesig.LaunchComp = tle.Match.Line1[14:15]

	tle.L1.Epoch.Year = NormalizeYear(Atoi(tle.Match.Line1[18:20]))
	tle.L1.Epoch.Day = Atof(tle.Match.Line1[20:32], false)

	tle.L1.MnMDrvs.First = Atof(tle.Match.Line1[33:43], false)
	tle.L1.MnMDrvs.Second = Atof(tle.Match.Line1[45:52], true)
	tle.L1.BSTAR = Atof(tle.Match.Line1[53:61], true)

	tle.L1.EphemerisType = Atoi(tle.Match.Line1[62:63])
	tle.L1.ElSetNum = Atoi(tle.Match.Line1[65:68])

	tle.L1.Checksum = Atoi(tle.Match.Line1[68:69])

	// ---------------- Line 2 ---------------- //

	tle.L2.Number = Atoi(tle.Match.Line2[:1])
	tle.L2.CatNum = tle.Match.Line2[2:7]

	tle.L2.Inc = Atof(tle.Match.Line2[9:16], false)
	tle.L2.LAN = Atof(tle.Match.Line2[17:25], false)
	tle.L2.Ecc = Atof(tle.Match.Line2[26:33], true)
	tle.L2.AgP = Atof(tle.Match.Line2[34:42], false)
	tle.L2.MnA = Atof(tle.Match.Line2[43:51], false)
	tle.L2.MnM = Atof(tle.Match.Line2[52:63], false)

	tle.L2.RevN = Atoi(tle.Match.Line2[63:68])

	tle.L2.Checksum = Atoi(tle.Match.Line2[68:69])

	return &tle
}
