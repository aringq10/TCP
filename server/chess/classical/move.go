package classical

type Promotion int

const (
    NoPromotion Promotion = iota
    ToQueen
    ToRook
    ToBishop
    ToKnight
)

var charToPromotion = map[rune]Promotion {
    '-': NoPromotion,
    'Q': ToQueen,
    'R': ToRook,
    'B': ToBishop,
    'N': ToKnight,
}

type Move struct {
    Color Color
    From Square
    To Square
    Promotion Promotion
}

func (m Move) IsPromotion() bool {
    return m.Promotion != NoPromotion
}

func ParseMove(c Color, s string) (m Move, ok bool) {
    if len(s) != 7 {
        return m, false
    }

    from, okFrom := ParseSquare(s[0:2])
    to, okTo := ParseSquare(s[3:5])

    if okFrom && okTo && s[5] == ' ' {
        if p, okPromo := charToPromotion[rune(s[6])]; okPromo {
            return Move{Color: c, From: from, To: to, Promotion: p}, true
        }
    }

    return m, false
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
