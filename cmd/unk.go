package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var unkCmd = &cobra.Command{
	Use:   "unk",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(22)
	},
}

func init() {
	unkCmd.PersistentFlags().Int64P("threshold", "t", 0, "")

	rootCmd.AddCommand(unkCmd)
}
