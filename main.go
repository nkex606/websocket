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
		go s.start()
		go s.broadcast()

		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			wsHanlder(s, w, r)
		})
		log.Println("Start listening on port 8080...")
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		forever = make(chan struct{})
		u := newUser()
		u.dial()
		<-forever
	}
}
