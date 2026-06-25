package classical

func NewPawnsOnlyBoard() *Board {
    var b Board
    placePawns(&b.squares)
    return &b
}
