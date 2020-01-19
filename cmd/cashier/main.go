package main

import (
	"context"
	"fmt"
	"github.com/xyths/game-cashier/cmd/utils"
	"github.com/xyths/game-cashier/node"
	"gopkg.in/urfave/cli.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var app *cli.App

func init() {
	app = &cli.App{
		Name:    filepath.Base(os.Args[0]),
		Usage:   "game cashier on blockchain",
		Version: "0.1.2",
		Action:  serve,
	}

	app.Commands = []*cli.Command{
		serveCommand,
		pullCommand,
		notifyCommand,
		downloadCommand,
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
	n := node.Node{}
	n.Init(ctx.Context, ctx.String(utils.ConfigFlag.Name))

	srv := &http.Server{
		Addr:    n.Addr,
		Handler: n.Engine,
	}

	go func() {
		log.Println("Listen on", srv.Addr)
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx2); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
	return nil
}
