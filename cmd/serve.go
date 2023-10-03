package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"scram_tcp/internal"
)

const (
	defaultUsername = "tcp"
	defaultPassword = "tcp_password"
	defaultPort     = "9001"
)

var (
	username string
	password string
	port     string
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "SCRAM Server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			server, err := internal.NewServer(username, password, port)
			if err != nil {
				fmt.Printf("error on creating server: %s", err.Error())
				return
			}
			server.Start()
		},
	}
)

func init() {
	serveCmd.Flags().StringVar(&username, "username", defaultUsername, "Username (required if password is set)")
	serveCmd.Flags().StringVar(&password, "password", defaultPassword, "Password (required if username is set)")
	serveCmd.Flags().StringVar(&port, "port", defaultPort, "Port (required if username and password is set)")
	serveCmd.MarkFlagsRequiredTogether("username", "password", "port")
	rootCmd.AddCommand(serveCmd)
}
