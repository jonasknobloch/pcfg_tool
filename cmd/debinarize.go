package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var debinarizeCmd = &cobra.Command{
	Use:     "debinarize",
	Aliases: []string{"debinarise"},
	Short:   "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(22)
	},
}

func init() {
	rootCmd.AddCommand(debinarizeCmd)
}
