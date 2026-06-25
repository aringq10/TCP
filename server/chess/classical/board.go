package classical

import (
	"fmt"
	"slices"
	"strings"
)

type Board struct {
    // First index is row: 0 = rank 1, 7 = rank 8.
    // Second index is column: 0 = file a, 7 = file h.
    squares [8][8]Piece
    lastMove Move
    whoseTurn Color
    outcome Outcome
}

func NewBoard() *Board {
    var b Board
    s := &b.squares
    placePawns(s)
    s[0][0], s[0][7] = NewRook(White), NewRook(White)
    s[0][1], s[0][6] = NewKnight(White), NewKnight(White)
    s[0][2], s[0][5] = NewBishop(White), NewBishop(White)
    s[0][3] = NewQueen(White)
    s[0][4] = NewKing(White)
    s[7][0], s[7][7] = NewRook(Black), NewRook(Black)
    s[7][1], s[7][6] = NewKnight(Black), NewKnight(Black)
    s[7][2], s[7][5] = NewBishop(Black), NewBishop(Black)
    s[7][3] = NewQueen(Black)
    s[7][4] = NewKing(Black)
    return &b
}

func placePawns(s *[8][8]Piece) {
    for col := range 8 {
        s[1][col] = NewPawn(White)
        s[6][col] = NewPawn(Black)
    }
}

func (b *Board) IsOver() bool {
    return b.outcome.Result.IsOver()
}

func (b *Board) Outcome() Outcome {
    return b.outcome
}

func (b *Board) SetOutcome(o Outcome) {
    b.outcome = o
}

func (b *Board) WhoseTurn() Color {
    return b.whoseTurn
}

func (b *Board) String(color Color) string {
    var sb strings.Builder
    sb.WriteByte('\n')

    rows := make([]int, 8)
    cols := make([]int, 8)
    for i := range 8 {
        rows[i] = i
        cols[i] = i
    }
    if color == White {
        slices.Reverse(rows)
    } else {
        slices.Reverse(cols)
    }

    letters := "    a b c d e f g h"
    if color == Black {
        letters = "    h g f e d c b a"
    }

    sb.WriteString("  ┌─────────────────┐\n")
    for _, row := range rows {
        sb.WriteString(fmt.Sprintf("%d │ ", row+1))
        for _, col := range cols {
            if p := b.squares[row][col]; p != nil {
                sb.WriteString(p.String())
            } else {
                sb.WriteByte('*')
            }
            sb.WriteByte(' ')
        }
        sb.WriteString("│\n")
    }
    sb.WriteString("  └─────────────────┘\n")
    sb.WriteString(letters + "\n")
    return sb.String()
}

func (b *Board) MakeMove(m Move) bool {
    if b.IsOver() {
        return false
    }

    if m.Color != b.whoseTurn {
        return false
    }

    piece := b.squares[m.From.Row][m.From.Col]
    if piece == nil || piece.Color() != m.Color {
        return false
    }

    exec, valid := piece.IsValidMove(b, m.From, m.To)
    if !valid || exec == nil {
        return false
    }

    _, isPawn := piece.(*Pawn)
    backRank := 7
    if m.Color == Black {
        backRank = 0
    }
    landsLast := isPawn && m.To.Row == backRank

    // A move is a promotion iff a pawn reaches the back rank. Reject a promotion
    // flag on any other move, and a pawn reaching the back rank without one.
    if landsLast != m.IsPromotion() {
        return false
    }

    if landsLast {
        base := exec
        exec = func() {
            base()
            b.squares[m.To.Row][m.To.Col] = getPromotionPiece(m.Promotion, m.Color)
        }
    }

    snapshot := b.squares

    exec()
    if b.isChecked(piece.Color()) {
        b.squares = snapshot
        return false
    }

    piece.IncMoveCount()
    b.lastMove = m
    b.rotateTurn()

    // Did this move end the game for the side now to move?
    b.outcome = NewOutcome(b.terminationFor(b.whoseTurn), piece.Color())

    return true
}

func getPromotionPiece(p Promotion, c Color) Piece {
    switch p {
    case ToQueen:
        return NewQueen(c)
    case ToRook:
        return NewRook(c)
    case ToBishop:
        return NewBishop(c)
    case ToKnight:
        return NewKnight(c)
    default:
        return nil
    }
}

func (b *Board) rotateTurn() {
    b.whoseTurn = OppositeColor(b.whoseTurn)
}

func (b *Board) isAttacked(s Square, by Color) bool {
    attacked := b.squares[s.Row][s.Col]
    if attacked != nil && attacked.Color() == by {
        return false
    }

    for row := range 8 {
        for col := range 8 {
            attacker := b.squares[row][col]
            if attacker == nil || attacker.Color() != by {
                continue
            }

            from := Square{Row: row, Col: col}
            if _, valid := attacker.IsValidMove(b, from, s); valid {
                return true
            }
        }
    }

    return false
}

func (b *Board) tryFindKing(c Color) (s Square, ok bool) {
    for row := range 8 {
        for col := range 8 {
            piece := b.squares[row][col]
            if piece == nil {
                continue
            }
            if _, isKing := piece.(*King); piece.Color() == c && isKing {
                return Square{Row: row, Col: col}, true
            }
        }
    }

    return Square{}, false
}

func (b *Board) isChecked(c Color) bool {
    s, ok := b.tryFindKing(c)
    if !ok {
        return false
    }

    return b.isAttacked(s, OppositeColor(c))
}

// Doesn't generate pawn promotions, as it's only used for check-escape detection.
func (b *Board) getValidMoves(from Square) []func() {
    piece := b.squares[from.Row][from.Col]
    if piece == nil {
        return nil
    }

    moves := make([]func(), 0, 64)

    for row := range 8 {
        for col := range 8 {
            to := Square{Row: row, Col: col}
            if move, valid := piece.IsValidMove(b, from, to); valid {
                moves = append(moves, move)
            }
        }
    }

    return moves
}

func (b *Board) terminationFor(c Color) Termination {
    _, ok := b.tryFindKing(c)
    if !ok {
        return NoTermination
    }

    friendlySquares := make([]Square, 0, 64)

    // Collect squares containing friendly pieces
    for row := range 8 {
        for col := range 8 {
            piece := b.squares[row][col]
            if piece == nil {
                continue
            }
            if piece.Color() == c {
                friendlySquares = append(friendlySquares, Square{Row: row, Col: col})
            }
        }
    }

    snapshot := b.squares

    for _, s := range friendlySquares {
        moves := b.getValidMoves(s)
        for _, m := range moves {
            m()
            legal := !b.isChecked(c)
            b.squares = snapshot
            if legal {
                return NoTermination
            }
        }
    }

    // No legal move: it's checkmate if the king is in check, otherwise stalemate.
    if b.isChecked(c) {
        return Checkmate
    } else {
        return Stalemate
    }
}
