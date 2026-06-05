package classical

type Color int

const (
    White Color = iota
    Black
)

var colorNames = map[Color]string {
    White: "White",
    Black: "Black",
}

func (c Color) String() string {
    return colorNames[c]
}

func (c Color) Valid() bool {
    _, ok := colorNames[c]
    return ok
}

func OppositeColor(c Color) Color {
    if c == White {
        return Black
    } else {
        return White
    }
}
