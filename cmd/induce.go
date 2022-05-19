package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"pcfg_tool/internal/tool"
)

var induceCmd = &cobra.Command{
	Use:   "induce",
	Short: "A brief description of your command",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stdin := OpenStdin()
		defer stdin.Close()

		g := tool.Induce(stdin)

		if len(args) == 0 {
			g.Print()
			return
		}

		if err := g.Export(args[0]); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(induceCmd)
}
