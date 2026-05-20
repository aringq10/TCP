package classical

import (
    "testing"

    "github.com/aringq10/TCP/server/chess"
)

func TestRookIsValidMove(t *testing.T) {
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
            place: map[string]Piece{"a1": NewRook(White)},
            color: White, from: "a1", to: "a5",
            want: true,
        },
        {
            name:  "slide along rank",
            place: map[string]Piece{"a1": NewRook(White)},
            color: White, from: "a1", to: "h1",
            want: true,
        },
        {
            name:  "diagonal",
            place: map[string]Piece{"a1": NewRook(White)},
            color: White, from: "a1", to: "h8",
            want: false,
        },
        {
            name:  "knight-shape",
            place: map[string]Piece{"a1": NewRook(White)},
            color: White, from: "a1", to: "c2",
            want: false,
        },
        {
            name:  "stays in place",
            place: map[string]Piece{"a1": NewRook(White)},
            color: White, from: "a1", to: "a1",
            want: false,
        },
        {
            name:  "blocked by own piece on path",
            place: map[string]Piece{"a1": NewRook(White), "a3": NewPawn(White)},
            color: White, from: "a1", to: "a5",
            want: false,
        },
        {
            name:  "blocked by enemy piece on path",
            place: map[string]Piece{"a1": NewRook(White), "a3": NewPawn(Black)},
            color: White, from: "a1", to: "a5",
            want: false,
        },
        {
            name:  "capture enemy at destination",
            place: map[string]Piece{"a1": NewRook(White), "a5": NewPawn(Black)},
            color: White, from: "a1", to: "a5",
            want: true,
        },
        {
            name:  "own piece at destination",
            place: map[string]Piece{"a1": NewRook(White), "a5": NewPawn(White)},
            color: White, from: "a1", to: "a5",
            want: false,
        },
        {
            name:  "slide backwards along file",
            place: map[string]Piece{"a5": NewRook(White)},
            color: White, from: "a5", to: "a1",
            want: true,
        },
        {
            name:  "slide right along rank past empty squares",
            place: map[string]Piece{"d4": NewRook(Black)},
            color: Black, from: "d4", to: "g4",
            want: true,
        },
        {
            name:  "adjacent square along file",
            place: map[string]Piece{"a1": NewRook(White)},
            color: White, from: "a1", to: "a2",
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
            rook, ok := b.squares[from.Row][from.Col].(*Rook)
            if !ok {
                t.Fatalf("no rook at %s", tc.from)
            }

            _, got := rook.IsValidMove(&b, chess.Move{
                PlayerColor: tc.color,
                From:        from,
                To:          sq(tc.to),
            })
            if got != tc.want {
                t.Errorf("IsValidMove() = %v, want %v", got, tc.want)
            }
        })
    }
}

func TestRookExecuteMovesPiece(t *testing.T) {
    var b Board
    rook := NewRook(White)
    b.squares[0][0] = rook

    exec, valid := rook.IsValidMove(&b, chess.Move{
        PlayerColor: White,
        From:        sq("a1"),
        To:          sq("a5"),
    })
    if !valid {
        t.Fatal("move flagged invalid")
    }
    exec()

    if b.squares[0][0] != nil {
        t.Errorf("origin not cleared")
    }
    if b.squares[4][0] != rook {
        t.Errorf("rook not at destination")
    }
}

func TestRookExecuteCapturesEnemy(t *testing.T) {
    var b Board
    rook := NewRook(White)
    b.squares[0][0] = rook
    b.squares[0][4] = NewPawn(Black)

    exec, valid := rook.IsValidMove(&b, chess.Move{
        PlayerColor: White,
        From:        sq("a1"),
        To:          sq("e1"),
    })
    if !valid {
        t.Fatal("capture flagged invalid")
    }
    exec()

    if b.squares[0][4] != rook {
        t.Errorf("rook did not replace captured piece")
    }
    if b.squares[0][0] != nil {
        t.Errorf("origin not cleared")
    }
}
