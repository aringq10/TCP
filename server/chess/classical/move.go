package classical

type Move struct {
    PlayerColor Color
    From Square
    To Square
}

func NewMove(playerColor Color, from Square, to Square) Move {
    return Move{PlayerColor: playerColor, From: from, To: to}
}

type Square struct {
    Row int
    Col int
}

func ParseSquare(s string) (sq Square, ok bool) {
    if len(s) != 2 {
        return sq, false
    }

    col, row := int(s[0] - 'a'), int(s[1] - '1')
    if col < 0 || col > 7 || row < 0 || row > 7 {
        return sq, false
    }

    sq.Row = row
    sq.Col = col
    return sq, true
}
