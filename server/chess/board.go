package chess

type Move interface {
    IsMove()
}

type Board interface {
    MakeMove(m Move) bool
    String() string
}
