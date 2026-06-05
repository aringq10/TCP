package match

import (
	"log"
	"slices"

	"github.com/aringq10/TCP/server/chess/classical"
	"github.com/aringq10/TCP/server/conn"
	"github.com/aringq10/TCP/server/msg"
)

const MAX_INVL = 5

type Match struct {
    Board *classical.Board
    Players []*Player
    ended bool
}

func NewMatch(board *classical.Board, players ...*Player) *Match {
    return &Match{
        Board: board,
        Players: players,
    }
}

func (m *Match) SendColors() {
    for _, p := range m.Players {
        p.WriteString(msg.MatchColor[p.Color])
    }
}

func (m *Match) Broadcast(data []byte, excluded ...*Player) {
    for _, p := range m.Players {
        if slices.Contains(excluded, p) {
            continue
        }
        p.Write(data)
    }
}

func (m *Match) End(reason string) {
    if m.ended {
        return
    }
    m.ended = true
    for _, p := range m.Players {
        p.Conn.Close(msg.EndOfMatch + " " + reason)
    }
}

func (m *Match) HandleMessage(p *Player, data []byte) {
    if p.SbsqINVL > MAX_INVL {
        m.End("Too many invalid messages from " + p.Color.String())
        return
    }

    l := len(data)

    if l < 4 {
        p.WriteINVL()
        return
    }

    message := string(data[:4])

    switch message {
    case msg.Move:
        if l != 10 {
            p.WriteRJCT()
            break
        }

        move, okMove := classical.ParseMove(p.Color, string(data[5:10]))

        if !okMove {
            p.WriteRJCT()
            return
        }

        if !m.Board.MakeMove(move) {
            p.WriteRJCT()
            return
        }

        log.Print(m.Board.String(classical.White))

        p.WriteACPT()
        m.Broadcast(data[:10], p)

        if m.Board.IsOver() {
            m.End(m.Board.Outcome().String())
        }

    case msg.Resign:
        if l != 4 {
            p.WriteINVL()
            break
        }

        m.Board.SetOutcome(classical.NewOutcome(classical.Resignation, classical.OppositeColor(p.Color)))
        m.End(m.Board.Outcome().String())

    default:
        p.WriteINVL()
    }
}

func HandleMatch(connWhite *conn.Conn, connBlack *conn.Conn) {
    log.Println("Match: started", connWhite.WsConn.RemoteAddr(), connBlack.WsConn.RemoteAddr())
    defer log.Println("Match: ended", connWhite.WsConn.RemoteAddr(), connBlack.WsConn.RemoteAddr())

    connWhite.SignalMatchStart()
    connBlack.SignalMatchStart()

    whiteP := NewPlayer(connWhite, classical.White)
    blackP := NewPlayer(connBlack, classical.Black)
    m := NewMatch(classical.NewBoard(), whiteP, blackP)

    log.Print(m.Board.String(classical.White))

    m.SendColors()
    defer m.End("unexpected end of match")

    for {
        select {
        case data, ok := <-connWhite.OutCh:
            if !ok {
                m.End(classical.NewOutcome(classical.Abandonment, m.Players[1].Color).String())
                return
            }
            m.HandleMessage(m.Players[0], data)
        case data, ok := <-connBlack.OutCh:
            if !ok {
                m.End(classical.NewOutcome(classical.Abandonment, m.Players[0].Color).String())
                return
            }
            m.HandleMessage(m.Players[1], data)
        }
    }
}
