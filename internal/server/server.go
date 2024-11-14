package server

import (
	_ "embed"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"mobitec/internal/flipdot"
	"mobitec/internal/serialport"
	"net/http"
)

//go:embed assets/index.html
var uiFile string

func Start() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}}))

	//var err error
	//TODO: Configure option to start with this serialport configured instead of commenting out.
	//TODO: @see text.go flag handling.
	port, err := serialport.GetPort()
	if err != nil {
		return err
	}

	//port := io.Discard

	// TODO: make configurable.
	f := flipdot.NewFlipdot(112, 19, 0x07, port)

	// Make this smarter, because now the queue is maxed at 200, which seems like a lot.
	msgChan := make(chan flipdot.Message, 50)
	r.Mount("/flipdot", newControlRouter(msgChan))
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte(uiFile))
	})
	go manageFlipdot(f, msgChan)

	return http.ListenAndServe(":8080", r)
}

func writeJson(w http.ResponseWriter, status int, body any) {
	respBody, err := json.Marshal(body)

	if err != nil {
		writeResponse(w, 500, []byte("An error occurred"))
		return
	}
	writeResponse(w, status, respBody)
}

func writeResponse(w http.ResponseWriter, status int, body []byte) {
	w.WriteHeader(status)
	_, err := w.Write(body)
	if err != nil {
		// TODO: some structured logging goes here.
		log.Println(err)
	}
}
