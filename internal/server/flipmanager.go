package server

import (
	"fmt"
	"log"
	"mobitec/internal/flipdot"
	"time"
)

// TODO: Add some manner of cancelling these running go routines when server should stop.
func manageFlipdot(f *flipdot.Flipdot, msgChan chan flipdot.Message) {
	for {
		select {
		case msg := <-msgChan:
			handleMsg(f, msg)
		}
		// Always sleep when done with a msg, leave them up there.
		time.Sleep(time.Second * 2)
	}
}

func handleMsg(f *flipdot.Flipdot, msg flipdot.Message) {
	fmt.Println("Handling msg!")
	err := f.Send(msg)
	if err != nil {
		log.Println("error occurred during msg sending: ", err)
		return
	}

}
