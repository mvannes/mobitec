package cmd

import (
	"log"
	"mobitec/internal/listen"

	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen to the mobitec-api",
	Run: func(cmd *cobra.Command, args []string) {
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			log.Fatal(err)
		}

		noSerialPort, err := cmd.Flags().GetBool("no-serial-port")
		if err != nil {
			log.Fatal(err)
		}

		listen.Start(host, noSerialPort)
	},
}

func init() {
	listenCmd.Flags().String("host", "localhost:3000", "Host of the api")
	listenCmd.Flags().Bool("no-serial-port", false, "Do not attempt to connect with a serial port when starting the program")
	rootCmd.AddCommand(listenCmd)
}
