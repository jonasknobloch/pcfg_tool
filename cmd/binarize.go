package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var binarizeCmd = &cobra.Command{
	Use:     "binarize",
	Aliases: []string{"binarise"},
	Short:   "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(22)
	},
}

func init() {
	rootCmd.AddCommand(binarizeCmd)
}
