package listen

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"mobitec/internal/flipdot"
	"mobitec/internal/serialport"
	"net/url"
	"os"
	"os/signal"
	"time"
	"net/http"

	"github.com/gorilla/websocket"
)

type ApiMessageFrame struct {
	DelayMs int      `json:"delayMs"`
	Pixels  [][]byte `json:"pixels"`
}

type ApiMessage []ApiMessageFrame

func Start(host string, authKey string, noSerialPort bool) {
	var port io.Writer
	var err error
	if noSerialPort {
		port = io.Discard
	} else {
		port, err = serialport.GetPort()
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
	}
	f := flipdot.NewFlipdot(112, 19, 0x07, port)

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: host, Path: ""}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"auth-key": []string{authKey} })
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			processMessage(f, message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func processMessage(f *flipdot.Flipdot, message []byte) {
	var apiMessage ApiMessage
	err := json.Unmarshal(message, &apiMessage)
	if err != nil {
		log.Printf("Bad api message %v\n", message)
		return
	}

	if len(apiMessage) == 0 {
		log.Printf("Message is empty %v\n", apiMessage)
		return
	}

	go outputFrames(f, apiMessage)
}

func outputFrames(f *flipdot.Flipdot, apiMessage ApiMessage) {
	for _, frame := range apiMessage {
		time.Sleep(time.Duration(frame.DelayMs) * time.Millisecond)

		_, err := f.SendPixels(frame.Pixels)
		if err != nil {
			log.Fatal(err)
		}
	}
}
