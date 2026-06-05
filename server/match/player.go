package match

import (
    "github.com/aringq10/TCP/server/conn"
    "github.com/aringq10/TCP/server/msg"
	"github.com/aringq10/TCP/server/chess/classical"
)

type Player struct {
    Conn *conn.Conn
    SbsqINVL int // subsequent INVL message responses to Conn
    Color classical.Color
}

func NewPlayer(conn *conn.Conn, color classical.Color) *Player {
    return &Player{Conn: conn, Color: color}
}

func (p *Player) Write(data []byte) error {
    if string(data) == msg.Invalid {
        p.SbsqINVL++
    } else {
        p.SbsqINVL = 0
    }
    return p.Conn.Write([]byte(data))
}

func (p *Player) WriteString(data string) error {
    return p.Write([]byte(data))
}

func (p *Player) WriteINVL() error {
    return p.WriteString(msg.Invalid)
}

func (p *Player) WriteACPT() error {
    return p.WriteString(msg.Accept)
}

func (p *Player) WriteRJCT() error {
    return p.WriteString(msg.Reject)
}
