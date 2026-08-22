package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/trucore-ai/x402fuel/internal/api"
	"github.com/trucore-ai/x402fuel/internal/config"
	"github.com/trucore-ai/x402fuel/internal/dashboard"
	"github.com/trucore-ai/x402fuel/internal/events"
	"github.com/trucore-ai/x402fuel/internal/keystore"
	"github.com/trucore-ai/x402fuel/internal/policy"
	"github.com/trucore-ai/x402fuel/internal/types"
)

var (
	servePort   int
	serveConfig string
)

func init() {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the x402Fuel daemon (API + proxy + dashboard)",
		RunE:  runServe,
	}
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8420, "API server listen port")
	serveCmd.Flags().StringVarP(&serveConfig, "config", "c", "", "Path to YAML config file")
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(serveConfig)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if servePort != 8420 {
		cfg.Server.Port = servePort
	}

	ks, err := keystore.New(cfg.Wallet.KeyDir)
	if err != nil {
		return fmt.Errorf("init keystore: %w", err)
	}

	eventLog, err := events.NewLogger(cfg.Events.LogPath)
	if err != nil {
		return fmt.Errorf("init event log: %w", err)
	}
	defer eventLog.Close()

	sp := &types.SpendPolicy{
		MaxPerTxn: cfg.Proxy.MaxPerTxn,
		DailyCap:  cfg.Proxy.DailyCap,
		Allowlist: cfg.Proxy.Allowlist,
	}
	polEng := policy.NewEngine(sp)

	apiSrv := api.NewServer(ks, polEng)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	apiSrv.Register(r)

	dashMux := http.NewServeMux()
	dashboard.Register(dashMux, ks, polEng)
	r.Mount("/", dashMux)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("⚡ x402Fuel daemon starting on %s", addr)
	log.Printf("   Dashboard: http://localhost%s", addr)
	log.Printf("   API:       http://localhost%s/api", addr)
	log.Printf("   Keys:      %s", cfg.Wallet.KeyDir)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("shutting down...")
		eventLog.Close()
		os.Exit(0)
	}()

	return http.ListenAndServe(addr, r)
}