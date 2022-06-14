package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"pcfg_tool/internal/tool"
)

var outsideCmd = &cobra.Command{
	Use:   "outside",
	Short: "Calculate viterbi outside weights",
	Long: "Calculates Viterbi outside weights for each non-terminal of the grammar and prints them on the standard " +
		"output. If the optional argument GRAMMAR is given, then the outside weights are stored in the file " +
		"GRAMMAR.outside.",
	Args: cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		n := cmd.Flag("initial-nonterminal").Value.String()

		var outside string

		if len(args) == 3 {
			outside = args[2] + ".outside"
		}

		if err := tool.Outside(args[0], args[1], n, outside); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	outsideCmd.PersistentFlags().StringP("initial-nonterminal", "i", "ROOT", "")

	rootCmd.AddCommand(outsideCmd)
}
