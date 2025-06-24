package server

import (
	"log"
	"mobitec/internal/flipdot"
	"time"
)

type FlipdotChange struct {
	text       *flipdot.Message
	pixelState *flipdot.PixelState
}

// TODO: Add some manner of cancelling these running go routines when server should stop.
func manageFlipdot(f *flipdot.Flipdot, msgChan chan FlipdotChange) {
	for {
		select {
		case msg := <-msgChan:
			handleMsg(f, msg)
		}
		// Always sleep when done with a msg, leave them up there.
		time.Sleep(time.Second * 2)
	}
}

func handleMsg(f *flipdot.Flipdot, msg FlipdotChange) {
	var err error
	if msg.text != nil {
		_, sendErr := f.SendText(*msg.text)
		err = sendErr
	} else if msg.pixelState != nil {
		_, sendErr := f.SendPixels(*msg.pixelState)
		err = sendErr
	} else {
		log.Println("error occurred during msg sending: FlipdotChange has no pixelState or text")
		return
	}

	if err != nil {
		// TODO: create some form of err handling that is either better logging or
		// TODO: cancelling of the server as it is.
		// TODO: Either way, from context of a go routine there is little to do here right now.
		log.Println("error occurred during msg sending: ", err)
		return
	}

}
