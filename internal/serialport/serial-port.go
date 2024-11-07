package serialport

import (
	"errors"
	"fmt"
	"go.bug.st/serial"
	"log"
)

func GetPort() (serial.Port, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		return nil, errors.New("no serial ports found")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

	mode := &serial.Mode{
		BaudRate: 4800,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	return serial.Open(ports[1], mode)
}
