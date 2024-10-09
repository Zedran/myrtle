package main

import (
	"fmt"
	"strings"
	"time"
)

// Elements struct holding orbital parameters derived from TLE.
type Elements struct {
	// The name of the described object
	Name,

	// Trimmed TLE lines
	L1 string
	L2 string

	// Epoch in unix seconds
	Epoch int64

	// Mass of dominant body
	DM float64

	// Radius of dominant body
	DR float64

	// Semi-Major Axis
	SMa float64

	// Semi-Minor Axis
	SMi float64

	// Periapsis Radius
	PeR float64

	// Apoapsis Radius
	ApR float64

	// Radius at Epoch
	R float64

	// Orbital Eccenticity
	Ecc float64

	// Orbital Period
	T float64

	// Time to Periapsis
	PeT float64

	// Time to Apoapsis
	ApT float64

	// Velocity at Epoch
	Vel float64

	// Orbital Inclination
	Inc float64

	// Longitude of Ascending Node
	LAN float64

	// Longitude of Periapsis
	LPe float64

	// Argument of Periapsis
	AgP float64

	// True Anomaly
	TrA float64

	// True Longitude
	TrL float64

	// Mean Anomaly
	MnA float64

	// Mean Longitude
	MnL float64

	// Eccentric Anomaly
	EcA float64

	// Indicates that eccentric anomaly solution did not converge
	EcAConvErr bool
}

// Creates a title string consisting of the object name, dates and original set lines.
func (e *Elements) GetTitle() string {
	mjd := JDNToMJD(UnixToJDN(e.Epoch))
	date := time.Unix(e.Epoch, 0).Format("2006-01-02T15:04:05 UTC")

	return fmt.Sprintf("%s    MJD %.5f    %4s\n    %4s\n    %4s\n\n", e.Name, mjd, date, e.L1, e.L2)
}

// Converts the Elements struct fields into a slice of strings. If alt is true,
// the distance will be displayed in relation to the dominant body's surface
// (ASL) instead of measuring it from the body's center. If acc is true,
// the numbers are not crunched by FormatNumber function. This increases
// their precision, but reduces readability.
func (e *Elements) ToString(alt, acc bool) []string {
	var (
		// These variables are different depending on the reference point. If alt is true,
		// then periapsis, apoapsis and radius are converted to altitude ASL. The deltaD
		// (difference in distance) is equal to dominant body's radius in that case.
		// The aforementioned three variables are also named differently.
		deltaD    float64
		pe, ap, r string

		// Eccentric anomaly has '!' appended to its symbol if solution for it does not converge.
		eca string = "EcA"
	)

	if e.EcAConvErr {
		eca += "!"
	}

	if alt {
		deltaD = e.DR

		// Altitudes
		pe = "PeA"
		ap = "ApA"
		r = "Alt"
	} else {
		deltaD = 0

		// Radii
		pe = "PeR"
		ap = "ApR"
		r = "R"
	}

	return []string{
		ParamToString("SMa", e.SMa, acc),
		ParamToString("SMi", e.SMi, acc),
		ParamToString(pe, e.PeR-deltaD, acc),
		ParamToString(ap, e.ApR-deltaD, acc),
		ParamToString(r, e.R-deltaD, acc),
		ParamToString("Ecc", e.Ecc, acc),
		ParamToString("T", e.T, acc),
		ParamToString("PeT", e.PeT, acc),
		ParamToString("ApT", e.ApT, acc),
		ParamToString("Vel", e.Vel, acc),
		ParamToString("Inc", e.Inc, acc),
		ParamToString("LAN", e.LAN, acc),
		ParamToString("LPe", e.LPe, acc),
		ParamToString("AgP", e.AgP, acc),
		ParamToString("TrA", e.TrA, acc),
		ParamToString("TrL", e.TrL, acc),
		ParamToString("MnA", e.MnA, acc),
		ParamToString("MnL", e.MnL, acc),
		ParamToString(eca, e.EcA, acc),
	}
}

// Creates Elements struct from TLE. Accepts dominant body mass and radius as well.
func CalculateElements(tle *TLE, m, r float64) *Elements {
	var (
		e   Elements
		err error
	)

	e.Name = strings.Trim(tle.Match.Title, " ")
	e.L1 = strings.Trim(tle.Match.Line1, " ")
	e.L2 = strings.Trim(tle.Match.Line2, " ")

	e.Epoch = EpochToUnix(tle.L1.Epoch.Year, tle.L1.Epoch.Day)

	e.DM = m
	e.DR = r

	e.Ecc = tle.L2.Ecc
	e.Inc = tle.L2.Inc
	e.LAN = tle.L2.LAN
	e.AgP = tle.L2.AgP
	e.MnA = tle.L2.MnA

	e.T = Period(tle.L2.MnM)
	e.SMa = SemiMajorAxis(e.T, e.DM)
	e.SMi = SemiMinorAxis(e.SMa, e.Ecc)
	e.PeR = PeriapsisRadius(e.SMa, e.Ecc)
	e.ApR = ApoapsisRadius(e.SMa, e.Ecc)

	e.PeT = TimeToPeriapsis(e.MnA, e.T)
	e.ApT = TimeToApoapsis(e.MnA, e.T, e.PeT)

	e.LPe = LongitudeOfPeriapsis(e.LAN, e.AgP)
	e.MnL = MeanLongitude(e.MnA, e.LPe)

	e.EcA, err = EccentricAnomaly(e.Ecc, e.MnA)
	if err != nil {
		e.EcAConvErr = true
	}

	e.TrA = TrueAnomaly(e.Ecc, e.EcA)
	e.TrL = TrueLongitude(e.TrA, e.LPe)
	e.R = OrbitalRadius(e.SMa, e.Ecc, e.TrA)
	e.Vel = OrbitalVelocity(e.R, e.SMa, e.DM)

	return &e
}

// Ensures the proper display format of the orbital element depending
// on its type. If accurate is true, the value is not submitted
// to FormatNumber function. This means it will be represented with maximum
// precision and reduced readability.
func ParamToString(symbol string, value float64, accurate bool) string {
	if accurate {
		return fmt.Sprintf("%-5s%f", symbol, value)
	}

	var n string

	switch symbol {
	case "SMa", "SMi", "PeR", "ApR", "R":
		n = FormatNumber(value, 5, 3, false, true)
	case "PeA", "ApA", "Alt":
		n = FormatNumber(value, 5, 1, false, true)
	case "Ecc":
		n = FormatNumber(value, 6, 4, false, false)
	case "T", "PeT", "ApT", "Vel":
		n = FormatNumber(value, 5, 3, false, true)
	default: // Angles
		n = FormatNumber(value, 6, 2, true, false)
	}

	return fmt.Sprintf("%-5s%s", symbol, n)
}
