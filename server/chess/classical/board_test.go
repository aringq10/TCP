package classical

import (
    "testing"
)

func TestIsChecked(t *testing.T) {
    cases := []struct {
        name  string
        place map[string]Piece
        color Color
        want  bool
    }{
        {
            name:  "lone king not checked",
            place: map[string]Piece{"e1": NewKing(White)},
            color: White,
            want:  false,
        },
        {
            name:  "rook checks along file",
            place: map[string]Piece{"e1": NewKing(White), "e8": NewRook(Black)},
            color: White,
            want:  true,
        },
        {
            name:  "rook checks along rank",
            place: map[string]Piece{"a1": NewKing(White), "h1": NewRook(Black)},
            color: White,
            want:  true,
        },
        {
            name:  "rook check blocked by friendly piece",
            place: map[string]Piece{"e1": NewKing(White), "e4": NewPawn(White), "e8": NewRook(Black)},
            color: White,
            want:  false,
        },
        {
            name:  "rook check blocked by enemy piece",
            place: map[string]Piece{"e1": NewKing(White), "e4": NewPawn(Black), "e8": NewRook(Black)},
            color: White,
            want:  false,
        },
        {
            name:  "bishop checks along diagonal",
            place: map[string]Piece{"e1": NewKing(White), "a5": NewBishop(Black)},
            color: White,
            want:  true,
        },
        {
            name:  "knight checks",
            place: map[string]Piece{"e1": NewKing(White), "f3": NewKnight(Black)},
            color: White,
            want:  true,
        },
        {
            name:  "pawn checks diagonally",
            place: map[string]Piece{"e1": NewKing(White), "d2": NewPawn(Black)},
            color: White,
            want:  true,
        },
        {
            name:  "no king is not checked",
            place: map[string]Piece{"a1": NewRook(Black)},
            color: White,
            want:  false,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            var b Board
            for s, p := range tc.place {
                square := sq(s)
                b.squares[square.Row][square.Col] = p
            }

            if got := b.isChecked(tc.color); got != tc.want {
                t.Errorf("isChecked(%v) = %v, want %v", tc.color, got, tc.want)
            }
        })
    }
}

func TestTerminationFor(t *testing.T) {
    cases := []struct {
        name  string
        place map[string]Piece
        color Color
        want  Termination
    }{
        {
            name:  "two-rook back-rank checkmate",
            place: map[string]Piece{"a8": NewKing(Black), "h8": NewRook(White), "h7": NewRook(White)},
            color: Black,
            want:  Checkmate,
        },
        {
            name:  "king-and-queen corner stalemate",
            place: map[string]Piece{"a8": NewKing(Black), "b6": NewQueen(White)},
            color: Black,
            want:  Stalemate,
        },
        {
            name:  "in check but king can escape",
            place: map[string]Piece{"e8": NewKing(Black), "e1": NewRook(White)},
            color: Black,
            want:  NoTermination,
        },
        {
            name:  "no king yields no termination",
            place: map[string]Piece{"a1": NewRook(Black)},
            color: Black,
            want:  NoTermination,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            var b Board
            for s, p := range tc.place {
                square := sq(s)
                b.squares[square.Row][square.Col] = p
            }

            if got := b.terminationFor(tc.color); got != tc.want {
                t.Errorf("terminationFor(%v) = %v, want %v", tc.color, got, tc.want)
            }
        })
    }
}

func TestTerminationForStartingPosition(t *testing.T) {
    b := NewBoard()
    if got := b.terminationFor(White); got != NoTermination {
        t.Errorf("terminationFor(White) = %v, want %v", got, NoTermination)
    }
}
