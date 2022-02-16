/* All functions within this file accept and return angular values in degrees unless explicitly
 * stated otherwise. In all other cases, SI units are used.
 *  
 * The eccentric anomaly solution consist of 3 functions:
 *
 * 1. EccentricAnomaly
 * 2. keplerStart
 * 3. eps
 *
 * These functions were adapted from: 
 *
 * Marc A. Murison. 2006. A Practical Method for Solving the Kepler Equation. [on-line]
 * Available at http://murison.alpheratz.net/dynamics/twobody/KeplerIterations_summary.pdf
 * [accessed on 15.01.2022] U.S. Naval Observatory, Washington, DC.
 *
 */

package main

import (
	"errors"
	"math"
	"time"
)

const (
	// Gravitational constant
	G   float64 = 6.67430e-11

	// Mass of Earth
	M_E float64 = 5.97219e+24

	// Mean radius of Earth
	R_E float64 = 6.371008e+6
)

// Error returned if eccentric anomaly solution does not converge after 100 iterations.
var errNoConvergence = errors.New("Kepler's equation solution failed to converge.")

/* Calculates orbital period from mean motion. */
func Period(mnm float64) float64 {
	return 86400 / mnm
}

/* Calculates semi-major axis from orbital period and dominant body mass. */
func SemiMajorAxis(t, dominantMass float64) float64 {
	return math.Cbrt((G * dominantMass * math.Pow(t, 2)) / (4 * math.Pow(math.Pi, 2)))
}

/* Returns semi-minor axis given semi-major-axis and orbital eccentricity. */
func SemiMinorAxis(sma, ecc float64) float64 {
	return sma * math.Sqrt(1 - math.Pow(ecc, 2))
}

/* Converts argument of periapsis to longitude of periapsis given longitude of ascending node. */
func LongitudeOfPeriapsis(lan, agp float64) float64 {
	return math.Mod(lan + agp, 360)
}

/* Returns the radius of apoapsis calculated from semi-major axis and eccentricity. */
func ApoapsisRadius(sma, ecc float64) float64 {
	return sma * (1 + ecc)
}

/* Returns the radius of periapsis calculated from semi-major axis and eccentricity. */
func PeriapsisRadius(sma, ecc float64) float64 {
	return sma * (1 - ecc)
}

/* Calculates time to periapsis from mean anomaly and orbital period. */
func TimeToPeriapsis(mna, t float64) float64 {
	return t - (mna / sweepRate(t))
}

/* Calculates time to apoapsis from mean anomaly, orbital period and time to periapsis. */
func TimeToApoapsis(mna, t, pet float64) float64 {
	if mna <= 180 {
		return pet - t / 2
	}
	return pet + t / 2
}

/* Calculates mean longitude given mean anomaly and longitude of periapsis. */
func MeanLongitude(mna, lpe float64) float64 {
	return math.Mod(lpe + mna, 360)
}

/* Eccentric anomaly solution adapted from:
 *
 * Marc A. Murison. 2006. A Practical Method for Solving the Kepler Equation.
 * For the full source description, see the top of the file or README.
 *
 * The function returns an error if the solution is not within a tolerance limits after
 * 100 iterations.
 *
 * Accepts orbital eccentricity and mean anomaly as arguments.
 */
func EccentricAnomaly(ecc, mna float64) (float64, error) {
	const (
		maxIter   int     = 100
		tolerance float64 = 1.0e-14
	)

	mnaNorm := math.Mod(Rad(mna), 2 * math.Pi)

	// Starting value of the eccentric anomaly
	eca0    := keplerStart(ecc, mnaNorm)

	// The difference in eccentric anomaly value between iterations
	dE      := tolerance + 1

	var eca float64
	
	for i := 0; dE > tolerance; i++ {
		if i >= maxIter {
			return Deg(eca), errNoConvergence
		}

		eca  = eca0 - eps(ecc, mnaNorm, eca0)
		dE   = math.Abs(eca - eca0)
		eca0 = eca
	}

	return Deg(eca), nil
}

/* Calculates the starting value for eccentric anomaly solution (Murison, 2006). 
 * Intakes orbital eccentricity and mean anomaly [rad].
 */
func keplerStart(ecc, mna float64) float64 {
	t34 := math.Pow(ecc, 2)
	t35 := ecc * t34
	t33 := math.Cos(mna)

	return mna + (-0.5 * t35 + ecc + (t34 + 1.5 * t33 * t35) * t33) * math.Sin(mna)
}

/* Iteration function for the eccentric anomaly solution (Murison, 2006). Accepts orbital
 * eccentricity, mean anomaly and the result of last iteration or the starting value.
 * Angles are expressed in radians.
 */
func eps(ecc, mna, x float64) float64 {
	t1 := math.Cos(x)
	t2 := -1 + ecc * t1
	t3 := math.Sin(x)
	t4 := ecc * t3
	t5 := -x + t4 + mna
	t6 := t5 / (0.5 * t5 * t4 / t2 + t2)

	return t5 / ((0.5 * t3 - t1 * t6 / 6) * ecc * t6 + t2)
}

/* Calculates true anomaly from orbital eccentricity and eccentric anomaly. */
func TrueAnomaly(ecc, eca float64) float64 {
	ecaRad := Rad(eca)
	return 2 * Deg(math.Atan2(math.Sqrt(1 + ecc) * math.Sin(ecaRad / 2), math.Sqrt(1 - ecc) * math.Cos(ecaRad / 2)))
}

/* Returns orbital radius given semi-major axis, orbital eccentricity and true anomaly. */
func OrbitalRadius(sma, ecc, tra float64) float64 {
	return sma * ((1 - math.Pow(ecc, 2)) / (1 + ecc * math.Cos(Rad(tra))))
}

/* Calculates orbital velocity from semi-major axis, orbital radius and dominant body mass. */
func OrbitalVelocity(r, sma, dominantMass float64) float64 {
	return math.Sqrt(G * dominantMass * (2 / r - 1 / sma))
}

/* Calculates true longitude from true anomaly and longitude of periapsis. */
func TrueLongitude(tra, lpe float64) float64 {
	return math.Mod(tra + lpe, 360)
}

/* Converts epoch year and fraction of a day extracted from TLE to unix time in seconds. */
func EpochToUnix(epochYear int, epochDay float64) int64 {
	return time.Date(epochYear, time.January, 0, 0, 0, 0, 0, time.UTC).Unix() + int64(86400 * epochDay)
}

/* Calculates the average rate of sweep from orbital period [deg/sec]. */
func sweepRate(t float64) float64 {
	return 360 / t
}
