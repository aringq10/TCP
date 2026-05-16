package conn

import (
	"slices"
	"sync"
)

type ConnPool struct {
    conns []*Conn
    mu    sync.Mutex
}

func (p *ConnPool) AddConn(c *Conn) {
    p.conns = append(p.conns, c)
}

func (p *ConnPool) RemoveConn(c *Conn) {
    if i := slices.Index(p.conns, c); i >= 0 {
        p.conns[i] = p.conns[len(p.conns) - 1]
        p.conns = p.conns[:len(p.conns) - 1]
    }
}

func (p *ConnPool) Get(i int) *Conn {
    return p.conns[i]
}

func (p *ConnPool) Count() int {
    return len(p.conns)
}
