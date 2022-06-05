package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

type Message struct {
	Message string `json:"message"`
}

func main() {

	url := "ws://localhost:1234/"
	conn, err := websocket.Dial(url, "", url)
	if err != nil {
		fmt.Println("error dialing\n", err)
		return
	}
	defer conn.Close()

	msg := Message{Message: "hello"}
	json, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error marshalling json\n", err)
	}

	// send message
	if err := websocket.JSON.Send(conn, string(json)); err != nil {
		fmt.Println("error sending\n", err)
		return
	}

	// receive message
	var message interface{}
	if err := websocket.JSON.Receive(conn, &message); err != nil {
		fmt.Println("error receiving\n", err)
		return
	}

	fmt.Println("received from server: ", message)
}
