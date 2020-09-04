package flights

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/umahmood/haversine"
)

var (
	trackHeaders        = []string{"Timestamp", "UTC", "Callsign", "Position", "Altitude", "Speed", "Direction"}
	minimumTrackRecords = 3
)

type Flight struct {
	FlightInfo

	Stats *FlightStats

	Track []*TrackEntry
}

type FlightInfo struct {
	// IATA Airport Codes
	From string
	To   string
	// Flight Number
	Number string
	// Operating Airline
	Operator string
	// Aircraft Type
	Aircraft string
	// Aircraft Registration
	Registration string
	// Times
	ScheduledDeparture time.Time
	ScheduledArrival   time.Time
	// Times calculated from Track
	ActualDeparture time.Time `json:"-"`
	ActualArrival   time.Time `json:"-"`
}

type TrackEntry struct {
	// UNIX Time
	Timestamp int64
	// Position in Decimal degrees
	Lat, Lng float64
	// Altitude in feet
	Alt int
	// Ground speed in Knots
	Spd int
	// Heading in degrees
	Hdg int
}

func parseFlight(data io.Reader, trackData io.Reader) (*Flight, error) {
	f := &Flight{
		Stats: &FlightStats{},
		Track: []*TrackEntry{},
	}

	if err := json.NewDecoder(data).Decode(&f.FlightInfo); err != nil {
		return nil, fmt.Errorf("unable to parse flight meta data: %v", err)
	}

	records, err := csv.NewReader(trackData).ReadAll()
	if err != nil {
		return nil, err
	}

	// ensure we have at least some data
	if len(records) < minimumTrackRecords {
		return nil, fmt.Errorf("not enough data; header and at least two data points are required")
	}

	// ensure we have the expected number of columns in order
	if !reflect.DeepEqual(records[0], trackHeaders) {
		return nil, fmt.Errorf("headers")
	}

	var (
		lastTrack                  *TrackEntry
		minTimestamp, maxTimestamp int64
	)

	// process track data, skipping header row
	for i, r := range records[1:] {
		t, err := parseTrackEntry(r)
		if err != nil {
			return nil, fmt.Errorf("failed to parse track data at %d: %v", i, err)
		}

		// calculate stats
		f.Stats.MaxAltitude = maxInt(f.Stats.MaxAltitude, t.Alt)
		f.Stats.MaxSpeed = maxInt(f.Stats.MaxSpeed, t.Spd)

		if lastTrack != nil && t.Alt > 0 {
			from := haversine.Coord{Lat: lastTrack.Lat, Lon: lastTrack.Lng}
			to := haversine.Coord{Lat: t.Lat, Lon: t.Lng}
			_, km := haversine.Distance(from, to)
			f.Stats.DistanceTraveled += int(km)
		}

		// capture start and end times for journey
		if t.Alt > 0 {
			if minTimestamp == 0 {
				minTimestamp = t.Timestamp
			}

			minTimestamp = minInt64(t.Timestamp, minTimestamp)
			maxTimestamp = maxInt64(t.Timestamp, maxTimestamp)
		}

		lastTrack = t
		f.Track = append(f.Track, t)
	}

	f.Stats.Duration = time.Second * time.Duration(maxTimestamp-minTimestamp)

	return f, nil
}

func parseTrackEntry(record []string) (*TrackEntry, error) {
	if len(record) < len(trackHeaders) {
		return nil, fmt.Errorf("row length lol")
	}

	ts, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing Timestamp: %v", err)
	}

	latlon := strings.Split(record[3], ",")
	if len(latlon) != 2 {
		return nil, fmt.Errorf("expected Position column to contain 2 values got: %v", len(latlon))
	}

	lat, err := strconv.ParseFloat(latlon[0], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing Latitude: %v", err)
	}

	lng, err := strconv.ParseFloat(latlon[1], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing Longitude: %v", err)
	}

	alt, err := strconv.Atoi(record[4])
	if err != nil {
		return nil, fmt.Errorf("error parsing Altitude: %v", err)
	}

	spd, err := strconv.Atoi(record[5])
	if err != nil {
		return nil, fmt.Errorf("error parsing Speed: %v", err)
	}

	hdg, err := strconv.Atoi(record[6])
	if err != nil {
		return nil, fmt.Errorf("error parsing Heading: %v", err)
	}

	t := &TrackEntry{
		Timestamp: ts,
		Lat:       lat,
		Lng:       lng,
		Alt:       alt,
		Spd:       spd,
		Hdg:       hdg,
	}
	return t, nil
}
