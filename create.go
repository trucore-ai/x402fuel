package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trucore-ai/x402fuel/internal/config"
	"github.com/trucore-ai/x402fuel/internal/keystore"
)

var createLabel string

func init() {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new agent wallet",
		Long:  "Generates a new ECDSA key pair, encrypts it with your passphrase, and stores it locally. Keys never leave this machine.",
		RunE:  runCreate,
	}
	createCmd.Flags().StringVarP(&createLabel, "label", "l", "", "Human-readable wallet label")
	createCmd.Flags().StringVarP(&serveConfig, "config", "c", "", "Path to YAML config file")
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(serveConfig)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	fmt.Print("Enter passphrase: ")
	var passphrase string
	fmt.Scanln(&passphrase)
	if passphrase == "" {
		return fmt.Errorf("passphrase cannot be empty")
	}

	ks, err := keystore.New(cfg.Wallet.KeyDir)
	if err != nil {
		return fmt.Errorf("init keystore: %w", err)
	}

	addr, err := ks.Generate(passphrase)
	if err != nil {
		return fmt.Errorf("create wallet: %w", err)
	}

	fmt.Printf("\n✅ Wallet created!\n")
	fmt.Printf("   Address: %s\n", addr)
	if createLabel != "" {
		fmt.Printf("   Label:   %s\n", createLabel)
	}
	fmt.Printf("   Keys stored at: %s\n", cfg.Wallet.KeyDir)
	fmt.Printf("\nFund this address with USDC on Base to start making payments.\n")
	return nil
}