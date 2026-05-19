package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <host> <port>\n", os.Args[0])
		os.Exit(2)
	}
	addr := "ws://" + os.Args[1] + ":" + os.Args[2] + "/ws"
	log.Printf("connecting to %s", addr)

	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			mt, payload, err := conn.ReadMessage()
			if err != nil {
				if ce, ok := err.(*websocket.CloseError); ok {
					fmt.Printf("[close]  code=%d reason=%q\n", ce.Code, ce.Text)
					return
				}
				log.Printf("read: %v", err)
				return
			}
			switch mt {
			case websocket.TextMessage:
				fmt.Printf("[text]   payload=%q\n", payload)
			case websocket.BinaryMessage:
				fmt.Printf("[binary] payload=%x\n", payload)
			default:
				fmt.Printf("[mt=%d]  payload=%x\n", mt, payload)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := conn.WriteMessage(websocket.TextMessage, scanner.Bytes()); err != nil {
			log.Printf("write: %v", err)
			break
		}
	}

	<-done
}
