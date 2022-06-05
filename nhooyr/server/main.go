package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

// This is a simple websocket echo server.
// When a message is received from a client
// the same message is sent back.
func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
	if err != nil {
		fmt.Printf("error accepting websocket handshake: %+v\n", err)
	}
	defer conn.Close(websocket.StatusInternalError, "connection closed")

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	// receive message
	_, received, err := conn.Read(ctx)
	if err != nil {
		fmt.Printf("error reading: %+v\n", err)
	}

	fmt.Printf("received message from client: %+v\n", string(received))

	// send message
	if err := conn.Write(ctx, websocket.MessageText, received); err != nil {
		fmt.Println("error sending\n", err)
	}

	conn.Close(websocket.StatusNormalClosure, "")
}

func main() {
	http.Handle("/", http.HandlerFunc(echo))

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
