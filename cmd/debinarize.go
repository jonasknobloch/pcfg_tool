package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"pcfg_tool/internal/tool"
)

var debinarizeCmd = &cobra.Command{
	Use:     "debinarize",
	Aliases: []string{"debinarise"},
	Short:   "Debinarize binarized constitutent trees",
	Long: "Reads a sequence of (binarized) constituent trees from standard input and outputs the original " +
		"(non-binarized) constituent trees on standard output.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tool.Transform(os.Getenv("STDIN"), os.Getenv("STDOUT"), tool.Demarkovize()); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(debinarizeCmd)
}
