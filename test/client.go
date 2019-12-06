package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type client struct {
	serv *server
	send chan []byte
	conn *websocket.Conn
}

func newClient() *client {
	c := &client{
		send: make(chan []byte),
	}
	return c
}

func (c *client) connect() {
	var err error
	c.conn, _, err = websocket.DefaultDialer.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	log.Println("Connected!")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	go func() {
		<-sigint
		log.Println("Quit chat")

		close(forever)
	}()
}

func (c *client) read() {
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		fmt.Print("server: ", string(p))
	}
}

func wsHanlder(serv *server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := &client{
		serv: serv,
		send: make(chan []byte),
		conn: conn,
	}
	serv.connected <- client

	go client.read()
}
