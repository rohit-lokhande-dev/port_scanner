package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "metronet",
	Short: "MetroNet - Network Port Scanner",
	Long:  `MetroNet is a high-performance network port scanner with service detection and banner grabbing capabilities.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Add any global flags here if needed
}
