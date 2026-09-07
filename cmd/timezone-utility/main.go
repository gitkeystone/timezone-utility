package main

import (
	"os"

	_ "time/tzdata"

	"timezone-utility/internal/cli"
	"timezone-utility/internal/clock"
	"timezone-utility/internal/data"
	"timezone-utility/internal/output"
)

func main() {
	store, err := data.Load()
	if err != nil {
		output.WriteError(os.Stderr, err)
		os.Exit(2)
	}
	app := cli.App{
		Store: store,
		Clock: clock.RealClock{},
		Out:   os.Stdout,
		Err:   os.Stderr,
	}
	os.Exit(cli.Run(os.Args[1:], app))
}
