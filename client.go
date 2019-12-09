package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type client struct {
	serv *server
	conn *websocket.Conn
	send chan []byte
}

func connect() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	log.Println("Connected!")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	go func() {
		<-sigint
		conn.Close()
		log.Println("Quit chat")
		close(forever)
	}()

	go read(conn)
	go send(conn)
}

func (c *client) reader() {
	defer func() {
		c.serv.disconnected <- c
		c.conn.Close()
	}()
	for {

		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Fatal("read err: ", err.Error())
			return
		}
		fmt.Print("server: ", string(p))
	}
}

func (c *client) sender() {
	defer func() {
		c.conn.Close()
	}()
	reader := bufio.NewReader(os.Stdin)
	for {

		text, _ := reader.ReadString('\n')
		if err := c.conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			log.Fatal("send err: ", err.Error())
			return
		}
		fmt.Print("you: ", string(text))
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
	go c.sender()
}

func read(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		fmt.Print("server: ", string(p))
	}
}

func send(conn *websocket.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {

		text, _ := reader.ReadString('\n')
		if err := conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			log.Fatal(err)
		}
		fmt.Print("you: ", string(text))
	}
}
