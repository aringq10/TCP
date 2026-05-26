package classical

import (
    "testing"
)

func TestKnightIsValidMove(t *testing.T) {
    cases := []struct {
        name  string
        place map[string]Piece
        color Color
        from  string
        to    string
        want  bool
    }{
        {
            name:  "L up-right tall",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "e6",
            want: true,
        },
        {
            name:  "L up-right wide",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "f5",
            want: true,
        },
        {
            name:  "L up-left tall",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "c6",
            want: true,
        },
        {
            name:  "L up-left wide",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "b5",
            want: true,
        },
        {
            name:  "L down-right tall",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "e2",
            want: true,
        },
        {
            name:  "L down-right wide",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "f3",
            want: true,
        },
        {
            name:  "L down-left tall",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "c2",
            want: true,
        },
        {
            name:  "L down-left wide",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "b3",
            want: true,
        },
        {
            name:  "straight along file",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "d6",
            want: false,
        },
        {
            name:  "straight along rank",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "g4",
            want: false,
        },
        {
            name:  "diagonal",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "f6",
            want: false,
        },
        {
            name:  "adjacent square",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "d5",
            want: false,
        },
        {
            name:  "stays in place",
            place: map[string]Piece{"d4": NewKnight(White)},
            color: White, from: "d4", to: "d4",
            want: false,
        },
        {
            name:  "jumps over piece on path",
            place: map[string]Piece{"d4": NewKnight(White), "d5": NewPawn(White), "e5": NewPawn(White)},
            color: White, from: "d4", to: "e6",
            want: true,
        },
        {
            name:  "capture enemy at destination",
            place: map[string]Piece{"d4": NewKnight(White), "e6": NewPawn(Black)},
            color: White, from: "d4", to: "e6",
            want: true,
        },
        {
            name:  "own piece at destination",
            place: map[string]Piece{"d4": NewKnight(White), "e6": NewPawn(White)},
            color: White, from: "d4", to: "e6",
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
            knight, ok := b.squares[from.Row][from.Col].(*Knight)
            if !ok {
                t.Fatalf("no knight at %s", tc.from)
            }

            _, got := knight.IsValidMove(&b, Move{
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

func TestKnightExecuteMovesPiece(t *testing.T) {
    var b Board
    knight := NewKnight(White)
    b.squares[3][3] = knight

    exec, valid := knight.IsValidMove(&b, Move{
        PlayerColor: White,
        From:        sq("d4"),
        To:          sq("e6"),
    })
    if !valid {
        t.Fatal("move flagged invalid")
    }
    exec()

    if b.squares[3][3] != nil {
        t.Errorf("origin not cleared")
    }
    if b.squares[5][4] != knight {
        t.Errorf("knight not at destination")
    }
}

func TestKnightExecuteCapturesEnemy(t *testing.T) {
    var b Board
    knight := NewKnight(White)
    b.squares[3][3] = knight
    b.squares[5][4] = NewPawn(Black)

    exec, valid := knight.IsValidMove(&b, Move{
        PlayerColor: White,
        From:        sq("d4"),
        To:          sq("e6"),
    })
    if !valid {
        t.Fatal("capture flagged invalid")
    }
    exec()

    if b.squares[5][4] != knight {
        t.Errorf("knight did not replace captured piece")
    }
    if b.squares[3][3] != nil {
        t.Errorf("origin not cleared")
    }
}
