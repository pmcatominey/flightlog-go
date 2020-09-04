package app

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pmcatominey/flightlog-go/pkg/flights"
	"github.com/pmcatominey/flightlog-go/pkg/world"
)

type FlightDetail struct {
	flight *flights.Flight

	activeTab int
}

func (l *FlightDetail) Render(width, height int) *ui.Grid {
	f := l.flight

	details := widgets.NewParagraph()
	format := `
[Scheduled Departure](mod:bold): %s
[Scheduled Arrival](mod:bold):   %s
[Route](mod:bold):               %s => %s

[Number](mod:bold):              %s
[Operator](mod:bold):            %s
[Aircraft](mod:bold):            %s
[Registration](mod:bold):        %s

[Flight Time](mod:bold):         %s
[Distance (km)](mod:bold):       %d
`
	details.Text = fmt.Sprintf(format,
		f.ScheduledDeparture.String(),
		f.ScheduledArrival.String(),
		f.From, f.To,
		f.Number,
		f.Operator,
		f.Aircraft,
		f.Registration,
		f.Stats.Duration,
		f.Stats.DistanceTraveled)
	details.Title = "Schedule"

	tabs := widgets.NewTabPane("(d) Details", "(g) Graphs", "(m) Map")
	tabs.ActiveTabIndex = l.activeTab
	tabs.Title = "Back: esc"

	renderTab := func() ui.Drawable {
		switch tabs.ActiveTabIndex {
		case 1:
			return graphsTab(f)
		case 2:
			return world.DrawMap(width, height, f)
		}
		return details
	}

	g := ui.NewGrid()
	g.Set(
		ui.NewRow(0.05, tabs),
		ui.NewRow(0.95, renderTab()),
	)
	g.SetRect(0, 0, width, height)

	return g
}

func (l *FlightDetail) ProcessEvent(e ui.Event, a *App) bool {
	switch e.ID {
	case "d":
		l.activeTab = 0
		return true
	case "g":
		l.activeTab = 1
		return true
	case "m":
		l.activeTab = 2
		return true
	}
	return false
}

func graphsTab(flight *flights.Flight) *ui.Grid {
	altData := []float64{}
	for _, t := range flight.Track {
		altData = append(altData, float64(t.Alt))
	}
	altGraph := widgets.NewPlot()
	altGraph.Data = [][]float64{altData}
	altGraph.LineColors[0] = ui.ColorCyan
	altGraph.Title = "Altitude"
	altGraph.TitleStyle.Fg = ui.ColorWhite

	spdData := []float64{}
	for _, t := range flight.Track {
		spdData = append(spdData, float64(t.Spd))
	}
	speedGraph := widgets.NewPlot()
	speedGraph.Data = [][]float64{spdData}
	speedGraph.LineColors[0] = ui.ColorCyan
	speedGraph.Title = "Speed"
	speedGraph.TitleStyle.Fg = ui.ColorWhite
	speedGraph.HorizontalScale = 1

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/2, speedGraph),
		ui.NewRow(1.0/2, altGraph),
	)

	return grid
}
