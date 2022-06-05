package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nhooyr.io/websocket"
)

type Message struct {
	Message string `json:"message"`
}

func main() {

	msg := Message{Message: "hello"}
	json, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error marshalling json\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws://localhost:1234", &websocket.DialOptions{})
	if err != nil {
		fmt.Println("error dialing\n", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "connection closed")

	// send message
	if err := conn.Write(ctx, websocket.MessageText, json); err != nil {
		fmt.Println("error sending\n", err)
	}

	// receive message
	_, received, err := conn.Read(ctx)

	fmt.Println("received from server: ", received)

	fmt.Printf("received reply from server: %+v\n", string(received))
	conn.Close(websocket.StatusNormalClosure, "")

}
