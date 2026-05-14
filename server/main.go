package main

import (
	// "io"
	"errors"
	"fmt"
	"log"
	"slices"
	"sync"

	// "net"
	"net/http"

	"github.com/gorilla/websocket"
)

const PORT  = 6767
const MOUNT = "/ws"
const WS_READ_LIMIT = 64
var conns connPool

type connPool struct {
    conns []*websocket.Conn
    mu    sync.Mutex
}

func (p *connPool) addConn(c *websocket.Conn) {
    p.conns = append(p.conns, c)
}

func (p *connPool) removeConn(c *websocket.Conn) {
    if i := slices.Index(p.conns, c); i >= 0 {
        p.conns[i] = p.conns[len(p.conns) - 1]
        p.conns = p.conns[:len(p.conns) - 1]
    }
}

func (p *connPool) count() int {
    return len(p.conns)
}

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

func setupConn(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return nil, errors.New("upgrade failed:" + err.Error())
    }

    if conns.count() > 0 {
        closeConn(c, "queue is full")
        return nil, errors.New("connection ended: queues is full")
    }

    c.SetReadLimit(WS_READ_LIMIT)

    conns.addConn(c)

    return c, nil
}

func closeConn(c * websocket.Conn, reason string) {
    c.WriteMessage(
        websocket.CloseMessage,
        websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason),
    )
    c.Close()
    conns.removeConn(c)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("ws connection from", r.RemoteAddr)

    c, err := setupConn(w, r)
    if err != nil {
        log.Println(err)
    }

    defer closeConn(c, "twas good playing")

    for {
        messageType, p, err := c.ReadMessage()
        if err != nil {
            if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
                break
            }

            log.Println("read error:", err)
            break
        }

        if err := c.WriteMessage(messageType, p); err != nil {
            log.Println("write error:", err)
            break
        }
    }
}

func main() {
	address := fmt.Sprintf(":%d", PORT)

    log.Println("server listening on", PORT)

    http.HandleFunc(MOUNT, wsHandler)
    log.Fatal(http.ListenAndServe(address, nil))
}
