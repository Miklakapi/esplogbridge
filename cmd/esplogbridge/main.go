package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Miklakapi/esplogbridge/internal/cli"
	"github.com/Miklakapi/esplogbridge/internal/config"
	"github.com/Miklakapi/esplogbridge/internal/version"
)

func main() {
	_, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	flags, err := cli.ParseFlags()
	if err != nil {
		printErrorAndExit(err, 2)
	}

	if flags.Version {
		fmt.Println(version.VersionString())
		return
	}

	_, err = config.Load(flags.ConfigPath)
	if err != nil {
		printErrorAndExit(err, 2)
	}
}

func printErrorAndExit(err error, code int) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(code)
}
