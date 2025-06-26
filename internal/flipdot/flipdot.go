package flipdot

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type PixelState = [][]byte

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

func (f Flipdot) StringToPixelState(input string) PixelState {
	rows := make([][]byte, f.height)
	columnCount := f.width

	inputRows := strings.Fields(input)
	for rowIndex := 0; rowIndex < len(rows); rowIndex++ {
		if len(rows[rowIndex]) == 0 {
			rows[rowIndex] = make([]byte, columnCount)
		}

		if len(inputRows) < rowIndex+1 {
			continue
		}

		inputRowRunes := []rune(inputRows[rowIndex])
		for runeIndex := 0; runeIndex < len(inputRowRunes); runeIndex++ {
			if runeIndex+1 > len(rows[rowIndex]) {
				break
			}

			switch inputRowRunes[runeIndex] {
			case '1':
				rows[rowIndex][runeIndex] = 1
			case '0':
			default:
				rows[rowIndex][runeIndex] = 0
			}
		}
	}

	return rows
}

func (f Flipdot) SendPixels(pixelState PixelState) ([]byte, error) {
	msgBytes, err := makePixelMessage(f, pixelState)
	if err != nil {
		log.Fatal(err)
	}
	n, err := f.port.Write(msgBytes)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Sent %v bytes\n", n)
	return msgBytes, nil
}

func (f Flipdot) SendText(msg Message) ([]byte, error) {
	msgBytes, err := makeTextMessage(f, msg)
	if err != nil {
		log.Fatal(err)
	}
	n, err := f.port.Write(msgBytes)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Sent %v bytes\n", n)
	return msgBytes, nil
}

func NewFlipdot(width int, height int, signAddress byte, port io.Writer) *Flipdot {
	return &Flipdot{
		width:       width,
		height:      height,
		signAddress: signAddress,
		port:        port,
	}
}

func makePixelMessage(f Flipdot, pixelState PixelState) ([]byte, error) {
	dataSections, err := pixelStateToBitwiseDataSections(pixelState)
	if err != nil {
		log.Fatal(err)
	}
	return makeMessage(f, dataSections)
}

func makeTextMessage(f Flipdot, m Message) ([]byte, error) {
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

	var dataSections [][]byte
	dataSections = append(dataSections, data)

	return makeMessage(f, dataSections)
}

func makeMessage(f Flipdot, dataSections [][]byte) ([]byte, error) {
	header := makeHeader(f.signAddress, f.width, f.height)

	var data []byte
	for _, dataSection := range dataSections {
		data = append(data, dataSection...)
	}

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

func pixelStateToBitwiseDataSections(pixelState PixelState) ([][]byte, error) {
	const charColumnSize = 5

	fontHex, err := chooseFont("bitwise")
	if err != nil {
		return nil, err
	}

	var lines [][]byte
	for lineIndex := 0; lineIndex < len(pixelState); lineIndex += charColumnSize {
		var lineBody []byte

		emptyColumnsStart := 0
		emptyColumnChar := columnToBitwiseChar(0, 0, 0, 0, 0)
		for columnIndex := 0; columnIndex < len(pixelState[lineIndex]); columnIndex++ {
			r1 := getPixelFromPixelState(pixelState, lineIndex+0, columnIndex)
			r2 := getPixelFromPixelState(pixelState, lineIndex+1, columnIndex)
			r3 := getPixelFromPixelState(pixelState, lineIndex+2, columnIndex)
			r4 := getPixelFromPixelState(pixelState, lineIndex+3, columnIndex)
			r5 := getPixelFromPixelState(pixelState, lineIndex+4, columnIndex)
			char := columnToBitwiseChar(r1, r2, r3, r4, r5)

			if char == emptyColumnChar && columnIndex == emptyColumnsStart {
				emptyColumnsStart++
			} else {
				lineBody = append(lineBody, char)
			}
		}

		// Trim end
		for {
			if len(lineBody) > 0 && lineBody[len(lineBody)-1] == emptyColumnChar {
				lineBody = lineBody[:len(lineBody)-1]
			} else {
				break
			}
		}

		lineHeader := []byte{
			0xd2, // Horizontal offset
			byte(emptyColumnsStart),
			0xd3,                // Vertical offset
			byte(lineIndex + 4), // 4 is initial offset? Offset is from bottom?
			0xd4,                // Font
			fontHex,
		}

		lines = append(lines, lineHeader, lineBody)
	}

	return lines, nil
}

func getPixelFromPixelState(pixelState PixelState, row int, column int) byte {
	if len(pixelState) >= row+1 && len(pixelState[row]) >= column+1 {
		return pixelState[row][column]
	}
	return 0
}

func columnToBitwiseChar(r1, r2, r3, r4, r5 byte) byte {
	var charCode byte = 32

	if r1 == 1 {
		charCode += 1
	}

	if r2 == 1 {
		charCode += 2
	}

	if r3 == 1 {
		charCode += 4
	}

	if r4 == 1 {
		charCode += 8
	}

	if r5 == 1 {
		charCode += 16
	}

	return charCode
}

func textToBytes(text string) []byte {
	var bytes []byte
	for _, char := range text {
		bytes = append(bytes, byte(char))
	}

	return bytes
}
