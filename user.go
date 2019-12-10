package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type user struct {
	conn *websocket.Conn
}

func newUser() *user {
	u := &user{}
	return u
}

func (u *user) dial() {
	var err error
	u.conn, _, err = websocket.DefaultDialer.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	log.Println("Connected!")

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		log.Println("===>", <-sigint)
		u.conn.Close()
		close(forever)
	}()

	go u.read()
	go u.send()
}

func (u *user) read() {
	for {
		_, p, err := u.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Print("server: ", string(p))
	}
}

func (u *user) send() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if err := u.conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			break
		}
		fmt.Print("you: ", string(text))
	}
}
