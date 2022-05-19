package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var smoothCmd = &cobra.Command{
	Use:   "smooth",
	Short: "Smooth constituent trees",
	Long: "Reads a sequence of constituent trees from standard input and outputs the trees obtained by smoothing on " +
		"standard output.",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(22)
	},
}

func init() {
	smoothCmd.PersistentFlags().Int64P("threshold", "t", 0, "")

	rootCmd.AddCommand(smoothCmd)
}
