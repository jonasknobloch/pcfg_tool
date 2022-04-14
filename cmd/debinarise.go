package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var debinariseCmd = &cobra.Command{
	Use:   "debinarise",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(22)
	},
}

func init() {
	rootCmd.AddCommand(debinariseCmd)
}
