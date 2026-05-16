package chess

type Board [8][8]Piece

func NewBoard() *Board {
    var b Board
    for col := 0; col < 8; col++ {
        b[1][col] = &Pawn{basePiece{color: WHITE}}
        b[6][col] = &Pawn{basePiece{color: BLACK}}
    }
    return &b
}

func parseSquare(s string) (col, row int) {
    return int(s[0] - 'a'), int(s[1] - '1')
}

func (b *Board) LookupMove(from, to string) (piece *Piece, fromCol, fromRow, toCol, toRow int, ok bool) {
    fromCol, fromRow = parseSquare(from)
    toCol, toRow = parseSquare(to)
    if fromCol < 0 || fromCol > 7 || fromRow < 0 || fromRow > 7 ||
        toCol < 0 || toCol > 7 || toRow < 0 || toRow > 7 {
        return nil, 0, 0, 0, 0, false
    }
    return &b[fromRow][fromCol], fromCol, fromRow, toCol, toRow, true
}

func (b *Board) MakeMove(playerColor Color, from string, to string) bool {
    piece, fromCol, fromRow, toCol, toRow, ok := b.LookupMove(from, to)
    if !ok || piece == nil || (*piece).Color() != playerColor {
        return false
    }

    valid, exec := (*piece).IsValidMove(b, fromCol, fromRow, toCol, toRow)
    if !valid || exec == nil {
        return false
    }

    exec()
    return true
}
