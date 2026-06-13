package match

import (
	"log"
	"slices"
	"time"

	"github.com/aringq10/TCP/server/chess/classical"
	"github.com/aringq10/TCP/server/conn"
	"github.com/aringq10/TCP/server/msg"
)

const MAX_INVL = 5
const MATCH_DURATION = 10 * time.Second
const White = classical.White
const Black = classical.Black

type Match struct {
    Board *classical.Board
    Players map[classical.Color]*Player
    StartOfMatch time.Time
    MatchDuration time.Duration
    ended bool
}

func NewMatch(
    board *classical.Board,
    players map[classical.Color]*Player,
    startOfMatch time.Time,
    matchDuration time.Duration,
) *Match {
    return &Match{
        Board: board,
        Players: players,
        StartOfMatch: startOfMatch,
        MatchDuration: matchDuration,
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

func (m *Match) RotateTimeRemaining(p *Player) {
    totalTimeRemaining := -1 * p.TimeRemaining
    nextPlayer := m.Players[m.Board.WhoseTurn()]
    for _, p := range m.Players {
        if p == nextPlayer {
            p.Timer.Reset(nextPlayer.TimeRemaining)
        } else {
            p.Timer.Reset(2 * nextPlayer.TimeRemaining)
        }
        totalTimeRemaining += p.TimeRemaining
    }

    timeElapsed := time.Since(m.StartOfMatch)
    totalTime := time.Duration(len(m.Players)) * m.MatchDuration
    newTimeRemaining := totalTime - timeElapsed - totalTimeRemaining
    p.TimeRemaining = newTimeRemaining
}

func (m *Match) HandleMessage(p *Player, data []byte) {
    if p.SbsqINVL > MAX_INVL {
        m.End("too many invalid messages from " + p.Color.String())
        return
    }

    l := len(data)

    if l < 4 {
        p.WriteINVL()
        return
    }

    message := string(data[:4])

    switch message {
    case msg.Resign:
        if l != 4 {
            p.WriteINVL()
            break
        }

        winner := classical.OppositeColor(p.Color)
        m.End(Resignation(winner))
    case msg.Move:
        if l != 12 {
            p.WriteINVL()
            break
        }

        if !m.Board.MakeMove(p.Color, string(data[5:12])) {
            p.WriteRJCT()
            return
        }

        m.RotateTimeRemaining(p)

        p.WriteACPT()
        m.Broadcast(data[:12], p)

        log.Print(m.Board.String(classical.White))

        if m.Board.IsOver() {
            m.End(m.Board.Outcome().String())
        }
    default:
        p.WriteINVL()
    }
}

func HandleMatch(connWhite *conn.Conn, connBlack *conn.Conn) {
    whAddr := connWhite.WsConn.RemoteAddr()
    blAddr := connBlack.WsConn.RemoteAddr()
    log.Println("MATCH STARTED", whAddr, blAddr)
    defer log.Println("MATCH ENDED", whAddr, blAddr)

    connWhite.SignalMatchStart()
    connBlack.SignalMatchStart()

    players := map[classical.Color]*Player{
        White: NewPlayer(connWhite, White, MATCH_DURATION),
        Black: NewPlayer(connBlack, Black, MATCH_DURATION),
    }
    m := NewMatch(classical.NewBoard(), players, time.Now(), MATCH_DURATION)

    m.SendColors()
    defer m.End("unexpected end of match")

    log.Print(m.Board.String(classical.White))

    // Set to a longer duration than White's, so Timer is not nil
    // and White would timeout faster
    players[Black].Timer = time.NewTimer(2 * MATCH_DURATION)
    players[White].Timer = time.NewTimer(MATCH_DURATION)

    for {
        select {
        case data, ok := <-players[White].Conn.OutCh:
            if !ok {
                m.End(Abandonment(Black))
                return
            }
            m.HandleMessage(players[White], data)
        case data, ok := <-players[Black].Conn.OutCh:
            if !ok {
                m.End(Abandonment(White))
                return
            }
            m.HandleMessage(players[Black], data)
        case <-players[White].Timer.C:
            m.End(TimeForfeit(Black))
            return
        case <-players[Black].Timer.C:
            m.End(TimeForfeit(White))
            return
        }
    }
}

func Resignation(winner classical.Color) string {
    return classical.NewOutcome(classical.Resignation, winner).String()
}

func TimeForfeit(winner classical.Color) string {
    return classical.NewOutcome(classical.TimeForfeit, winner).String()
}

func Abandonment(winner classical.Color) string {
    return classical.NewOutcome(classical.Abandonment, winner).String()
}
