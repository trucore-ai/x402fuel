package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trucore-ai/x402fuel/internal/config"
	"github.com/trucore-ai/x402fuel/internal/keystore"
)

func init() {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show wallet summary and daemon status",
		RunE:  runStatus,
	}
	statusCmd.Flags().StringVarP(&serveConfig, "config", "c", "", "Path to YAML config file")
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(serveConfig)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	fmt.Println("⚡ x402Fuel Status")
	fmt.Println("─────────────────")
	fmt.Printf("Config:     %s\n", serveConfig)
	fmt.Printf("Key dir:    %s\n", cfg.Wallet.KeyDir)
	fmt.Printf("RPC:        %s\n", cfg.Chain.RPCURL)
	fmt.Printf("Chain ID:   %d\n", cfg.Chain.ChainID)
	fmt.Printf("Server:     port %d\n", cfg.Server.Port)
	fmt.Printf("Max/txn:    %s USDC\n", cfg.Proxy.MaxPerTxn)
	fmt.Printf("Daily cap:  %s USDC\n", cfg.Proxy.DailyCap)
	fmt.Printf("Allowlist:  %v\n", cfg.Proxy.Allowlist)
	fmt.Printf("Telemetry:  %v\n", cfg.Telemetry.Enabled)

	ks, _ := keystore.New(cfg.Wallet.KeyDir)
	addrs, err := ks.List()
	if err != nil {
		return fmt.Errorf("list wallets: %w", err)
	}
	fmt.Printf("\nWallets:    %d\n", len(addrs))
	for _, a := range addrs {
		fmt.Printf("  - %s\n", a)
	}

	return nil
}