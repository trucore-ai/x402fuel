package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

var (
	servePort   int
	serveConfig string
)

func init() {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the x402fuel daemon with health endpoint",
		Long:  `Starts the x402fuel wallet daemon, serving a health-check endpoint and the proxy gateway.`,
		RunE:  runServe,
	}

	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8420, "API server listen port")
	serveCmd.Flags().StringVarP(&serveConfig, "config", "c", "", "Path to YAML config file (optional)")

	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	addr := fmt.Sprintf(":%d", servePort)
	fmt.Printf("x402fuel daemon listening on %s\n", addr)
	return http.ListenAndServe(addr, r)
}