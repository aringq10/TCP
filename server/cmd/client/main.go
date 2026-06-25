package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "sync"

    "github.com/aringq10/TCP/server/chess/classical"
    "github.com/aringq10/TCP/server/msg"
    "github.com/gorilla/websocket"
)

type state struct {
    mu          sync.Mutex
    myColor     classical.Color
    pendingMove string
    board       *classical.Board
}

func (s *state) setColor(c classical.Color) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.myColor = c
}

func (s *state) color() classical.Color {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.myColor
}

func (s *state) setPending(m string) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    if s.pendingMove != "" {
        return false
    }
    s.pendingMove = m
    return true
}

func (s *state) clearPending() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.pendingMove = ""
}

func (s *state) printBoard() {
    fmt.Println(s.board.String(s.myColor))
}

func main() {
    if len(os.Args) != 3 {
        fmt.Fprintf(os.Stderr, "usage: %s <host> <port>\n", os.Args[0])
        os.Exit(2)
    }

    addr := "ws://" + os.Args[1] + ":" + os.Args[2] + "/ws"
    conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
    if err != nil {
        log.Fatalf("could not connect to %s: %v", addr, err)
    }
    defer conn.Close()

    fmt.Printf("connected to %s, waiting for match...\n", addr)

    s := &state{board: classical.NewBoard()}
    done := make(chan struct{})

    go readLoop(conn, s, done)
    go writeLoop(conn, s)

    <-done
    fmt.Println("disconnected.")
}

func readLoop(conn *websocket.Conn, s *state, done chan struct{}) {
    defer close(done)
    for {
        mt, payload, err := conn.ReadMessage()
        if err != nil {
            if ce, ok := err.(*websocket.CloseError); ok {
                reason := ce.Text
                if strings.HasPrefix(reason, "EOM: ") {
                    fmt.Printf("\nmatch ended: %s\n", reason[5:])
                } else {
                    fmt.Printf("\nconnection closed: %s\n", reason)
                }
                return
            }
            log.Printf("read error: %v", err)
            return
        }
        if mt != websocket.TextMessage || len(payload) < 4 {
            fmt.Println("received unknown message")
            continue
        }
        handleMessage(payload, s)
    }
}

func handleMessage(payload []byte, s *state) {
    tag := string(payload[:4])
    switch tag {
    case msg.Black:
        s.setColor(classical.Black)
        fmt.Println("you are playing as Black")
        s.printBoard()
    case msg.White:
        s.setColor(classical.White)
        fmt.Println("you are playing as White")
        s.printBoard()
    case msg.Accept:
        m, _ := classical.ParseMove(s.color(), s.pendingMove)
        s.board.MakeMove(m)
        s.clearPending()
        fmt.Println("move accepted")
        s.printBoard()
    case msg.Reject:
        s.clearPending()
        fmt.Println("move rejected, try again")
    case msg.Move:
        if len(payload) < 10 {
            fmt.Println("received malformed MOVE message")
            return
        }
        moveStr := string(payload[5:10])
        m, _ := classical.ParseMove(classical.OppositeColor(s.color()), moveStr)
        s.board.MakeMove(m)
        fmt.Printf("opponent played: %s\n", moveStr)
        s.printBoard()
    case msg.Invalid:
        s.clearPending()
        fmt.Println("invalid move, use format e2 e4 -")
    default:
        fmt.Printf("received unknown message: %q\n", payload)
    }
    prompt()
}

func writeLoop(conn *websocket.Conn, s *state) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        input := strings.TrimSpace(scanner.Text())
        if !s.setPending(input) {
            fmt.Printf("already waiting on move %q\n", s.pendingMove)
            prompt()
            continue
        }
        msg := "MOVE " + input
        if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
            log.Printf("write error: %v", err)
            return
        }
    }
}

func prompt() {
    fmt.Print("> ")
}
