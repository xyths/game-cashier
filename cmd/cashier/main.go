package main

import (
	"fmt"
	"github.com/xyths/game-cashier/cmd/utils"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"path/filepath"
)

var app *cli.App

func init() {
	app = &cli.App{
		Name:    filepath.Base(os.Args[0]),
		Usage:   "game cashier on blockchain",
		Version: "0.0.1",
		Action:  serve,
	}

	app.Commands = []*cli.Command{
		serveCommand,
		pullCommand,
		notifyCommand,
	}
	app.Flags = []cli.Flag{
		utils.ConfigFlag,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func serve(ctx *cli.Context) error {
	log.Println("serve started")
	return nil
}
