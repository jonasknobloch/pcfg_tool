package cmd

import (
	"log"
	"os"
	"pcfg_tool/internal/tool"

	"github.com/spf13/cobra"
)

var binarizeCmd = &cobra.Command{
	Use:     "binarize",
	Aliases: []string{"binarise"},
	Short:   "Binarize constituent trees",
	Long: "Reads a sequence of constituent trees from standard input and outputs the corresponding binarized " +
		"constituent trees on standard output.",
	Run: func(cmd *cobra.Command, args []string) {
		var horizontal int
		var vertical int

		var err error

		if horizontal, err = cmd.Flags().GetInt("horizontal"); err != nil {
			log.Fatal(err)
		}

		if vertical, err = cmd.Flags().GetInt("vertical"); err != nil {
			log.Fatal(err)
		}

		if err := tool.Binarize(os.Getenv("STDIN"), os.Getenv("STDOUT"), horizontal, vertical); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	binarizeCmd.PersistentFlags().IntP("horizontal", "h", 999, "")
	binarizeCmd.PersistentFlags().IntP("vertical", "v", 1, "")

	rootCmd.AddCommand(binarizeCmd)
}
