package app

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pmcatominey/flightlog-go/pkg/flights"
)

type FlightList struct {
	flightLog *flights.Log

	selectedRow int
}

func (l *FlightList) Render(width, height int) *ui.Grid {
	p := widgets.NewParagraph()
	format := "[Flights](mod:bold): %d    [Time Flying](mod:bold): %s    [Distance Flown](mod:bold): %dkm"
	p.Text = fmt.Sprintf(format, len(l.flightLog.Flights), l.flightLog.StatsTotals.Duration, l.flightLog.StatsTotals.DistanceTraveled)
	p.Title = "Total Stats"

	t := widgets.NewTable()
	t.Rows = [][]string{
		{"DATE", "DEPARTURE", "NUMBER", "ROUTE", "OPERATOR", "DISTANCE (KM)", "DURATION"},
	}
	for _, f := range l.flightLog.Flights {
		t.Rows = append(t.Rows, []string{f.ScheduledDeparture.Format("02 Jan"), f.ScheduledDeparture.Format("15:04"), f.Number, f.From + " --> " + f.To, f.Operator, fmt.Sprintf("%d", f.Stats.DistanceTraveled), f.Stats.Duration.String()})
	}
	t.RowStyles[0] = ui.Style{Fg: ui.ColorWhite, Modifier: ui.ModifierBold}
	t.RowStyles[1] = ui.Style{Fg: ui.ColorWhite}
	t.RowStyles[l.selectedRow+1] = ui.Style{Fg: ui.ColorYellow}
	t.Title = "Flights [ Scroll: up/down | Select: enter ]"
	t.PaddingTop = 1

	g := ui.NewGrid()
	g.Set(
		ui.NewRow(0.05, p),
		ui.NewRow(0.95, t),
	)
	g.SetRect(0, 0, width, height)

	return g
}

func (l *FlightList) ProcessEvent(e ui.Event, a *App) bool {
	switch e.ID {
	case "<Down>":
		if l.selectedRow == len(l.flightLog.Flights)-1 {
			return false
		}
		l.selectedRow++
		return true
	case "<Up>":
		if l.selectedRow < 1 {
			return false
		}
		l.selectedRow--
		return true
	case "<Enter>":
		a.PushPage(&FlightDetail{flight: l.flightLog.Flights[l.selectedRow]})
		return false
	}
	return false
}
