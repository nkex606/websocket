package main

import (
	"fmt"
	"log"
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
		fmt.Println(s.clients)
	}
}
