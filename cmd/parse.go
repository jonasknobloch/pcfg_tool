package cmd

import (
	"github.com/spf13/cobra"
	"pcfg_tool/internal/tool"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse natural language sentences",
	Long: "Reads a sequence of natural language sentences from standard input and outputs the associated best parse " +
		"trees in PTB format or (NOPARSE <sentence>) on standard output. RULES and LEXICON are the file names of the " +
		"PCFG.",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		stdin := OpenStdin()
		defer stdin.Close()

		n := cmd.Flag("initial-nonterminal").Value.String()

		tool.Parse(args[0], args[1], n, stdin)
	},
}

func init() {
	parseCmd.PersistentFlags().StringP("paradigma", "p", "", "")
	parseCmd.PersistentFlags().StringP("initial-nonterminal", "i", "ROOT", "")

	parseCmd.PersistentFlags().BoolP("unking", "u", false, "")
	parseCmd.PersistentFlags().BoolP("smoothing", "s", false, "")

	parseCmd.PersistentFlags().Int64P("threshold-beam", "t", 0, "")
	parseCmd.PersistentFlags().Int64P("rank-beam", "r", 0, "")
	parseCmd.PersistentFlags().Int64P("kbest", "k", 0, "")

	parseCmd.PersistentFlags().StringP("astar", "a", "", "")

	rootCmd.AddCommand(parseCmd)
}
