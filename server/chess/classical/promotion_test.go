package classical

import (
    "testing"
)

func setupBoard(place map[string]Piece, turn Color) *Board {
    var b Board
    for s, p := range place {
        q := sq(s)
        b.squares[q.Row][q.Col] = p
    }
    b.whoseTurn = turn
    return &b
}

func TestMakeMovePromotion(t *testing.T) {
    cases := []struct {
        name      string
        place     map[string]Piece
        turn      Color
        move      string
        wantOK    bool
        wantPiece string // String() of the piece expected on the destination square
    }{
        {
            name:  "white push promote to queen",
            place: map[string]Piece{"e7": NewPawn(White)},
            turn:  White, move: "e7 e8 Q",
            wantOK: true, wantPiece: "Q",
        },
        {
            name:  "white push promote to rook",
            place: map[string]Piece{"e7": NewPawn(White)},
            turn:  White, move: "e7 e8 R",
            wantOK: true, wantPiece: "R",
        },
        {
            name:  "white push promote to bishop",
            place: map[string]Piece{"e7": NewPawn(White)},
            turn:  White, move: "e7 e8 B",
            wantOK: true, wantPiece: "B",
        },
        {
            name:  "white push promote to knight",
            place: map[string]Piece{"e7": NewPawn(White)},
            turn:  White, move: "e7 e8 N",
            wantOK: true, wantPiece: "N",
        },
        {
            name:  "black push promote to queen",
            place: map[string]Piece{"e2": NewPawn(Black)},
            turn:  Black, move: "e2 e1 Q",
            wantOK: true, wantPiece: "q",
        },
        {
            name:  "white capture promote to queen",
            place: map[string]Piece{"a7": NewPawn(White), "b8": NewRook(Black)},
            turn:  White, move: "a7 b8 Q",
            wantOK: true, wantPiece: "Q",
        },
        {
            name:  "missing promotion on last rank rejected",
            place: map[string]Piece{"e7": NewPawn(White)},
            turn:  White, move: "e7 e8 -",
            wantOK: false,
        },
        {
            name:  "promotion flag off last rank rejected",
            place: map[string]Piece{"e2": NewPawn(White)},
            turn:  White, move: "e2 e4 Q",
            wantOK: false,
        },
        {
            name:  "promotion flag on non-pawn rejected",
            place: map[string]Piece{"a1": NewRook(White)},
            turn:  White, move: "a1 a8 Q",
            wantOK: false,
        },
        {
            name:  "promotion push onto occupied square rejected",
            place: map[string]Piece{"e7": NewPawn(White), "e8": NewPawn(Black)},
            turn:  White, move: "e7 e8 Q",
            wantOK: false,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            b := setupBoard(tc.place, tc.turn)

            if got := b.MakeMove(tc.turn, tc.move); got != tc.wantOK {
                t.Fatalf("MakeMove(%q) = %v, want %v", tc.move, got, tc.wantOK)
            }
            if !tc.wantOK {
                return
            }

            landSq := tc.move[3:5]
            land := sq(landSq)
            got := b.squares[land.Row][land.Col]
            if got == nil {
                t.Fatalf("no piece on %s after promotion", landSq)
            }
            if got.String() != tc.wantPiece {
                t.Errorf("piece on %s = %q, want %q", landSq, got.String(), tc.wantPiece)
            }
            if _, stillPawn := got.(*Pawn); stillPawn {
                t.Errorf("piece on %s is still a pawn", landSq)
            }

            from := sq(tc.move[0:2])
            if b.squares[from.Row][from.Col] != nil {
                t.Errorf("origin %s not cleared", tc.move[0:2])
            }
        })
    }
}

func TestMakeMovePromotionIntoCheckRejected(t *testing.T) {
    // White king on a7, pawn on b7, black rook on h7: promoting the pawn vacates
    // b7 and exposes the king along the 7th rank, so the move must be rejected.
    b := setupBoard(map[string]Piece{
        "a7": NewKing(White),
        "b7": NewPawn(White),
        "h7": NewRook(Black),
    }, White)

    if b.MakeMove(White, "b7 b8 Q") {
        t.Fatal("promotion exposing own king was accepted")
    }

    from := sq("b7")
    if _, ok := b.squares[from.Row][from.Col].(*Pawn); !ok {
        t.Errorf("pawn not restored on b7")
    }
    to := sq("b8")
    if b.squares[to.Row][to.Col] != nil {
        t.Errorf("b8 not restored to empty")
    }
}

func TestMakeMovePromotionDeliversCheckmate(t *testing.T) {
    // White king f7 supports a queen promotion on g8, mating the black king on h8.
    b := setupBoard(map[string]Piece{
        "h8": NewKing(Black),
        "f7": NewKing(White),
        "g7": NewPawn(White),
    }, White)

    if !b.MakeMove(White, "g7 g8 Q") {
        t.Fatal("promotion delivering checkmate was rejected")
    }
    if !b.IsOver() {
        t.Fatal("game not marked over after checkmating promotion")
    }
    if got := b.Outcome(); got.Termination != Checkmate || got.Result != WhiteWins {
        t.Errorf("outcome = %+v, want checkmate / white wins", got)
    }
}
