package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var outsideCmd = &cobra.Command{
	Use:   "outside",
	Short: "Calculate viterbi outside weights",
	Long: "Calculates Viterbi outside weights for each non-terminal of the grammar and prints them on the standard " +
		"output. If the optional argument GRAMMAR is given, then the outside weights are stored in the file " +
		"GRAMMAR.outside.",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(22)
	},
}

func init() {
	outsideCmd.PersistentFlags().StringP("initial-nonterminal", "i", "ROOT", "")

	rootCmd.AddCommand(outsideCmd)
}
