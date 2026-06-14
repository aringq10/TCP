package classical

import (
    "testing"
)

// moved returns p after marking it as having moved, so castling eligibility
// (which requires MoveCount() == 0) can be exercised.
func moved(p Piece) Piece {
    p.IncMoveCount()
    return p
}

func TestKingCastleIsValidMove(t *testing.T) {
    cases := []struct {
        name  string
        place map[string]Piece
        from  string
        to    string
        want  bool
    }{
        // --- Legal castles, both sides, both colors ---
        {
            name:  "white king side",
            place: map[string]Piece{"e1": NewKing(White), "h1": NewRook(White)},
            from:  "e1", to: "g1", want: true,
        },
        {
            name:  "white queen side",
            place: map[string]Piece{"e1": NewKing(White), "a1": NewRook(White)},
            from:  "e1", to: "c1", want: true,
        },
        {
            name:  "black king side",
            place: map[string]Piece{"e8": NewKing(Black), "h8": NewRook(Black)},
            from:  "e8", to: "g8", want: true,
        },
        {
            name:  "black queen side",
            place: map[string]Piece{"e8": NewKing(Black), "a8": NewRook(Black)},
            from:  "e8", to: "c8", want: true,
        },

        // --- Blocked by a piece between king and rook ---
        {
            name:  "king side blocked on f1",
            place: map[string]Piece{"e1": NewKing(White), "h1": NewRook(White), "f1": NewBishop(White)},
            from:  "e1", to: "g1", want: false,
        },
        {
            name:  "king side blocked on g1",
            place: map[string]Piece{"e1": NewKing(White), "h1": NewRook(White), "g1": NewKnight(White)},
            from:  "e1", to: "g1", want: false,
        },
        {
            name:  "queen side blocked on d1",
            place: map[string]Piece{"e1": NewKing(White), "a1": NewRook(White), "d1": NewQueen(White)},
            from:  "e1", to: "c1", want: false,
        },
        {
            name:  "queen side blocked on c1",
            place: map[string]Piece{"e1": NewKing(White), "a1": NewRook(White), "c1": NewBishop(White)},
            from:  "e1", to: "c1", want: false,
        },
        {
            // b1 sits between the queen-side rook and king, so it must be empty.
            name:  "queen side blocked on b1",
            place: map[string]Piece{"e1": NewKing(White), "a1": NewRook(White), "b1": NewKnight(White)},
            from:  "e1", to: "c1", want: false,
        },
        {
            name:  "king side landing square holds enemy",
            place: map[string]Piece{"e1": NewKing(White), "h1": NewRook(White), "g1": NewPawn(Black)},
            from:  "e1", to: "g1", want: false,
        },

        // --- King or rook has already moved ---
        {
            name:  "king has moved",
            place: map[string]Piece{"e1": moved(NewKing(White)), "h1": NewRook(White)},
            from:  "e1", to: "g1", want: false,
        },
        {
            name:  "rook has moved",
            place: map[string]Piece{"e1": NewKing(White), "h1": moved(NewRook(White))},
            from:  "e1", to: "g1", want: false,
        },
        {
            name:  "no rook on corner",
            place: map[string]Piece{"e1": NewKing(White)},
            from:  "e1", to: "g1", want: false,
        },
        {
            name:  "non-rook on corner",
            place: map[string]Piece{"e1": NewKing(White), "h1": NewBishop(White)},
            from:  "e1", to: "g1", want: false,
        },

        // --- King may not castle out of, through, or into check ---
        {
            name:  "king in check",
            place: map[string]Piece{"e1": NewKing(White), "h1": NewRook(White), "e3": NewRook(Black)},
            from:  "e1", to: "g1", want: false,
        },
        {
            name:  "king side passes through attacked f1",
            place: map[string]Piece{"e1": NewKing(White), "h1": NewRook(White), "f3": NewRook(Black)},
            from:  "e1", to: "g1", want: false,
        },
        {
            name:  "king side lands on attacked g1",
            place: map[string]Piece{"e1": NewKing(White), "h1": NewRook(White), "g3": NewRook(Black)},
            from:  "e1", to: "g1", want: false,
        },
        {
            name:  "queen side passes through attacked d1",
            place: map[string]Piece{"e1": NewKing(White), "a1": NewRook(White), "d3": NewRook(Black)},
            from:  "e1", to: "c1", want: false,
        },
        {
            name:  "queen side lands on attacked c1",
            place: map[string]Piece{"e1": NewKing(White), "a1": NewRook(White), "c3": NewRook(Black)},
            from:  "e1", to: "c1", want: false,
        },
        {
            // b1 is not crossed by the king, so an attack on it does not block
            // queen-side castling (only occupancy does).
            name:  "queen side allowed when only b1 is attacked",
            place: map[string]Piece{"e1": NewKing(White), "a1": NewRook(White), "b3": NewRook(Black)},
            from:  "e1", to: "c1", want: true,
        },

        // --- Not a castle: king off its home rank ---
        {
            name:  "two-square move off back rank",
            place: map[string]Piece{"e2": NewKing(White), "h2": NewRook(White)},
            from:  "e2", to: "g2", want: false,
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

func TestKingCastleExecuteKingSide(t *testing.T) {
    var b Board
    king := NewKing(White)
    rook := NewRook(White)
    b.squares[0][4] = king
    b.squares[0][7] = rook

    exec, valid := king.IsValidMove(&b, sq("e1"), sq("g1"))
    if !valid {
        t.Fatal("king-side castle flagged invalid")
    }
    exec()

    if b.squares[0][6] != king {
        t.Errorf("king not on g1")
    }
    if b.squares[0][5] != rook {
        t.Errorf("rook not on f1")
    }
    if b.squares[0][4] != nil {
        t.Errorf("e1 not cleared")
    }
    if b.squares[0][7] != nil {
        t.Errorf("h1 not cleared")
    }
}

func TestKingCastleExecuteQueenSide(t *testing.T) {
    var b Board
    king := NewKing(White)
    rook := NewRook(White)
    b.squares[0][4] = king
    b.squares[0][0] = rook

    exec, valid := king.IsValidMove(&b, sq("e1"), sq("c1"))
    if !valid {
        t.Fatal("queen-side castle flagged invalid")
    }
    exec()

    if b.squares[0][2] != king {
        t.Errorf("king not on c1")
    }
    if b.squares[0][3] != rook {
        t.Errorf("rook not on d1")
    }
    if b.squares[0][4] != nil {
        t.Errorf("e1 not cleared")
    }
    if b.squares[0][0] != nil {
        t.Errorf("a1 not cleared")
    }
}

// The castle exec closure must not mutate piece move counts. MakeMove and
// terminationFor execute moves speculatively and roll them back by restoring
// the squares array, which does not undo changes to a piece's internal state.
// If exec bumped the rook's move count, a speculative castle would permanently
// make the rook ineligible to castle.
func TestCastleExecLeavesRookMoveCountUntouched(t *testing.T) {
    var b Board
    king := NewKing(White)
    rook := NewRook(White)
    b.squares[0][4] = king
    b.squares[0][7] = rook

    exec, valid := king.IsValidMove(&b, sq("e1"), sq("g1"))
    if !valid {
        t.Fatal("king-side castle flagged invalid")
    }
    exec()

    if rook.MoveCount() != 0 {
        t.Errorf("rook MoveCount() = %d after castle exec, want 0", rook.MoveCount())
    }
}

func TestMakeMoveCastleKingSide(t *testing.T) {
    b := setupBoard(map[string]Piece{
        "e1": NewKing(White),
        "h1": NewRook(White),
        "e8": NewKing(Black),
    }, White)

    if !makeMove(b, White, "e1 g1 -") {
        t.Fatal("king-side castle rejected by MakeMove")
    }
    if _, ok := b.squares[0][6].(*King); !ok {
        t.Errorf("king not on g1 after castle")
    }
    if _, ok := b.squares[0][5].(*Rook); !ok {
        t.Errorf("rook not on f1 after castle")
    }
}

func TestMakeMoveCastleQueenSide(t *testing.T) {
    b := setupBoard(map[string]Piece{
        "e8": NewKing(Black),
        "a8": NewRook(Black),
        "e1": NewKing(White),
    }, Black)

    if !makeMove(b, Black, "e8 c8 -") {
        t.Fatal("queen-side castle rejected by MakeMove")
    }
    if _, ok := b.squares[7][2].(*King); !ok {
        t.Errorf("king not on c8 after castle")
    }
    if _, ok := b.squares[7][3].(*Rook); !ok {
        t.Errorf("rook not on d8 after castle")
    }
}

// Castling must still be available after unrelated moves have been played, so
// long as the king and rook themselves never moved.
func TestMakeMoveCastleAfterUnrelatedMoves(t *testing.T) {
    b := setupBoard(map[string]Piece{
        "e1": NewKing(White),
        "h1": NewRook(White),
        "a2": NewPawn(White),
        "e8": NewKing(Black),
        "a7": NewPawn(Black),
    }, White)

    if !makeMove(b, White, "a2 a3 -") {
        t.Fatal("white pawn move rejected")
    }
    if !makeMove(b, Black, "a7 a6 -") {
        t.Fatal("black pawn move rejected")
    }
    if !makeMove(b, White, "e1 g1 -") {
        t.Fatal("king-side castle rejected after unrelated moves")
    }
    if _, ok := b.squares[0][6].(*King); !ok {
        t.Errorf("king not on g1 after castle")
    }
    if _, ok := b.squares[0][5].(*Rook); !ok {
        t.Errorf("rook not on f1 after castle")
    }
}
