package conn

import (
	"errors"
	"github.com/gorilla/websocket"
	"slices"
	"sync"
    "net/http"
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

type ConnPool struct {
    conns []*websocket.Conn
    mu    sync.Mutex
}

func (p *ConnPool) AddConn(c *websocket.Conn) {
    p.conns = append(p.conns, c)
}

func (p *ConnPool) RemoveConn(c *websocket.Conn) {
    if i := slices.Index(p.conns, c); i >= 0 {
        p.conns[i] = p.conns[len(p.conns) - 1]
        p.conns = p.conns[:len(p.conns) - 1]
    }
}

func (p *ConnPool) Get(i int) *websocket.Conn {
    return p.conns[i]
}

func (p *ConnPool) Count() int {
    return len(p.conns)
}

func SetupConn(w http.ResponseWriter, r *http.Request) error {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return errors.New("upgrade failed:" + err.Error())
    }

    if Conns.Count() >= QUEUE_LIMIT {
        Conns.RemoveConn(c)
        CloseConn(c, "queue is full")
        return errors.New("connection ended: queue is full")
    }

    c.SetReadLimit(WS_READ_LIMIT)

    Conns.AddConn(c)

    return nil
}

func CloseConn(c *websocket.Conn, reason string) {
    c.WriteMessage(
        websocket.CloseMessage,
        websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason),
    )
    c.Close()
}
