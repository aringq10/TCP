package conn

import (
	"errors"
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
        c.Close()
        return nil, errors.New("connection ended: queue is full")
    }

    c.WsConn.SetReadLimit(WS_READ_LIMIT)

    Conns.AddConn(c)

    return c, nil
}

func (c *Conn) Close() error {
    writeErr := c.WsConn.WriteMessage(
        websocket.CloseMessage,
        websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
    )
    closeErr := c.WsConn.Close()

    return errors.Join(writeErr, closeErr)
}

func (c *Conn) Write(data []byte) error {
    return c.WsConn.WriteMessage(websocket.TextMessage, data)
}

func (c *Conn) Read() (messageType int, data []byte, err error) {
    return c.WsConn.ReadMessage()
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
            default:
                // network / protocol error
            }
            break
        }
        c.OutCh <- data
    }
    close(c.OutCh)
}
