package main

import (
	"github.com/xyths/game-cashier/cmd/utils"
	"github.com/xyths/game-cashier/puller"
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
	config, err := utils.ParseConfig(ctx.String(utils.ConfigFlag.Name))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("network: %s, api_key: %s, manager: %s",
		config.Server.Network, config.Server.ApiKey, config.Server.Manager)
	p := new(puller.Puller)

	if err = p.Init(config.Server.Network, config.Server.ApiKey, config.Server.Manager); err != nil {
		log.Fatal(err)
	}
	if err = p.Pull(ctx.Context); err != nil {
		log.Fatal(err)
	}
	return nil
}

func notify(ctx *cli.Context) error {
	log.Println("notify started")
	return nil
}
