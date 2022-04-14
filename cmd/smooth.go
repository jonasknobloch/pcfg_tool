package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var smoothCmd = &cobra.Command{
	Use:   "smooth",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(22)
	},
}

func init() {
	smoothCmd.PersistentFlags().Int64P("threshold", "t", 0, "")

	rootCmd.AddCommand(smoothCmd)
}
