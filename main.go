package main

import (
	"fmt"
	"os"

	"github.com/pmcatominey/flightlog-go/pkg/app"
	"github.com/pmcatominey/flightlog-go/pkg/flights"
)

func main() {
	l, err := flights.NewLog("./data")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	app.New(l).Run()
}
