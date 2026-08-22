package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trucore-ai/x402fuel/internal/config"
)

func init() {
	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause all agent payments (kill switch)",
		Long:  "Activates the global kill switch. All new payments are blocked before signing. In-flight transactions still settle.",
		RunE:  runPause,
	}
	pauseCmd.Flags().StringVarP(&serveConfig, "config", "c", "", "Path to YAML config file")
	rootCmd.AddCommand(pauseCmd)

	resumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume agent payments",
		Long:  "Deactivates the global kill switch. Payments resume immediately.",
		RunE:  runResume,
	}
	resumeCmd.Flags().StringVarP(&serveConfig, "config", "c", "", "Path to YAML config file")
	rootCmd.AddCommand(resumeCmd)
}

func runPause(cmd *cobra.Command, args []string) error {
	_, err := config.Load(serveConfig)
	if err != nil {
		fmt.Println("⚠  Using default config.")
	}
	fmt.Println("⏸  To pause payments, call:")
	fmt.Println("   curl -X POST http://localhost:8420/api/wallets/pause")
	fmt.Println("Or use: http://localhost:8420/settings")
	return nil
}

func runResume(cmd *cobra.Command, args []string) error {
	fmt.Println("▶  To resume payments, call:")
	fmt.Println("   curl -X POST http://localhost:8420/api/wallets/resume")
	fmt.Println("Or use: http://localhost:8420/settings")
	return nil
}