package conn

import (
	"errors"
	"log"
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
    SbsqINVL int // subsequent INVL message responses to WsConn
    OutCh chan []byte
    MatchStartCh chan struct{}
}

func SetupConn(w http.ResponseWriter, r *http.Request) (*Conn, error) {
    wsConn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return nil, errors.New("upgrade failed:" + err.Error())
    }

    c := &Conn{WsConn: wsConn, OutCh: make(chan []byte), MatchStartCh: make(chan struct{})}

    if Conns.Count() >= QUEUE_LIMIT {
        c.Close("queue is full")
        return nil, errors.New("connection ended: queue is full")
    }

    c.WsConn.SetReadLimit(WS_READ_LIMIT)

    Conns.AddConn(c)

    return c, nil
}

func (c *Conn) Write(data []byte) error {
    if string(data) == "INVL" {
        c.SbsqINVL++
    } else {
        c.SbsqINVL = 0
    }
    return c.WsConn.WriteMessage(websocket.TextMessage, data)
}

func (c *Conn) WriteString(data string) error {
    return c.Write([]byte(data))
}

func (c *Conn) WriteINVL() error {
    return c.WriteString("INVL")
}

func (c *Conn) WriteACPT() error {
    return c.WriteString("ACPT")
}

func (c *Conn) WriteRJCT() error {
    return c.WriteString("RJCT")
}

func (c *Conn) Read() (messageType int, data []byte, err error) {
    return c.WsConn.ReadMessage()
}

func (c *Conn) Close(reason string) error {
    var errors []error
    errors = append(errors,
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

func (c *Conn) SignalMatchStart() {
    c.MatchStartCh <- struct{}{}
}

func (c *Conn) ReadToChan() {
    for {
        _, data, err := c.Read()
        if err != nil {
            switch {
            case errors.Is(err, websocket.ErrReadLimit):
                // client sent oversized message
            case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
                // expected disconnect
            case websocket.IsUnexpectedCloseError(err):
                // peer left abnormally — worth logging
            default:
                // network / protocol error — also bail
                log.Println("reading error:", err)
            }
            break
        }
        c.OutCh <- data
    }
    close(c.OutCh)
}
