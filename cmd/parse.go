package cmd

import (
	"github.com/spf13/cobra"
	"pcfg_tool/internal/pcfg"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		stdin := OpenStdin()
		defer stdin.Close()

		n := cmd.Flag("initial-nonterminal").Value.String()

		pcfg.Parse(args[0], args[1], n, stdin)
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
