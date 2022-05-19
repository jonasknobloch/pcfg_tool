package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var debinarizeCmd = &cobra.Command{
	Use:     "debinarize",
	Aliases: []string{"debinarise"},
	Short:   "Debinarize binarized constitutent trees",
	Long: "Reads a sequence of (binarized) constituent trees from standard input and outputs the original " +
		"(non-binarized) constituent trees on standard output.",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(22)
	},
}

func init() {
	rootCmd.AddCommand(debinarizeCmd)
}
