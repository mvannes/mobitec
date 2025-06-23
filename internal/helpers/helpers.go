package helpers

import (
	"fmt"
	"log"
)

func PrintHex(bytes []byte) {
	log.Print("Written bites")

	for i := 0; i < len(bytes); i++ {
		fmt.Printf("0x%X (%d)\n", bytes[i], bytes[i])
	}
}
