package match

import (
	"log"

	"github.com/aringq10/TCP/server/chess"
	"github.com/aringq10/TCP/server/conn"
)

const MAX_INVL = 5

type Player struct {
    Conn *conn.Conn
    Color chess.Color
}

type Match struct {
    Board *chess.Board
    Players []Player
    WhoseTurn chess.Color
    endSignal chan string
}

func (m *Match) SendColors() {
    for _, p := range m.Players {
        p.Conn.WriteString(p.Color.String())
    }
}

func (m *Match) SignalEnd(reason string) {
    m.endSignal <- reason
}

func (m *Match) End(reason string) {
    for _, p := range m.Players {
        p.Conn.Close("EOM: " + reason)
    }
}

func (m *Match) NextTurn() {
    if m.WhoseTurn == chess.WHITE {
        m.WhoseTurn = chess.BLACK
    } else {
        m.WhoseTurn = chess.WHITE
    }
}

func (m *Match) HandleMessage(p *Player, data []byte) {
    if p.Conn.SbsqInvl > MAX_INVL {
        m.SignalEnd("Too many invalid messages from " + p.Color.String())
        return
    }
    if len(data) < 4 {
        p.Conn.WriteInvl()
        return
    }

    msg := string(data[:4])

    switch msg {
    case "MOVE":
        if len(data) < 10 {
            p.Conn.WriteInvl()
            break
        }

        from := string(data[5:7])
        to := string(data[8:10])

        if p.Color == m.WhoseTurn && m.Board.MakeMove(p.Color, from, to) {
            m.NextTurn()
            p.Conn.WriteString("ACPT")
            for _, opp := range m.Players {
                if opp != *p {
                    opp.Conn.Write(data[:10])
                }
            }
        } else {
            p.Conn.WriteString("RJCT")
        }
    default:
        p.Conn.WriteInvl()
    }
}

func HandleMatch(connWhite *conn.Conn, connBlack *conn.Conn) {
    log.Println("Match: started", connWhite.WsConn.RemoteAddr(), connBlack.WsConn.RemoteAddr())
    defer log.Println("Match: ended", connWhite.WsConn.RemoteAddr(), connBlack.WsConn.RemoteAddr())

    // Signal that that match has started, stop discarding messages
    connWhite.DoneCh <- struct{}{}
    connBlack.DoneCh <- struct{}{}

    var m Match
    m.Board = chess.NewBoard()
    m.WhoseTurn = chess.WHITE
    m.Players = []Player{
        {connWhite, chess.WHITE},
        {connBlack, chess.BLACK},
    }

    m.SendColors()

    for {
        select {
        case data, ok := <-connWhite.OutCh:
            if !ok {
                m.SignalEnd(m.Players[0].Color.String() + " disconnected")
                return
            }
            m.HandleMessage(&m.Players[0], data)
        case data, ok := <-connBlack.OutCh:
            if !ok {
                m.SignalEnd(m.Players[1].Color.String() + " disconnected")
                return
            }
            m.HandleMessage(&m.Players[1], data)
        case reason := <-m.endSignal:
            m.End(reason)
        }
    }
}
