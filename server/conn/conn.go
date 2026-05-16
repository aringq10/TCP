package conn

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

const WS_READ_LIMIT = 64
const QUEUE_LIMIT = 2

var Conns ConnPool

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // allow all origins in dev; tighten in prod
    },
}

type Conn struct {
    WsConn *websocket.Conn
    SbsqInvl int // subsequent INVL message responses to WsConn
}

func SetupConn(w http.ResponseWriter, r *http.Request) (*Conn, error) {
    wsConn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return nil, errors.New("upgrade failed:" + err.Error())
    }

    c := &Conn{WsConn: wsConn}

    if Conns.Count() >= QUEUE_LIMIT {
        Conns.RemoveConn(c)
        c.Close("queue is full")
        return nil, errors.New("connection ended: queue is full")
    }

    c.WsConn.SetReadLimit(WS_READ_LIMIT)

    Conns.AddConn(c)

    return c, nil
}

func (c *Conn) Write(data []byte) error {
    if string(data) == "INVL" {
        c.SbsqInvl++
    } else {
        c.SbsqInvl = 0
    }
    return c.WsConn.WriteMessage(websocket.TextMessage, data)
}

func (c *Conn) WriteString(data string) error {
    return c.Write([]byte(data))
}

func (c *Conn) WriteInvl() error {
    return c.WriteString("INVL")
}

func (c *Conn) Read() (messageType int, data []byte, err error) {
    return c.WsConn.ReadMessage()
}

func (c *Conn) Close(reason string) error {
    var errors []error
    errors = append(errors,
        c.WriteString(reason),
        c.WsConn.WriteMessage(
            websocket.CloseMessage,
            websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason),
        ),
        c.WsConn.Close(),
    )

    for _, e := range errors {
        if e != nil {
            return e
        }
    }

    return nil
}

func (c *Conn) ReadToChan(ch chan<- []byte) {
    for {
        _, data, err := c.Read()
        if err != nil {
            fmt.Println("reading error:", err)
            break
        }
        ch <- data
    }
    close(ch)
}
