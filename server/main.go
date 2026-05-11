package main

import (
	// "io"
	"fmt"
	"log"
	// "net"
	"net/http"

	"github.com/gorilla/websocket"
)

const PORT = 6767

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // allow all origins in dev; tighten in prod
    },
}

// func connToChan(c net.Conn, ch chan []byte) {
//     b := make([]byte, 1024)
//     for {
//         n, err := c.Read(b)
//         if err == io.EOF {
//             break
//         }
//         if err != nil {
//             fmt.Println("reading error:", err)
//             continue
//         }
//
//         out := make([]byte, n)
//         copy(out, b)
//
//         ch <- out
//     }
// }
//
// func handleMatch(connWhite net.Conn, connBlack net.Conn) {
//     defer connWhite.Close()
//     defer connBlack.Close()
//
//     fmt.Println("Connecting", connWhite.RemoteAddr(), "with", connBlack.RemoteAddr())
//
//     chWhite := make(chan []byte)
//     chBlack := make(chan []byte)
//
//     go connToChan(connWhite, chWhite)
//     go connToChan(connBlack, chBlack)
//
//     fmt.Println("Starting echo loop for match")
//     var data []byte
//     for {
//         select {
//         case data = <-chWhite:
//             fmt.Println("Received data from chWhite:", data)
//             connBlack.Write(data)
//         case data = <-chBlack:
//             fmt.Println("Received data from chBlack:", data)
//             connWhite.Write(data)
//         }
//     }
// }

func wsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("ws connection from", r.RemoteAddr)
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("upgrade failed:", err)
        return
    }
    defer conn.Close()

    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
                break
            }
            log.Println("read error:", err)
            break
        }

        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println("write error:", err)
            break
        }
    }
}

func main() {
	address := fmt.Sprintf(":%d", PORT)

    log.Println("server listening on", PORT)

    http.HandleFunc("/ws", wsHandler)
	log.Fatal(http.ListenAndServe(address, nil))
}
