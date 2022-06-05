package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

// This is a simple websocket echo server.
// When a message is received from a client
// the same message is sent back.
func Echo(ws *websocket.Conn) {

	for {
		var msg string

		if err := websocket.Message.Receive(ws, &msg); err != nil {
			fmt.Printf("error receiving: %+v\n", err)
			break
		}

		fmt.Printf("received message from client: %+v\n", msg)

		fmt.Println("sending message to client")
		err := websocket.Message.Send(ws, msg)
		if err == nil {
			fmt.Println("message sent ok")
			break
		}
		fmt.Printf("error sending: %+v", err)
	}
}

func main() {
	http.Handle("/", websocket.Handler(Echo))

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
