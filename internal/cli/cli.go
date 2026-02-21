package cli

import (
	"errors"
	"flag"
	"strings"
)

type CLI struct {
	ConfigPath string
	Version    bool
}

func ParseFlags() (CLI, error) {
	var cfg CLI

	flag.StringVar(&cfg.ConfigPath, "config", "", "Path to YAML config file (e.g. /etc/gometrum.yaml)")
	flag.StringVar(&cfg.ConfigPath, "c", "", "Shorthand for --config")

	flag.BoolVar(&cfg.Version, "version", false, "Show version and exit")
	flag.BoolVar(&cfg.Version, "v", false, "Shorthand for --version")

	flag.Parse()

	if strings.TrimSpace(cfg.ConfigPath) == "" {
		cfg.ConfigPath = "./gometrum.yaml"
	}

	if err := validateFlags(cfg); err != nil {
		return CLI{}, err
	}

	return cfg, nil
}

func validateFlags(c CLI) error {
	exitModes := 0
	if c.Version {
		exitModes++
	}

	if exitModes > 1 {
		return errors.New("choose only one of: --version,")
	}

	return nil
}
