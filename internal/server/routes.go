package server

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"mobitec/internal/flipdot"
	"net/http"
)

type messageSendRequest struct {
	Text             string
	HorizontalOffset int
	VerticalOffset   int
	Font             string
}

func newControlRouter(flipdotMessageChan chan flipdot.Message) chi.Router {
	r := chi.NewRouter()

	r.Get("/fonts", func(resp http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		supportedFonts := []string{
			"text_5px",
			"text_6px",
			"text_7px",
			"text_7px_bold",
			"text_9px",
			"text_9px_bold",
			"text_9px_bolder",
			"text_13px",
			"text_13px_bold",
			"text_13px_bolder",
			"text_13px_boldest",
			"numbers_14px",
			"text_15px",
			"text_16px",
			"text_16px_bold",
			"text_16px_bolder",
			"symbols",
		}
		writeJson(resp, 200, supportedFonts)
	})

	r.Post("/enqueue/text", func(resp http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		var message messageSendRequest
		err := json.NewDecoder(req.Body).Decode(&message)
		if err != nil {
			writeJson(resp, 500, err)
			return
		}

		flipMsg, err := flipdot.NewMessage(message.Text, message.Font, message.HorizontalOffset, message.VerticalOffset)
		if err != nil {
			var invalidMsgErr flipdot.InvalidMessageError
			if errors.As(err, &invalidMsgErr) {
				writeJson(resp, 400, invalidMsgErr.Messages)
				return
			}

			log.Println(err)
			writeResponse(resp, 500, []byte("An error occurred during message validation"))
			return
		}

		select {
		case flipdotMessageChan <- flipMsg:
			writeResponse(resp, 200, []byte("Message enqueued"))
		default:
			writeResponse(resp, 429, []byte("queue is full, come back later."))
		}
	})

	return r
}
