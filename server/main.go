package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/aringq10/TCP/server/conn"
)

const PORT  = 6767
const MOUNT = "/ws"

func isMoveValid() bool {
    return true
}

func connToChan(c *websocket.Conn, ch chan []byte) {
    for {
        _, data, err := c.ReadMessage()

        if err != nil {
            fmt.Println("reading error:", err)
            break
        }

        ch <- data
    }

    close(ch)
}

func handleMatch(connWhite *websocket.Conn, connBlack *websocket.Conn) {
    conn.Conns.RemoveConn(connWhite)
    conn.Conns.RemoveConn(connBlack)

    defer conn.CloseConn(connWhite, "twas good playing")
    defer conn.CloseConn(connBlack, "twas good playing")

    fmt.Println("Match started between", connWhite.RemoteAddr(), "with", connBlack.RemoteAddr())

    chWhite := make(chan []byte)
    chBlack := make(chan []byte)

    go connToChan(connWhite, chWhite)
    go connToChan(connBlack, chBlack)

    var data []byte
    var ok bool
    var matchOnGoing bool = true
	var whiteMove bool = true

    for matchOnGoing {
        select {
        case data, ok = <-chWhite:
            if !ok {
                matchOnGoing = false
                break
            }
            if len(data) < 4 {
                connWhite.WriteMessage(websocket.TextMessage, []byte("INVL"))
                continue
            }
            switch string(data[:4]) {
            case "MOVE":
                if len(data) < 10 {
                    connWhite.WriteMessage(websocket.TextMessage, []byte("INVL"))
                    continue
                }
                from := data[5:7]
                to := data[8:10]

                if isMoveValid() && whiteMove {
					whiteMove = false
					fmt.Println("WHITE MOVED", from, to)
                    connWhite.WriteMessage(websocket.TextMessage, []byte("ACPT"))
                    connBlack.WriteMessage(websocket.TextMessage, data[:10])
                } else {
                    connWhite.WriteMessage(websocket.TextMessage, []byte("RJCT"))
                }
            default:
                connWhite.WriteMessage(websocket.TextMessage, []byte("INVL"))
            }
        case data, ok = <-chBlack:
            if !ok {
                matchOnGoing = false
                break
            }
            if len(data) < 4 {
                connBlack.WriteMessage(websocket.TextMessage, []byte("INVL"))
                continue
            }
            switch string(data[:4]) {
            case "MOVE":
                if len(data) < 10 {
                    connBlack.WriteMessage(websocket.TextMessage, []byte("INVL"))
                    continue
                }
                from := data[5:7]
                to := data[8:10]
                
                if isMoveValid() && !whiteMove {
					whiteMove = true
					fmt.Println("BLACK MOVED", from, to)
                    connBlack.WriteMessage(websocket.TextMessage, []byte("ACPT"))
                    connWhite.WriteMessage(websocket.TextMessage, data[:10])
                } else {
                    connBlack.WriteMessage(websocket.TextMessage, []byte("RJCT"))
                }
            default:
                connBlack.WriteMessage(websocket.TextMessage, []byte("INVL"))
            }
        }
    }

    fmt.Println("Match ended between", connWhite.RemoteAddr(), "with", connBlack.RemoteAddr())
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("ws connection from", r.RemoteAddr)

    err := conn.SetupConn(w, r)
    if err != nil {
        log.Println(err)
        return
    }

    if conn.Conns.Count() >= 2 {
        go handleMatch(conn.Conns.Get(0), conn.Conns.Get(1))
    }
}

func main() {
	address := fmt.Sprintf(":%d", PORT)

    log.Println("server listening on", PORT)

    http.HandleFunc(MOUNT, wsHandler)
    log.Fatal(http.ListenAndServe(address, nil))
}
