package classical

type King struct {
    BasePiece
}

func NewKing(c Color) *King {
    return &King{NewBasePiece(c)}
}

func (k *King) String() string {
    if k.Color() == White {
        return "K"
    } else {
        return "k"
    }
}

func (k *King) IsValidMove(b *Board, from, to Square) (execute func(), valid bool) {
    s := &b.squares
    fromRow, fromCol, toRow, toCol := from.Row, from.Col, to.Row, to.Col

    execCastle, castleValid := k.isValidCastle(b, fromRow, fromCol, toRow, toCol)
    if castleValid {
        return execCastle, true
    }

    dRow := toRow - fromRow
    dCol := toCol - fromCol
    if abs(dRow) > 1 || abs(dCol) > 1 || abs(dRow) + abs(dCol) == 0 {
        return nil, false
    }

    target := s[toRow][toCol]

    if target != nil && target.Color() == k.Color() {
        return nil, false
    }

    return func() {
        s[toRow][toCol] = s[fromRow][fromCol]
        s[fromRow][fromCol] = nil
    }, true
}

func (k *King) isValidCastle(b *Board, fromRow, fromCol, toRow, toCol int) (execute func(), valid bool) {
    s := &b.squares
    var castleRow int
    if k.Color() == White {
        castleRow = 0
    } else {
        castleRow = 7
    }

    var mustBeEmpty []Square
    var mustNotBeAttacked []Square
    var isRookValid bool
    var isKingValid bool = k.MoveCount() == 0
    var moveRook func()

    if (fromRow == toRow) && (fromRow == castleRow) {
        var rook Piece
        if fromCol == 4 && toCol == 6 { // King side castle
            mustNotBeAttacked = []Square{{Row: castleRow, Col: 5}, {Row: castleRow, Col: 6}}
            rook = s[castleRow][7]
            moveRook = func() {
                s[castleRow][5] = s[castleRow][7]
                s[castleRow][7] = nil
            }
        } else if fromCol == 4 && toCol == 2 { // Queen side castle
            mustNotBeAttacked = []Square{{Row: castleRow, Col: 3}, {Row: castleRow, Col: 2}}
            mustBeEmpty = []Square{{Row: castleRow, Col: 1}}
            rook = s[castleRow][0]
            moveRook = func() {
                s[castleRow][3] = s[castleRow][0]
                s[castleRow][0] = nil
            }
        }
        mustBeEmpty = append(mustBeEmpty, mustNotBeAttacked...)
        if rook != nil {
            r, isRook := rook.(*Rook)
            if isRook && r.MoveCount() == 0 {
                isRookValid = true
            }
        }
    }

    if isKingValid && isRookValid && !b.isAttacked(Square{Row: castleRow, Col: 4}, OppositeColor(k.Color())) {
        isCastleValid := true
        for _, sq := range mustBeEmpty {
            if s[sq.Row][sq.Col] != nil {
                isCastleValid = false
                break
            }
        }
        for _, sq := range mustNotBeAttacked {
            if b.isAttacked(sq, OppositeColor(k.Color())) {
                isCastleValid = false
                break
            }
        }
        if isCastleValid {
            return func() {
                s[toRow][toCol] = s[fromRow][fromCol]
                s[fromRow][fromCol] = nil
                moveRook()
            }, true
        }
    }

    return nil, false
}
