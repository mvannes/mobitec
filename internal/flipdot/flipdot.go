package flipdot

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type Flipdot struct {
	width       int
	height      int
	signAddress byte
	port        io.Writer
}

type Message struct {
	Text             string
	Font             string
	HorizontalOffset int
	VerticalOffset   int
}

type InvalidMessageError struct {
	Messages []string
}

func (i InvalidMessageError) Error() string {
	return "invalid message provided, errors were: " + strings.Join(i.Messages, ", ")
}

func NewMessage(text string, font string, horizontalOffset int, verticalOffset int) (Message, error) {
	m := Message{
		Text:             text,
		Font:             font,
		HorizontalOffset: horizontalOffset,
		VerticalOffset:   verticalOffset,
	}

	return m, validateMessage(m)
}

func validateMessage(m Message) error {
	var errs []string
	if len(strings.TrimSpace(m.Text)) == 0 {
		errs = append(errs, "non empty text must be provided")
	}
	if m.HorizontalOffset < 0 {
		errs = append(errs, "horizontal offset must be a positive integer")
	}

	if m.VerticalOffset < 0 {
		errs = append(errs, "vertical offset must be a positive integer")
	}

	_, err := chooseFont(m.Font)
	if err != nil {
		errs = append(errs, "invalid font provided")
	}

	if len(errs) > 0 {
		return InvalidMessageError{Messages: errs}
	}
	return nil
}

func (f Flipdot) Send(msg Message) error {
	msgBytes, err := makeMessage(f, msg)
	if err != nil {
		log.Fatal(err)
	}
	n, err := f.port.Write(msgBytes)
	if err != nil {
		return err
	}
	fmt.Printf("Sent %v bytes\n", n)
	return nil
}

func NewFlipdot(width int, height int, signAddress byte, port io.Writer) *Flipdot {
	return &Flipdot{
		width:       width,
		height:      height,
		signAddress: signAddress,
		port:        port,
	}
}

func makeMessage(f Flipdot, m Message) ([]byte, error) {
	header := makeHeader(f.signAddress, f.width, f.height)
	fontHex, err := chooseFont(m.Font)
	if err != nil {
		return nil, err
	}
	data := []byte{
		0xd2, // Horizontal offset
		byte(m.HorizontalOffset),
		0xd3, // Vertical offset
		byte(m.VerticalOffset),
		0xd4, // Font
		fontHex,
	}

	data = append(data, textToBytes(m.Text)...)
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
	// Because these would obviously conflict with the end of the Message.
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
		return 0x00, errors.New("unknown Font given")
	}
}

func textToBytes(text string) []byte {
	var bytes []byte
	for _, char := range text {
		bytes = append(bytes, byte(char))
	}

	return bytes
}
