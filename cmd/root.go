package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "x402fuel",
	Short: "Non-custodial HTTP 402 wallet daemon for AI agents",
	Long: `x402Fuel gives AI agents a USDC wallet on Base that speaks HTTP 402.
Agents can autonomously pay for APIs, data, and compute against x402 paywalls.
Keys never leave your machine.`,
}

func Execute() error {
	return rootCmd.Execute()
}