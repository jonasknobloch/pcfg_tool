package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pcfg_tool",
	Short: "Tools for PCFG-based parsing of natural language sentences",
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "", false, "")
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
