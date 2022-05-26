package cmd

import (
	"github.com/spf13/cobra"
	"log"
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
		stdin := OpenStdin()
		defer stdin.Close()

		g := tool.Induce(stdin)

		var grammar string

		if len(args) > 0 {
			grammar = args[0]
		}

		if err := g.Export(grammar); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(induceCmd)
}
