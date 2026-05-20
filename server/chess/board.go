package chess

type Color interface {
    String() string
}

type Square struct {
    Row int
    Col int
}

type Move struct {
    PlayerColor Color
    From Square
    To Square
}

func NewMove(playerColor Color, from Square, to Square) Move {
    return Move{PlayerColor: playerColor, From: from, To: to}
}

type Board interface {
    MakeMove(m Move) bool
    String() string
    ParseSquare(s string) (sq Square, ok bool)
}
