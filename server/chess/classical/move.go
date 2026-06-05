package classical

type Move struct {
    PlayerColor Color
    From Square
    To Square
}

func ParseMove(c Color, s string) (m Move, ok bool) {
    if len(s) != 5 {
        return m, false
    }

    from, ok1 := ParseSquare(s[0:2])
    to, ok2 := ParseSquare(s[3:5])

    if !ok1 || !ok2 || !c.Valid() {
        return m, false
    }

    return Move{PlayerColor: c, From: from, To: to}, true
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
