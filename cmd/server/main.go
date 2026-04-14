package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/JetManiack/alertmanager-template-preview/internal/api"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "alertmanager-template-preview",
		Usage: "A web application for previewing Prometheus Alertmanager templates",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "listen-address",
				Aliases: []string{"l"},
				Value:   ":8080",
				Usage:   "Address to listen on for HTTP requests",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			addr := cmd.String("listen-address")
			fmt.Printf("Starting server on %s...\n", addr)

			router := api.SetupRouter()
			return router.Run(addr)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
