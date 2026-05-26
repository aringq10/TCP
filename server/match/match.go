package match

import (
	"log"
	"slices"

	"github.com/aringq10/TCP/server/chess"
	"github.com/aringq10/TCP/server/chess/classical"
	"github.com/aringq10/TCP/server/conn"
)

const MAX_INVL = 5

type Player struct {
    Conn *conn.Conn
    Color classical.Color
}

func NewPlayer(conn *conn.Conn, color classical.Color) *Player {
    return &Player{Conn: conn, Color: color}
}

type Match struct {
    Board chess.Board
    Players []*Player
}

func NewMatch(board chess.Board, players ...*Player) *Match {
    return &Match{
        Board: board,
        Players: players,
    }
}

func (m *Match) SendColors() {
    for _, p := range m.Players {
        p.Conn.WriteString(p.Color.String())
    }
}

func (m *Match) Broadcast(data []byte, excluded ...*Player) {
    for _, p := range m.Players {
        if slices.Contains(excluded, p) {
            continue
        }
        p.Conn.Write(data)
    }
}

func (m *Match) End(reason string) {
    for _, p := range m.Players {
        p.Conn.Close("EOM: " + reason)
    }
}

func (m *Match) HandleMessage(p *Player, data []byte) {
    if p.Conn.SbsqINVL > MAX_INVL {
        m.End("Too many invalid messages from " + p.Color.String())
        return
    }

    l := len(data)

    if l < 4 {
        p.Conn.WriteINVL()
        return
    }

    msg := string(data[:4])

    switch msg {
    case "MOVE":
        if l < 10 {
            p.Conn.WriteINVL()
            break
        }

        from, okFrom := classical.ParseSquare(string(data[5:7]))
        to,   okTo   := classical.ParseSquare(string(data[8:10]))

        if !okFrom || !okTo {
            p.Conn.WriteRJCT()
            return
        }

        if !m.Board.MakeMove(classical.NewMove(p.Color, from, to)) {
            p.Conn.WriteRJCT()
            return
        }

        log.Print(m.Board.String())

        p.Conn.WriteACPT()
        m.Broadcast(data[:10], p)

    default:
        p.Conn.WriteINVL()
    }
}

func HandleMatch(connWhite *conn.Conn, connBlack *conn.Conn) {
    log.Println("Match: started", connWhite.WsConn.RemoteAddr(), connBlack.WsConn.RemoteAddr())
    defer log.Println("Match: ended", connWhite.WsConn.RemoteAddr(), connBlack.WsConn.RemoteAddr())

    connWhite.SignalMatchStart()
    connBlack.SignalMatchStart()

    whiteP := NewPlayer(connWhite, classical.White)
    blackP := NewPlayer(connBlack, classical.Black)
    m := NewMatch(classical.NewPawnsOnlyBoard(), whiteP, blackP)

    log.Print(m.Board.String())

    m.SendColors()
    defer m.End("unexpected end of match")

    for {
        select {
        case data, ok := <-connWhite.OutCh:
            if !ok {
                m.End(m.Players[0].Color.String() + " disconnected")
                return
            }
            m.HandleMessage(m.Players[0], data)
        case data, ok := <-connBlack.OutCh:
            if !ok {
                m.End(m.Players[1].Color.String() + " disconnected")
                return
            }
            m.HandleMessage(m.Players[1], data)
        }
    }
}
