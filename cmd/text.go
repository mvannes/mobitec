package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"io"
	"log"
	"mobitec/internal/flipdot"
	"mobitec/internal/serialport"
)

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "place text on the mobitec destination board",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal(errors.New("a single argument containing the text to display must be provided"))
		}
		text := args[0]
		horizontalOffset, err := cmd.Flags().GetInt("horizontal-offset")
		if err != nil {
			log.Fatal(err)
		}
		verticalOffset, err := cmd.Flags().GetInt("vertical-offset")
		if err != nil {
			log.Fatal(err)
		}

		noSerialPort, err := cmd.Flags().GetBool("no-serial-port")
		if err != nil {
			log.Fatal(err)
		}
		m, err := flipdot.NewMessage(
			text,
			"text_13px_bold",
			horizontalOffset,
			verticalOffset,
		)
		if err != nil {
			log.Fatal(err)
		}

		var port io.Writer
		if noSerialPort {
			// If flag is passed to disable serial port behaviour, use a discarder instead.
			port = io.Discard
		} else {
			// unclosed port happens here, should make a custom discarder that is a io.WriteCloser instead.
			port, err = serialport.GetPort()
			if err != nil {
				log.Fatal(err)
			}
		}

		f := flipdot.NewFlipdot(112, 19, 0x07, port)
		err = f.Send(m)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	textCmd.Flags().IntP("horizontal-offset", "o", 0, "Positive number denoting the horizontal offset")
	textCmd.Flags().IntP("vertical-offset", "v", 0, "Positive number denoting the vertical offset")
	textCmd.Flags().Bool("no-serial-port", false, "Do not attempt to connect with a serial port when starting the program")
	rootCmd.AddCommand(textCmd)
}
