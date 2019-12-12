package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type client struct {
	serv *server
	conn *websocket.Conn
	send chan []byte
}

func (c *client) reader() {
	defer func() {
		c.serv.disconnected <- c
		c.conn.Close()
	}()
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Print("client: ", string(p))
		c.serv.broadcastCh <- p
	}
}

func (c *client) writer() {
	defer c.conn.Close()
	for {
		select {
		case msg := <-c.send:
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}
		}
	}
}

func wsHanlder(serv *server, w http.ResponseWriter, r *http.Request) {
	// server side conn
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	c := &client{
		serv: serv,
		send: make(chan []byte),
		conn: conn,
	}
	c.serv.connected <- c

	go c.reader()
	go c.writer()
}
