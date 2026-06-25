package classical

import (
    "testing"
)

func TestQueenIsValidMove(t *testing.T) {
    cases := []struct {
        name  string
        place map[string]Piece
        color Color
        from  string
        to    string
        want  bool
    }{
        {
            name:  "slide along file",
            place: map[string]Piece{"d1": NewQueen(White)},
            color: White, from: "d1", to: "d5",
            want: true,
        },
        {
            name:  "slide along rank",
            place: map[string]Piece{"d1": NewQueen(White)},
            color: White, from: "d1", to: "h1",
            want: true,
        },
        {
            name:  "diagonal up-right",
            place: map[string]Piece{"d1": NewQueen(White)},
            color: White, from: "d1", to: "h5",
            want: true,
        },
        {
            name:  "diagonal up-left",
            place: map[string]Piece{"d1": NewQueen(White)},
            color: White, from: "d1", to: "a4",
            want: true,
        },
        {
            name:  "knight-shape",
            place: map[string]Piece{"d1": NewQueen(White)},
            color: White, from: "d1", to: "e3",
            want: false,
        },
        {
            name:  "stays in place",
            place: map[string]Piece{"d1": NewQueen(White)},
            color: White, from: "d1", to: "d1",
            want: false,
        },
        {
            name:  "blocked by own piece on file",
            place: map[string]Piece{"d1": NewQueen(White), "d3": NewPawn(White)},
            color: White, from: "d1", to: "d5",
            want: false,
        },
        {
            name:  "blocked by enemy piece on file",
            place: map[string]Piece{"d1": NewQueen(White), "d3": NewPawn(Black)},
            color: White, from: "d1", to: "d5",
            want: false,
        },
        {
            name:  "blocked by own piece on diagonal",
            place: map[string]Piece{"d1": NewQueen(White), "f3": NewPawn(White)},
            color: White, from: "d1", to: "h5",
            want: false,
        },
        {
            name:  "blocked by enemy piece on diagonal",
            place: map[string]Piece{"d1": NewQueen(White), "f3": NewPawn(Black)},
            color: White, from: "d1", to: "h5",
            want: false,
        },
        {
            name:  "capture enemy at destination along file",
            place: map[string]Piece{"d1": NewQueen(White), "d5": NewPawn(Black)},
            color: White, from: "d1", to: "d5",
            want: true,
        },
        {
            name:  "capture enemy at destination diagonal",
            place: map[string]Piece{"d1": NewQueen(White), "h5": NewPawn(Black)},
            color: White, from: "d1", to: "h5",
            want: true,
        },
        {
            name:  "own piece at destination",
            place: map[string]Piece{"d1": NewQueen(White), "d5": NewPawn(White)},
            color: White, from: "d1", to: "d5",
            want: false,
        },
        {
            name:  "slide backwards along file",
            place: map[string]Piece{"d5": NewQueen(Black)},
            color: Black, from: "d5", to: "d1",
            want: true,
        },
        {
            name:  "adjacent diagonal",
            place: map[string]Piece{"d1": NewQueen(White)},
            color: White, from: "d1", to: "e2",
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
            queen, ok := b.squares[from.Row][from.Col].(*Queen)
            if !ok {
                t.Fatalf("no queen at %s", tc.from)
            }

            _, got := queen.IsValidMove(&b, from, sq(tc.to))
            if got != tc.want {
                t.Errorf("IsValidMove() = %v, want %v", got, tc.want)
            }
        })
    }
}

func TestQueenExecuteMovesPieceDiagonal(t *testing.T) {
    var b Board
    queen := NewQueen(White)
    b.squares[0][3] = queen

    exec, valid := queen.IsValidMove(&b, sq("d1"), sq("h5"))
    if !valid {
        t.Fatal("move flagged invalid")
    }
    exec()

    if b.squares[0][3] != nil {
        t.Errorf("origin not cleared")
    }
    if b.squares[4][7] != queen {
        t.Errorf("queen not at destination")
    }
}

func TestQueenExecuteCapturesEnemy(t *testing.T) {
    var b Board
    queen := NewQueen(White)
    b.squares[0][3] = queen
    b.squares[4][3] = NewPawn(Black)

    exec, valid := queen.IsValidMove(&b, sq("d1"), sq("d5"))
    if !valid {
        t.Fatal("capture flagged invalid")
    }
    exec()

    if b.squares[4][3] != queen {
        t.Errorf("queen did not replace captured piece")
    }
    if b.squares[0][3] != nil {
        t.Errorf("origin not cleared")
    }
}
