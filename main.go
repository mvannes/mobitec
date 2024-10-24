package main

import (
	"log"
	"mobitec/flipdot"
	serialport "mobitec/serial-port"
	"time"
)

func main() {
	port, err := serialport.GetPort()
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	f := flipdot.NewFlipdot(112, 19, 0x07, port)

	messages := []string{"SISYPHOS", "PYTHIOS", "PYTHSYPHOS"}
	currentMsg := 0
	for {
		currentMsg++
		if currentMsg > 2 {
			currentMsg = 0
		}
		m := flipdot.Message{
			Text:             messages[currentMsg],
			Font:             "text_13px_bold",
			HorizontalOffset: 9,
			VerticalOffset:   19,
		}

		err := f.Send(m)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 60)
	}
}
