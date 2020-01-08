package main

import (
	"gopkg.in/urfave/cli.v2"
	"log"
)

var (
	serveCommand = &cli.Command{
		Action:  serve,
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "serve the API for recharge history.",
		Flags: []cli.Flag{
		},
	}
	pullCommand = &cli.Command{
		Action:  pull,
		Name:    "pull",
		Aliases: []string{"p"},
		Usage:   "pull the recharge history to database.",
		Flags: []cli.Flag{
		},
	}
	notifyCommand = &cli.Command{
		Action:  notify,
		Name:    "notify",
		Aliases: []string{"n"},
		Usage:   "pull the recharge history to database.",
		Flags: []cli.Flag{
		},
	}
)

func pull(ctx *cli.Context) error {
	log.Println("pull started")
	return nil
}

func notify(ctx *cli.Context) error {
	log.Println("notify started")
	return nil
}