package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var unkCmd = &cobra.Command{
	Use:   "unk",
	Short: "Unk constituent trees",
	Long: "Reads a sequence of constituent trees from standard input and outputs the trees obtained by trivial " +
		"unking on standard output.",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(22)
	},
}

func init() {
	unkCmd.PersistentFlags().Int64P("threshold", "t", 0, "")

	rootCmd.AddCommand(unkCmd)
}
