package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"mobitec/internal/server"
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
