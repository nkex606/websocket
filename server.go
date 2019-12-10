package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

type server struct {
	clients      map[*client]struct{}
	connected    chan *client
	disconnected chan *client
}

func newServer() *server {
	s := &server{
		clients:      make(map[*client]struct{}),
		connected:    make(chan *client),
		disconnected: make(chan *client),
	}
	return s
}

func (s *server) start() {
	for {
		select {
		case client := <-s.connected:
			log.Println("A client connect.")
			s.clients[client] = struct{}{}
		case client := <-s.disconnected:
			log.Println("A client disconnect.")
			delete(s.clients, client)
		}
	}
}

func (s *server) broadcast() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		for client := range s.clients {
			if err := client.conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
				log.Fatal("broadcast err: ", err.Error())
				break
			}
		}
		fmt.Print("broadcast: ", string(text))
	}
}
