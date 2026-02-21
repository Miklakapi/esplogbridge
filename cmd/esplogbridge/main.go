package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Miklakapi/esplogbridge/internal/bridge"
	"github.com/Miklakapi/esplogbridge/internal/cli"
	"github.com/Miklakapi/esplogbridge/internal/config"
	"github.com/Miklakapi/esplogbridge/internal/version"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	flags, err := cli.ParseFlags()
	if err != nil {
		printErrorAndExit(err, 2)
	}

	if flags.Version {
		fmt.Println(version.VersionString())
		return
	}

	cfg, err := config.Load(flags.ConfigPath)
	if err != nil {
		printErrorAndExit(err, 2)
	}

	app, err := bridge.New(cfg)
	if err != nil {
		printErrorAndExit(err, 1)
	}

	if err := app.Run(ctx); err != nil {
		printErrorAndExit(err, 1)
	}
}

func printErrorAndExit(err error, code int) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(code)
}
