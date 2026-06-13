package match

import (
	"time"

	"github.com/aringq10/TCP/server/chess/classical"
	"github.com/aringq10/TCP/server/conn"
)

type Player struct {
    Conn *conn.Conn
    SbsqINVL int // subsequent INVL message responses to Conn
    Timer *time.Timer
    Color classical.Color
    TimeRemaining time.Duration
}

func NewPlayer(conn *conn.Conn,color classical.Color, timeRemaining time.Duration) *Player {
    return &Player{Conn: conn, Color: color, TimeRemaining: timeRemaining}
}

func (p *Player) Write(message string) error {
    return p.Conn.Write([]byte(message))
}
