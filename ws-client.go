package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	log.Println("Connected!")
	defer conn.Close()

	forever := make(chan struct{})

	go reader(conn)
	go sender(conn)

	<-forever

}

func reader(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		fmt.Print("server: ", string(p))
	}
}

func sender(conn *websocket.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {

		text, _ := reader.ReadString('\n')
		if err := conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			log.Fatal(err)
		}
		fmt.Print("you: ", string(text))
	}
}

