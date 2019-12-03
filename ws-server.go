package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsHanlder)
	log.Println("Start listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsHanlder(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("get conn fail, err: ", err.Error())
	}
	log.Println("Client Connected")

	go reader(conn)
	go sender(conn)
}

func reader(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			panic(err)
			// return
		}
		fmt.Print("client: ", string(p))
	}
}

func sender(conn *websocket.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if err := conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			log.Fatal(err)
			// return
		}
		fmt.Print("you: ", string(text))
	}
}

