package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var outsideCmd = &cobra.Command{
	Use:   "outside",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(22)
	},
}

func init() {
	outsideCmd.PersistentFlags().StringP("initial-nonterminal", "i", "ROOT", "")

	rootCmd.AddCommand(outsideCmd)
}
