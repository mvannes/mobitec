package main

import (
	"log"
	"mobitec/cmd"
)

func main() {

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
