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
		Usage: "A web application for previewing Prometheus and Alertmanager templates",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "listen-address",
				Aliases: []string{"l"},
				Value:   ":8080",
				Usage:   "Address to listen on for HTTP requests",
			},
			&cli.StringFlag{
				Name:    "prometheus-url",
				Aliases: []string{"p"},
				Usage:   "Prometheus server URL for 'query' functions in templates",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			addr := cmd.String("listen-address")
			prometheusURL := cmd.String("prometheus-url")
			fmt.Printf("Starting server on %s...\n", addr)
			if prometheusURL != "" {
				fmt.Printf("Using Prometheus server at %s for queries\n", prometheusURL)
			}

			router := api.SetupRouter(prometheusURL)
			return router.Run(addr)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
