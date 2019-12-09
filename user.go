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
		// TODO: notify server disconnect
		close(forever)
	}()

	go u.read()
	go u.send()
}

func (u *user) read() {
	defer u.conn.Close()
	for {
		_, p, err := u.conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		fmt.Print("server: ", string(p))
	}
}

func (u *user) send() {
	defer u.conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if err := u.conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			log.Fatal(err)
		}
		fmt.Print("you: ", string(text))
	}
}
