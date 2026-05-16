package match

import (
    "log"

    "github.com/aringq10/TCP/server/conn"
    "github.com/aringq10/TCP/server/chess"
)

const MAX_INVL = 5

type Player struct {
    Conn *conn.Conn
    Color chess.Color
}

type Match struct {
    Board *chess.Board
    Players []Player
}

func (m *Match) SendColors() {
    for _, p := range m.Players {
        p.Conn.WriteString(p.Color.String())
    }
}

func (m *Match) End() {
    for _, p := range m.Players {
        p.Conn.Close("End Of Match")
    }
}

func (m *Match) HandleMessage(p *Player, data []byte) {
    if p.Conn.SbsqInvl > MAX_INVL {
        m.End()
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

        if m.Board.MakeMove(p.Color, from, to) {
            p.Conn.WriteString("ACPT")
            for _, v := range m.Players {
                if v != *p {
                    v.Conn.Write(data[:10])
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
    conn.Conns.RemoveConn(connWhite)
    conn.Conns.RemoveConn(connBlack)

    log.Println("Match started between", connWhite.WsConn.RemoteAddr(), "with", connBlack.WsConn.RemoteAddr())
    defer log.Println("Match ended between", connWhite.WsConn.RemoteAddr(), "with", connBlack.WsConn.RemoteAddr())

    var m Match

    m.Board = chess.NewBoard()

    m.Players = []Player{
        {connWhite, chess.WHITE},
        {connBlack, chess.BLACK},
    }

    defer m.End()

    m.SendColors()

    chWhite := make(chan []byte)
    chBlack := make(chan []byte)

    go connWhite.ReadToChan(chWhite)
    go connBlack.ReadToChan(chBlack)

    var data []byte
    var ok bool

    for {
        select {
        case data, ok = <-chWhite:
            if !ok {
                return
            }
            m.HandleMessage(&m.Players[0], data)
        case data, ok = <-chBlack:
            if !ok {
                return
            }
            m.HandleMessage(&m.Players[1], data)
        }
    }
}
