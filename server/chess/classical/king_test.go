package classical

import (
    "testing"
)

func TestKingIsValidMove(t *testing.T) {
    cases := []struct {
        name  string
        place map[string]Piece
        color Color
        from  string
        to    string
        want  bool
    }{
        {
            name:  "up",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "e5",
            want: true,
        },
        {
            name:  "down",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "e3",
            want: true,
        },
        {
            name:  "left",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "d4",
            want: true,
        },
        {
            name:  "right",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "f4",
            want: true,
        },
        {
            name:  "up-right",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "f5",
            want: true,
        },
        {
            name:  "up-left",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "d5",
            want: true,
        },
        {
            name:  "down-right",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "f3",
            want: true,
        },
        {
            name:  "down-left",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "d3",
            want: true,
        },
        {
            name:  "two squares along file",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "e6",
            want: false,
        },
        {
            name:  "two squares along rank",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "g4",
            want: false,
        },
        {
            name:  "two squares diagonal",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "g6",
            want: false,
        },
        {
            name:  "knight jump",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "f6",
            want: false,
        },
        {
            name:  "stays in place",
            place: map[string]Piece{"e4": NewKing(White)},
            color: White, from: "e4", to: "e4",
            want: false,
        },
        {
            name:  "capture enemy at destination",
            place: map[string]Piece{"e4": NewKing(White), "f5": NewPawn(Black)},
            color: White, from: "e4", to: "f5",
            want: true,
        },
        {
            name:  "own piece at destination",
            place: map[string]Piece{"e4": NewKing(White), "f5": NewPawn(White)},
            color: White, from: "e4", to: "f5",
            want: false,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            var b Board
            for s, p := range tc.place {
                square := sq(s)
                b.squares[square.Row][square.Col] = p
            }

            from := sq(tc.from)
            king, ok := b.squares[from.Row][from.Col].(*King)
            if !ok {
                t.Fatalf("no king at %s", tc.from)
            }

            _, got := king.IsValidMove(&b, from, sq(tc.to))
            if got != tc.want {
                t.Errorf("IsValidMove() = %v, want %v", got, tc.want)
            }
        })
    }
}

func TestKingExecuteMovesPiece(t *testing.T) {
    var b Board
    king := NewKing(White)
    b.squares[3][4] = king

    exec, valid := king.IsValidMove(&b, sq("e4"), sq("e5"))
    if !valid {
        t.Fatal("move flagged invalid")
    }
    exec()

    if b.squares[3][4] != nil {
        t.Errorf("origin not cleared")
    }
    if b.squares[4][4] != king {
        t.Errorf("king not at destination")
    }
}

func TestKingExecuteCapturesEnemy(t *testing.T) {
    var b Board
    king := NewKing(White)
    b.squares[3][4] = king
    b.squares[4][5] = NewPawn(Black)

    exec, valid := king.IsValidMove(&b, sq("e4"), sq("f5"))
    if !valid {
        t.Fatal("capture flagged invalid")
    }
    exec()

    if b.squares[4][5] != king {
        t.Errorf("king did not replace captured piece")
    }
    if b.squares[3][4] != nil {
        t.Errorf("origin not cleared")
    }
}
