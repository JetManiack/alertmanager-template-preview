package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JetManiack/alertmanager-template-preview/internal/api"
	"github.com/urfave/cli/v3"
)

var Version = "dev"

func main() {
	cmd := &cli.Command{
		Name:    "alertmanager-template-preview",
		Usage:   "A web application for previewing Prometheus and Alertmanager templates",
		Version: Version,
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
			fmt.Printf("Starting server on %s (version: %s)...\n", addr, Version)
			if prometheusURL != "" {
				fmt.Printf("Using Prometheus server at %s for queries\n", prometheusURL)
			}

			router := api.SetupRouter(prometheusURL)

			srv := &http.Server{
				Addr:    addr,
				Handler: router,
			}

			// Graceful shutdown logic
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Printf("listen error: %v", err)
				}
			}()

			// Wait for interrupt signal to gracefully shut down the server with
			// a timeout of 5 seconds.
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			log.Println("Shutting down server...")

			// The context is used to inform the server it has 5 seconds to finish
			// the request it is currently handling
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("server forced to shutdown: %w", err)
			}

			log.Println("Server exiting")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
