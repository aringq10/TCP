package classical

type Color int

const (
    White Color = iota
    Black
)

var colorNames = map[Color]string {
    White: "WHTE",
    Black: "BLCK",
}

func (c Color) String() string {
    return colorNames[c]
}

func oppositeColor(c Color) Color {
    if c == White {
        return Black
    } else {
        return White
    }
}
