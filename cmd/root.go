package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "scram_tcp",
	Short: "SCRAM TCP Server/Client example",
	Long: `SCRAM TCP Server/Client example
		- Implemented using Go language and spf13/cobra tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use commands serve or client")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
