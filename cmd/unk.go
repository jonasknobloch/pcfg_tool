package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"pcfg_tool/internal/tool"
)

var unkCmd = &cobra.Command{
	Use:   "unk",
	Short: "Unk constituent trees",
	Long: "Reads a sequence of constituent trees from standard input and outputs the trees obtained by trivial " +
		"unking on standard output.",
	Run: func(cmd *cobra.Command, args []string) {
		threshold, err := cmd.Flags().GetInt("threshold")

		if err != nil {
			log.Fatal(err)
		}

		if err := tool.Unk(os.Getenv("STDIN"), os.Getenv("STDOUT"), threshold); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	unkCmd.PersistentFlags().IntP("threshold", "t", 0, "")

	rootCmd.AddCommand(unkCmd)
}
