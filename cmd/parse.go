package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
	"os"
	"pcfg_tool/internal/tool"
)

const ParadigmCYK = "cyk"
const ParadigmDeductive = "deductive"

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse natural language sentences",
	Long: "Reads a sequence of natural language sentences from standard input and outputs the associated best parse " +
		"trees in PTB format or (NOPARSE <sentence>) on standard output. RULES and LEXICON are the file names of the " +
		"PCFG.",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		p := cmd.Flag("paradigm").Value.String()
		n := cmd.Flag("initial-nonterminal").Value.String()
		a := cmd.Flag("astar").Value.String()

		var u bool
		var s bool
		var t float64
		var r int
		var k int

		var err error

		if u, err = cmd.Flags().GetBool("unking"); err != nil {
			log.Fatal(err)
		}

		if s, err = cmd.Flags().GetBool("smoothing"); err != nil {
			log.Fatal(err)
		}

		if t, err = cmd.Flags().GetFloat64("threshold-beam"); err != nil {
			log.Fatal(err)
		}

		if r, err = cmd.Flags().GetInt("rank-beam"); err != nil {
			log.Fatal(err)
		}

		if k, err = cmd.Flags().GetInt("kbest"); err != nil {
			log.Fatal(err)
		}

		if p != ParadigmDeductive {
			if p == ParadigmCYK {
				os.Exit(22)
			}

			log.Fatal(errors.New("unknown parser paradigm"))
		}

		if s || k != 0 {
			os.Exit(22)
		}

		if err := tool.Parse(args[0], args[1], n, u, t, r, a, os.Getenv("STDIN")); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	parseCmd.PersistentFlags().StringP("paradigm", "p", "deductive", "")
	parseCmd.PersistentFlags().StringP("initial-nonterminal", "i", "ROOT", "")

	parseCmd.PersistentFlags().BoolP("unking", "u", false, "")
	parseCmd.PersistentFlags().BoolP("smoothing", "s", false, "")

	parseCmd.PersistentFlags().Float64P("threshold-beam", "t", 0, "")
	parseCmd.PersistentFlags().IntP("rank-beam", "r", 0, "")
	parseCmd.PersistentFlags().IntP("kbest", "k", 0, "")

	parseCmd.PersistentFlags().StringP("astar", "a", "", "")

	rootCmd.AddCommand(parseCmd)
}
