package match

import (
	"fmt"
	"log"
	"time"

	"github.com/aringq10/TCP/server/chess/classical"
	"github.com/aringq10/TCP/server/conn"
	"github.com/aringq10/TCP/server/msg"
)

// Disqualify a player after this many invalid messages in a row.
const MAX_INVL = 6
// Set to a negative value for no time limit
const MATCH_DURATION = 120 * time.Second
const White = classical.White
const Black = classical.Black

type Match struct {
    board *classical.Board
    white *Player
    black *Player
    startOfMatch time.Time
    matchDuration time.Duration
    ended bool
}

func NewMatch(
    board *classical.Board,
    white *Player,
    black *Player,
    startOfMatch time.Time,
    matchDuration time.Duration,
) *Match {
    return &Match{
        board: board,
        white: white,
        black: black,
        startOfMatch: startOfMatch,
        matchDuration: matchDuration,
    }
}

func (m *Match) sendColors() {
    m.white.Write(msg.MatchColor[White] + m.getTimers())
    m.black.Write(msg.MatchColor[Black] + m.getTimers())
}

func (m *Match) sendInvalid(p *Player) {
    p.SbsqINVL++
    p.Write(msg.Invalid)
    if p.SbsqINVL >= MAX_INVL {
        m.sendEndOfMatch(disqualification(classical.OppositeColor(p.Color)))
    }
}

func (m *Match) sendReject(p *Player) {
    p.SbsqINVL = 0
    p.Write(msg.Reject + m.getTimers())
}

func (m *Match) sendAccept(p *Player) {
    p.SbsqINVL = 0
    p.Write(msg.Accept + m.getTimers())
}

func (m *Match) sendMove(p *Player, moveString string) {
    p.Write(msg.Move + " " + moveString + m.getTimers())
}

func (m *Match) sendEndOfMatch(reason string) {
    if m.ended {
        return
    }
    m.ended = true
    m.white.Write(msg.EndOfMatch + " " + reason)
    m.black.Write(msg.EndOfMatch + " " + reason)
    m.white.Conn.Close()
    m.black.Conn.Close()
}

func (m *Match) noTimeLimit() bool {
    return m.matchDuration < 0
}

func (m *Match) timeRemainingFor(c classical.Color) time.Duration {
    if m.noTimeLimit() {
        return 0.0
    }

    me, opp := m.getPlayers(c)
    if c != m.board.WhoseTurn() {
        return me.TimeRemaining
    }

    timeElapsed := time.Since(m.startOfMatch)
    totalTime := 2 * m.matchDuration
    return totalTime - timeElapsed - opp.TimeRemaining
}

func (m *Match) UpdateTimers() {
    if m.noTimeLimit() {
        return
    }

    // currPlayer is whose turn just started.
    // We need to update the previous player's remaining time.
    currPlayer, prevPlayer := m.getPlayers(m.board.WhoseTurn())
    prevPlayer.Timer.Stop()
    currPlayer.Timer.Reset(currPlayer.TimeRemaining)

    timeElapsed := time.Since(m.startOfMatch)
    totalTime := 2 * m.matchDuration
    prevPlayer.TimeRemaining = totalTime - timeElapsed - currPlayer.TimeRemaining
}

func (m *Match) HandleMessage(p *Player, data []byte) {
    l := len(data)

    if l < 4 {
        m.sendInvalid(p)
        return
    }

    message := string(data[:4])

    switch message {
    case msg.Resign:
        if l != 4 {
            m.sendInvalid(p)
            return
        }

        winner := classical.OppositeColor(p.Color)
        m.sendEndOfMatch(resignation(winner))
    case msg.Move:
        if l != 12 {
            m.sendInvalid(p)
            return
        }

        moveString := string(data[5:12])
        move, okMove := classical.ParseMove(p.Color, moveString)
        if !okMove {
            m.sendInvalid(p)
            return
        }

        if !m.board.MakeMove(move) {
            m.sendReject(p)
            return
        }

        opp := m.getPlayer(classical.OppositeColor(p.Color))
        m.UpdateTimers()
        m.sendAccept(p)
        m.sendMove(opp, moveString)

        log.Print(m.board.String(classical.White))

        if m.board.IsOver() {
            m.sendEndOfMatch(m.board.Outcome().String())
        }
    default:
        m.sendInvalid(p)
    }
}

func HandleMatch(connWhite *conn.Conn, connBlack *conn.Conn) {
    whAddr := connWhite.WsConn.RemoteAddr()
    blAddr := connBlack.WsConn.RemoteAddr()
    log.Println("MATCH STARTED", whAddr, blAddr)
    defer log.Println("MATCH ENDED", whAddr, blAddr)

    pw := NewPlayer(connWhite, White, MATCH_DURATION)
    pb := NewPlayer(connBlack, Black, MATCH_DURATION)
    m := NewMatch(classical.NewBoard(), pw, pb, time.Now(), MATCH_DURATION)

    m.white.Conn.SignalMatchStart()
    m.black.Conn.SignalMatchStart()

    m.sendColors()

    // Immediately stop black's timer so it doesn't fire before white's.
    // We need to create it here since select can't take a nil Timer.
    m.black.Timer = time.NewTimer(MATCH_DURATION)
    m.black.Timer.Stop()
    m.white.Timer = time.NewTimer(MATCH_DURATION)
    if m.noTimeLimit() {
        m.white.Timer.Stop()
    }

    log.Print(m.board.String(classical.White))

    for {
        select {
        case data, ok := <-m.white.Conn.OutCh:
            if !ok {
                m.sendEndOfMatch(abandonment(Black))
                return
            }
            m.HandleMessage(m.white, data)
        case data, ok := <-m.black.Conn.OutCh:
            if !ok {
                m.sendEndOfMatch(abandonment(White))
                return
            }
            m.HandleMessage(m.black, data)
        case <-m.white.Timer.C:
            m.sendEndOfMatch(timeForfeit(Black))
            return
        case <-m.black.Timer.C:
            m.sendEndOfMatch(timeForfeit(White))
            return
        }
    }
}

func resignation(winner classical.Color) string {
    return classical.NewOutcome(classical.Resignation, winner).String()
}

func timeForfeit(winner classical.Color) string {
    return classical.NewOutcome(classical.TimeForfeit, winner).String()
}

func abandonment(winner classical.Color) string {
    return classical.NewOutcome(classical.Abandonment, winner).String()
}

func disqualification(winner classical.Color) string {
    return classical.NewOutcome(classical.Disqualification, winner).String()
}

func (m *Match) getPlayers(yourColor classical.Color) (you *Player, opp *Player) {
    if yourColor == White {
        return m.white, m.black
    } else {
        return m.black, m.white
    }
}

func (m *Match) getPlayer(c classical.Color) *Player {
    if c == White {
        return m.white
    } else {
        return m.black
    }
}

func (m *Match) getTimers() string {
    return fmt.Sprintf(" %.3f %.3f", m.timeRemainingFor(White).Seconds(), m.timeRemainingFor(Black).Seconds())
}
