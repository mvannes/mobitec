package main

import (
	"errors"
	"fmt"
	"go.bug.st/serial"
	"log"
	"time"
)

func main() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
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
	port, err := serial.Open(ports[0], mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	f := flipdot{
		width:       112,
		height:      19,
		signAddress: 0x07,
	}

	messages := []string{"SISYPHOS", "PYTHIOS", "PYTHSYPHOS"}

	currentMsg := 0
	for {
		currentMsg++
		if currentMsg > 2 {
			currentMsg = 0
		}
		m := message{
			text:             messages[currentMsg],
			font:             "text_13px_bold",
			horizontalOffset: 9,
			verticalOffset:   19,
		}

		msg, err := makeMessage(f, m)
		if err != nil {
			log.Fatal(err)
		}
		n, err := port.Write(msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %v bytes\n", n)
		time.Sleep(time.Second * 10)

	}

}

type flipdot struct {
	width       int
	height      int
	signAddress byte
}

type message struct {
	text             string
	font             string
	horizontalOffset int
	verticalOffset   int
}

func makeMessage(f flipdot, m message) ([]byte, error) {
	header := makeHeader(f.signAddress, f.width, f.height)
	fontHex, err := chooseFont(m.font)
	if err != nil {
		return nil, err
	}
	data := []byte{
		0xd2, // Horizontal offset
		byte(m.horizontalOffset),
		0xd3, // Vertical offset
		byte(m.verticalOffset),
		0xd4, // Font
		fontHex,
	}

	data = append(data, textToBytes(m.text)...)
	footer := makeFooter(header, data)
	return append(append(header, data...), footer...), nil
}

func makeHeader(signAddres byte, width int, height int) []byte {
	return []byte{
		0xff,       // Starting byte
		signAddres, // Sign address
		0xa2,       // Always a2
		0xd0,       // width marker
		byte(width),
		0xd1, // height marker
		byte(height),
	}
}

func makeFooter(header []byte, data []byte) []byte {
	checkSum := 0

	for _, b := range append(header[1:], data...) {
		checkSum += int(b)
	}

	var checkSumBytes []byte
	checkSumByte := byte(checkSum & 0xff)

	// Some bytes are handled differently
	// The stop byte and what the stop byte turns into.
	// Because these would obviously conflict with the end of the message.
	if checkSumByte == 0xff {
		checkSumBytes = []byte{0xfe, 0x01}
	} else if checkSumByte == 0xfe {
		checkSumBytes = []byte{0xfe, 0x00}
	} else {
		checkSumBytes = []byte{checkSumByte}
	}

	// Stop byte at the end.
	return append(checkSumBytes, 0xff)
}

func chooseFont(font string) (byte, error) {
	switch font {
	case "text_5px":
		return 0x72, nil
	case "text_6px":
		return 0x66, nil
	case "text_7px":
		return 0x65, nil
	case "text_7px_bold":
		return 0x64, nil
	case "text_9px":
		return 0x75, nil
	case "text_9px_bold":
		return 0x70, nil
	case "text_9px_bolder":
		return 0x62, nil
	case "text_13px":
		return 0x73, nil
	case "text_13px_bold":
		return 0x69, nil
	case "text_13px_bolder":
		return 0x61, nil
	case "text_13px_boldest":
		return 0x79, nil
	case "numbers_14px":
		return 0x00, nil
	case "text_15px":
		return 0x71, nil
	case "text_16px":
		return 0x68, nil
	case "text_16px_bold":
		return 0x78, nil
	case "text_16px_bolder":
		return 0x74, nil
	case "symbols":
		return 0x67, nil
	case "bitwise":
		return 0x77, nil
	default:
		return 0x00, errors.New("unknown font given")
	}
}

func textToBytes(text string) []byte {
	var bytes []byte
	for _, char := range text {
		bytes = append(bytes, byte(char))
	}

	return bytes
}
