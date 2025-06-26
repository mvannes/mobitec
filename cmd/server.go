package cmd

import (
	"log"
	"mobitec/internal/server"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start an http server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal(server.Start())
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
