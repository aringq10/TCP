package classical

import (
    "testing"
)

func TestBishopIsValidMove(t *testing.T) {
    cases := []struct {
        name  string
        place map[string]Piece
        color Color
        from  string
        to    string
        want  bool
    }{
        {
            name:  "diagonal up-right",
            place: map[string]Piece{"c1": NewBishop(White)},
            color: White, from: "c1", to: "h6",
            want: true,
        },
        {
            name:  "diagonal up-left",
            place: map[string]Piece{"f1": NewBishop(White)},
            color: White, from: "f1", to: "a6",
            want: true,
        },
        {
            name:  "diagonal down-right",
            place: map[string]Piece{"a8": NewBishop(Black)},
            color: Black, from: "a8", to: "h1",
            want: true,
        },
        {
            name:  "diagonal down-left",
            place: map[string]Piece{"h8": NewBishop(Black)},
            color: Black, from: "h8", to: "a1",
            want: true,
        },
        {
            name:  "straight along file",
            place: map[string]Piece{"c1": NewBishop(White)},
            color: White, from: "c1", to: "c5",
            want: false,
        },
        {
            name:  "straight along rank",
            place: map[string]Piece{"c1": NewBishop(White)},
            color: White, from: "c1", to: "h1",
            want: false,
        },
        {
            name:  "knight-shape",
            place: map[string]Piece{"c1": NewBishop(White)},
            color: White, from: "c1", to: "d3",
            want: false,
        },
        {
            name:  "stays in place",
            place: map[string]Piece{"c1": NewBishop(White)},
            color: White, from: "c1", to: "c1",
            want: false,
        },
        {
            name:  "blocked by own piece on path",
            place: map[string]Piece{"c1": NewBishop(White), "e3": NewPawn(White)},
            color: White, from: "c1", to: "h6",
            want: false,
        },
        {
            name:  "blocked by enemy piece on path",
            place: map[string]Piece{"c1": NewBishop(White), "e3": NewPawn(Black)},
            color: White, from: "c1", to: "h6",
            want: false,
        },
        {
            name:  "capture enemy at destination",
            place: map[string]Piece{"c1": NewBishop(White), "h6": NewPawn(Black)},
            color: White, from: "c1", to: "h6",
            want: true,
        },
        {
            name:  "own piece at destination",
            place: map[string]Piece{"c1": NewBishop(White), "h6": NewPawn(White)},
            color: White, from: "c1", to: "h6",
            want: false,
        },
        {
            name:  "adjacent diagonal",
            place: map[string]Piece{"c1": NewBishop(White)},
            color: White, from: "c1", to: "d2",
            want: true,
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
            bishop, ok := b.squares[from.Row][from.Col].(*Bishop)
            if !ok {
                t.Fatalf("no bishop at %s", tc.from)
            }

            _, got := bishop.IsValidMove(&b, from, sq(tc.to))
            if got != tc.want {
                t.Errorf("IsValidMove() = %v, want %v", got, tc.want)
            }
        })
    }
}

func TestBishopExecuteMovesPiece(t *testing.T) {
    var b Board
    bishop := NewBishop(White)
    b.squares[0][2] = bishop

    exec, valid := bishop.IsValidMove(&b, sq("c1"), sq("h6"))
    if !valid {
        t.Fatal("move flagged invalid")
    }
    exec()

    if b.squares[0][2] != nil {
        t.Errorf("origin not cleared")
    }
    if b.squares[5][7] != bishop {
        t.Errorf("bishop not at destination")
    }
}

func TestBishopExecuteCapturesEnemy(t *testing.T) {
    var b Board
    bishop := NewBishop(White)
    b.squares[0][2] = bishop
    b.squares[5][7] = NewPawn(Black)

    exec, valid := bishop.IsValidMove(&b, sq("c1"), sq("h6"))
    if !valid {
        t.Fatal("capture flagged invalid")
    }
    exec()

    if b.squares[5][7] != bishop {
        t.Errorf("bishop did not replace captured piece")
    }
    if b.squares[0][2] != nil {
        t.Errorf("origin not cleared")
    }
}
