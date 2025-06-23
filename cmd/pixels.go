package cmd

import (
	"errors"
	"io"
	"log"
	"mobitec/internal/flipdot"
	"mobitec/internal/helpers"
	"mobitec/internal/serialport"

	"github.com/spf13/cobra"
)

var pixelsCmd = &cobra.Command{
	Use:   "pixels",
	Short: "set pixels on the mobitec destination board",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal(errors.New("a single argument containing the pixels to display must be provided"))
		}
		stringPixelState := args[0]

		noSerialPort, err := cmd.Flags().GetBool("no-serial-port")
		if err != nil {
			log.Fatal(err)
		}

		debug, err := cmd.Flags().GetBool("debug")
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
		pixelState := f.StringToPixelState(stringPixelState)
		writtenBytes, err := f.SendPixels(pixelState)
		if err != nil {
			log.Fatal(err)
		} else if debug {
			helpers.PrintHex(writtenBytes)
		}
	},
}

func init() {
	pixelsCmd.Flags().Bool("no-serial-port", false, "Do not attempt to connect with a serial port when starting the program")
	pixelsCmd.Flags().Bool("debug", false, "Output the bytes that are written to the serial port")
	rootCmd.AddCommand(pixelsCmd)
}
