package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pmcatominey/flightlog-go/pkg/app"
	"github.com/pmcatominey/flightlog-go/pkg/flights"
)

var data = flag.String("d", "./data", "data directory")

func main() {
	flag.Parse()

	l, err := flights.NewLog(*data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	app.New(l).Run()
}
