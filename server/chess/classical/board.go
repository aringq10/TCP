package classical

import (
    "github.com/aringq10/TCP/server/chess"
)

func NewBoard() *chess.Board {
    var b chess.Board
    for col := range 8 {
        b[1][col] = NewPawn(chess.WHITE)
        b[6][col] = NewPawn(chess.BLACK)
    }
    return &b
}
