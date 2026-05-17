package conn

import (
	"slices"
	"sync"
)

type ConnPool struct {
    conns []*Conn
    mu    sync.RWMutex
}

func (p *ConnPool) AddConn(c *Conn) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.conns = append(p.conns, c)
}

func (p *ConnPool) RemoveConn(c *Conn) {
    p.mu.Lock()
    defer p.mu.Unlock()
    if i := slices.Index(p.conns, c); i >= 0 {
        p.conns[i] = p.conns[len(p.conns) - 1]
        p.conns = p.conns[:len(p.conns) - 1]
    }
}

func (p *ConnPool) Get(i int) *Conn {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return p.conns[i]
}

func (p *ConnPool) Count() int {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return len(p.conns)
}

func (p *ConnPool) TryDequeue2() (p1 *Conn, p2 *Conn, ok bool) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if len(p.conns) < 2 {
        return nil, nil, false
    }

    p1 = p.conns[0]
    p2 = p.conns[1]
    p.conns = p.conns[2:]

    return p1, p2, true
}
