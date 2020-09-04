package world

import (
	"image"
	"math"

	ui "github.com/gizak/termui/v3"
	"github.com/pmcatominey/flightlog-go/pkg/countries"
	"github.com/pmcatominey/flightlog-go/pkg/flights"
)

var (
	zoom    float64 = 3
	xoffset float64 = 0
	yoffset float64 = 0

	deltaOffset float64 = 25
	deltaZoom   float64 = 0.25

	flightIndex = 0
)

func DrawMap(width, height int, f *flights.Flight) *ui.Canvas {
	width = int(float64(width) * 1.5)
	canvas := ui.NewCanvas()

	for _, c := range countries.Countries {
		for i, p := range c {
			// the last one fucks it all
			if i == len(c)-1 {
				continue
			}

			x, y := mercator(width, height, p.Lat, p.Lng)
			p2 := c[i+1]
			x2, y2 := mercator(width, height, p2.Lat, p2.Lng)
			if x < 0 || y < 0 || x2 < 0 || y2 < 0 {
				continue
			}

			canvas.SetLine(image.Point{int(x), int(y)}, image.Point{int(x2), int(y2)}, ui.ColorGreen)
		}
	}

	for k, t := range f.Track {
		x, y := mercator(width, height, t.Lat, t.Lng)
		switch {
		case k < len(f.Track)-1:
			t2 := f.Track[k+1]
			x2, y2 := mercator(width, height, t2.Lat, t2.Lng)
			canvas.SetLine(image.Point{int(x), int(y)}, image.Point{int(x2), int(y2)}, ui.ColorRed)
		}
	}

	canvas.SetRect(0, 0, width, height)

	return canvas
}

func mercators(mapWidth, mapHeight int, lat, lng float64) (x, y float64) {
	TILE_SIZE := float64(mapHeight)
	siny := math.Sin((lat * math.Pi) / 180)

	// Truncating to 0.9999 effectively limits latitude to 89.189. This is
	// about a third of a tile past the edge of the world tile.
	siny = math.Min(math.Max(siny, -0.9999), 0.9999)

	x = TILE_SIZE * (0.5 + lng/360)
	y = TILE_SIZE * (0.5 - math.Log((1+siny)/(1-siny))/(4*math.Pi))

	return
}

func mercator(mapWidth, mapHeight int, lat, lng float64) (x, y float64) {
	x, y = mercators(mapWidth, mapHeight, lat, lng)

	//zoom := float64(zoom)
	scale := math.Pow(2, zoom)

	x = math.Floor(x * scale)
	y = math.Floor(y * scale)

	x += xoffset
	y += yoffset
	return
}
