package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"pcfg_tool/internal/tool"
)

var induceCmd = &cobra.Command{
	Use:   "induce",
	Short: "Induce probabilistic context-free grammar",
	Long: "Reads a sequence of constituent trees from standard input and outputs a PCFG induced from these trees on " +
		"standard output. If the optional GRAMMAR argument is specified, then the PCFG is stored in the " +
		"GRAMMAR.rules, GRAMMAR.lexicon, and GRAMMAR.words files instead.`,",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var rules string
		var lexicon string
		var words string

		if len(args) == 1 {
			rules = args[0] + ".rules"
			lexicon = args[0] + ".lexicon"
			words = args[0] + ".words"
		}

		if err := tool.Induce(os.Getenv("STDIN"), rules, lexicon, words); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(induceCmd)
}
