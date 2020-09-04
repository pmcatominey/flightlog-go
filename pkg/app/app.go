package app

import (
	ui "github.com/gizak/termui/v3"
	"github.com/pmcatominey/flightlog-go/pkg/flights"
)

type App struct {
	width, height int
	pages         []Page

	flightLog *flights.Log
}

type Page interface {
	Render(width, height int) *ui.Grid
	ProcessEvent(e ui.Event, a *App) bool
}

func New(l *flights.Log) *App {
	return &App{
		pages:     []Page{},
		flightLog: l,
	}
}

func (a *App) Run() error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	a.width, a.height = ui.TerminalDimensions()

	list := &FlightList{
		flightLog: a.flightLog,
	}
	a.PushPage(list)

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			case "<Escape>":
				a.PopPage()
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				a.width = payload.Width
				a.height = payload.Height
				a.draw()
			default:
				if a.pages[len(a.pages)-1].ProcessEvent(e, a) {
					a.draw()
				}
			}
		}
	}
}

func (a *App) draw() {
	ui.Clear()
	widget := a.pages[len(a.pages)-1]
	ui.Render(widget.Render(a.width, a.height))
}

func (a *App) PushPage(p Page) {
	a.pages = append(a.pages, p)
	a.draw()
}

func (a *App) PopPage() {
	if len(a.pages) > 1 {
		a.pages = a.pages[:len(a.pages)-1]
		a.draw()
	}
}
