package flights

import "time"

type FlightStats struct {
	// Distance Traveled including on ground in KM
	DistanceTraveled int
	// Duration spent flying
	Duration time.Duration
	// Highest Altitude Reached in Feet
	MaxAltitude int
	// Highest Speed reached in Knots
	MaxSpeed int
}

func sumFlightStats(a, b *FlightStats) *FlightStats {
	return &FlightStats{
		DistanceTraveled: a.DistanceTraveled + b.DistanceTraveled,
		Duration:         a.Duration + b.Duration,
		MaxAltitude:      maxInt(a.MaxAltitude, b.MaxAltitude),
		MaxSpeed:         maxInt(a.MaxSpeed, b.MaxSpeed),
	}
}
