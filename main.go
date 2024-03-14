package main

import (
	"errors"
	"os"

	"github.com/kolah/github-batch-updater/cli"
	"github.com/kolah/github-batch-updater/di"
)

func recoverPanic(app *di.Container) {
	if r := recover(); r != nil {
		var err error
		switch x := r.(type) {
		case string:
			err = errors.New(x)
		case error:
			err = x
		default:
			err = errors.New("unknown panic")
		}
		app.Logger().Fatal("panic", "error", err)

		os.Exit(1)
	}
}

func run() int {
	app := di.NewContainer()

	defer recoverPanic(app)

	if err := cli.RootCmd(app).Execute(); err != nil {
		app.Logger().Error("CLI command returned an error", err)

		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
