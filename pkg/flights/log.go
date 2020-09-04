package flights

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	searchExt = ".json"
	trackExt  = ".csv"
)

type Log struct {
	Flights []*Flight

	StatsTotals *FlightStats
}

func NewLog(dir string) (*Log, error) {
	l := &Log{
		Flights:     []*Flight{},
		StatsTotals: &FlightStats{},
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing %q: %v", path, err)
		}

		// ignore dirs / non-csv files
		if info.IsDir() || filepath.Ext(path) != searchExt {
			return nil
		}

		data, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening %q: %v", path, err)
		}

		prefix := strings.TrimSuffix(path, searchExt)
		track, err := os.Open(prefix + trackExt)
		if err != nil {
			return fmt.Errorf("error opening track %q: %v", prefix+trackExt, err)
		}

		flight, err := parseFlight(data, track)
		if err != nil {
			return fmt.Errorf("error parsing flight data %q: %v", path, err)
		}

		l.StatsTotals = sumFlightStats(l.StatsTotals, flight.Stats)

		l.Flights = append(l.Flights, flight)
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.SliceStable(l.Flights, func(i, j int) bool {
		return l.Flights[i].ScheduledDeparture.Before(l.Flights[j].ScheduledDeparture)
	})

	return l, nil
}
