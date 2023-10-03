package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"scram_tcp/internal"
)

var (
	clientUsername string
	clientPassword string
	clientAddress  string
	clientCmd      = &cobra.Command{
		Use:   "client",
		Short: "SCRAM Client",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			clt, err := internal.NewClient(clientAddress, clientUsername, clientPassword)
			if err != nil {
				fmt.Printf("error on creating client: %s\n", err.Error())
				return
			}
			clt.MakeRequest()
		},
	}
)

func init() {
	clientCmd.Flags().StringVar(&clientUsername, "username", "", "Username (required if password is set)")
	clientCmd.Flags().StringVar(&clientPassword, "password", "", "Password (required if username is set)")
	clientCmd.Flags().StringVar(&clientAddress, "address", "", "Address (required if username and password is set)")
	clientCmd.MarkFlagsRequiredTogether("username", "password", "address")
	rootCmd.AddCommand(clientCmd)
}
