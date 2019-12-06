package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var mode string
var forever chan struct{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func init() {
	flag.StringVar(&mode, "m", "", "server or client mode")
}

func main() {
	flag.Parse()

	if mode == "server" {
		s := newServer()
		go s.run()

		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			wsHanlder(s, w, r)
		})
		log.Println("Start listening on port 8080...")
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else {

	}
}