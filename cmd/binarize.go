package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var binarizeCmd = &cobra.Command{
	Use:     "binarize",
	Aliases: []string{"binarise"},
	Short:   "Binarize constituent trees",
	Long: "Reads a sequence of constituent trees from standard input and outputs the corresponding binarized " +
		"constituent trees on standard output.",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(22)
	},
}

func init() {
	rootCmd.AddCommand(binarizeCmd)
}
