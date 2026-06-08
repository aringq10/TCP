package match

import (
	"time"

	"github.com/aringq10/TCP/server/chess/classical"
	"github.com/aringq10/TCP/server/conn"
	"github.com/aringq10/TCP/server/msg"
)

type Player struct {
    Conn *conn.Conn
    SbsqINVL int // subsequent INVL message responses to Conn
    Color classical.Color
    Timer *time.Timer
    TimeRemaining time.Duration
}

func NewPlayer(conn *conn.Conn, color classical.Color, timeRemaining time.Duration) *Player {
    return &Player{Conn: conn, Color: color, TimeRemaining: timeRemaining}
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
