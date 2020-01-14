package utils

import "gopkg.in/urfave/cli.v2"

var (
	ConfigFlag = &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "config.json",
		Usage:   "Load configuration from `file`",
	}
	DryRunFlag = &cli.BoolFlag{
		Name:  "dry-run",
		Value: false,
		Usage: "dry run.",
	}
	AfterFlag = &cli.StringFlag{
		Name:  "after",
		Usage: "the history after this time, eg. 2006-01-02T15:04:05.000-07:00",
	}
	BeforeFlag = &cli.StringFlag{
		Name:  "before",
		Usage: "the history before this time, eg. 2006-01-02T15:04:05.000-07:00",
	}
)
