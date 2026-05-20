package classical

import (
    "testing"

    "github.com/aringq10/TCP/server/chess"
)

func sq(s string) chess.Square {
    return chess.Square{Row: int(s[1] - '1'), Col: int(s[0] - 'a')}
}

func TestPawnIsValidMove(t *testing.T) {
    cases := []struct {
        name  string
        place map[string]Piece
        last  chess.Move
        color Color
        from  string
        to    string
        want  bool
    }{
        {
            name:  "white single push",
            place: map[string]Piece{"e2": NewPawn(White)},
            color: White, from: "e2", to: "e3",
            want: true,
        },
        {
            name:  "white single push blocked",
            place: map[string]Piece{"e2": NewPawn(White), "e3": NewPawn(Black)},
            color: White, from: "e2", to: "e3",
            want: false,
        },
        {
            name:  "white double push from start",
            place: map[string]Piece{"e2": NewPawn(White)},
            color: White, from: "e2", to: "e4",
            want: true,
        },
        {
            name:  "white double push not from start",
            place: map[string]Piece{"e3": NewPawn(White)},
            color: White, from: "e3", to: "e5",
            want: false,
        },
        {
            name:  "white double push blocked at intermediate",
            place: map[string]Piece{"e2": NewPawn(White), "e3": NewPawn(Black)},
            color: White, from: "e2", to: "e4",
            want: false,
        },
        {
            name:  "white double push blocked at destination",
            place: map[string]Piece{"e2": NewPawn(White), "e4": NewPawn(Black)},
            color: White, from: "e2", to: "e4",
            want: false,
        },
        {
            name:  "white backward push",
            place: map[string]Piece{"e3": NewPawn(White)},
            color: White, from: "e3", to: "e2",
            want: false,
        },
        {
            name:  "white sideways",
            place: map[string]Piece{"e3": NewPawn(White)},
            color: White, from: "e3", to: "f3",
            want: false,
        },
        {
            name:  "white stays in place",
            place: map[string]Piece{"e3": NewPawn(White)},
            color: White, from: "e3", to: "e3",
            want: false,
        },
        {
            name:  "white diagonal capture",
            place: map[string]Piece{"e4": NewPawn(White), "d5": NewPawn(Black)},
            color: White, from: "e4", to: "d5",
            want: true,
        },
        {
            name:  "white diagonal own pawn",
            place: map[string]Piece{"e4": NewPawn(White), "d5": NewPawn(White)},
            color: White, from: "e4", to: "d5",
            want: false,
        },
        {
            name:  "white diagonal empty no en passant",
            place: map[string]Piece{"e4": NewPawn(White)},
            color: White, from: "e4", to: "d5",
            want: false,
        },
        {
            name:  "white en passant",
            place: map[string]Piece{"e5": NewPawn(White), "d5": NewPawn(Black)},
            last:  chess.Move{From: sq("d7"), To: sq("d5")},
            color: White, from: "e5", to: "d6",
            want: true,
        },
        {
            name:  "white en passant target moved only one square",
            place: map[string]Piece{"e5": NewPawn(White), "d5": NewPawn(Black)},
            last:  chess.Move{From: sq("d6"), To: sq("d5")},
            color: White, from: "e5", to: "d6",
            want: false,
        },
        {
            name:  "white en passant last move different file",
            place: map[string]Piece{"e5": NewPawn(White), "d5": NewPawn(Black)},
            last:  chess.Move{From: sq("a7"), To: sq("a5")},
            color: White, from: "e5", to: "d6",
            want: false,
        },
        {
            name:  "white en passant adjacent is own pawn",
            place: map[string]Piece{"e5": NewPawn(White), "d5": NewPawn(White)},
            last:  chess.Move{From: sq("d7"), To: sq("d5")},
            color: White, from: "e5", to: "d6",
            want: false,
        },
        {
            name:  "black single push",
            place: map[string]Piece{"e7": NewPawn(Black)},
            color: Black, from: "e7", to: "e6",
            want: true,
        },
        {
            name:  "black double push from start",
            place: map[string]Piece{"e7": NewPawn(Black)},
            color: Black, from: "e7", to: "e5",
            want: true,
        },
        {
            name:  "black backward push",
            place: map[string]Piece{"e6": NewPawn(Black)},
            color: Black, from: "e6", to: "e7",
            want: false,
        },
        {
            name:  "black diagonal capture",
            place: map[string]Piece{"e5": NewPawn(Black), "d4": NewPawn(White)},
            color: Black, from: "e5", to: "d4",
            want: true,
        },
        {
            name:  "black en passant",
            place: map[string]Piece{"e4": NewPawn(Black), "d4": NewPawn(White)},
            last:  chess.Move{From: sq("d2"), To: sq("d4")},
            color: Black, from: "e4", to: "d3",
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
            b.lastMove = tc.last

            from := sq(tc.from)
            pawn, ok := b.squares[from.Row][from.Col].(*Pawn)
            if !ok {
                t.Fatalf("no pawn at %s", tc.from)
            }

            _, got := pawn.IsValidMove(&b, chess.Move{
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

func TestPawnExecuteSinglePush(t *testing.T) {
    var b Board
    pawn := NewPawn(White)
    b.squares[1][4] = pawn

    exec, valid := pawn.IsValidMove(&b, chess.Move{
        PlayerColor: White,
        From:        sq("e2"),
        To:          sq("e3"),
    })
    if !valid {
        t.Fatal("single push flagged invalid")
    }
    exec()

    if b.squares[1][4] != nil {
        t.Errorf("origin not cleared")
    }
    if b.squares[2][4] != pawn {
        t.Errorf("pawn not at destination")
    }
}

func TestPawnExecuteDoublePush(t *testing.T) {
    var b Board
    pawn := NewPawn(White)
    b.squares[1][4] = pawn

    exec, valid := pawn.IsValidMove(&b, chess.Move{
        PlayerColor: White,
        From:        sq("e2"),
        To:          sq("e4"),
    })
    if !valid {
        t.Fatal("double push flagged invalid")
    }
    exec()

    if b.squares[1][4] != nil {
        t.Errorf("origin not cleared")
    }
    if b.squares[3][4] != pawn {
        t.Errorf("pawn not at destination")
    }
}

func TestPawnExecuteCapture(t *testing.T) {
    var b Board
    pawn := NewPawn(White)
    enemy := NewPawn(Black)
    b.squares[3][4] = pawn
    b.squares[4][3] = enemy

    exec, valid := pawn.IsValidMove(&b, chess.Move{
        PlayerColor: White,
        From:        sq("e4"),
        To:          sq("d5"),
    })
    if !valid {
        t.Fatal("capture flagged invalid")
    }
    exec()

    if b.squares[3][4] != nil {
        t.Errorf("origin not cleared")
    }
    if b.squares[4][3] != pawn {
        t.Errorf("pawn did not replace captured piece")
    }
}

func TestPawnExecuteEnPassant(t *testing.T) {
    var b Board
    white := NewPawn(White)
    black := NewPawn(Black)
    b.squares[4][4] = white
    b.squares[4][3] = black
    b.lastMove = chess.Move{From: sq("d7"), To: sq("d5")}

    exec, valid := white.IsValidMove(&b, chess.Move{
        PlayerColor: White,
        From:        sq("e5"),
        To:          sq("d6"),
    })
    if !valid {
        t.Fatal("en passant flagged invalid")
    }
    exec()

    if b.squares[4][4] != nil {
        t.Errorf("origin not cleared")
    }
    if b.squares[5][3] != white {
        t.Errorf("attacker not at destination")
    }
    if b.squares[4][3] != nil {
        t.Errorf("captured pawn not removed")
    }
}
