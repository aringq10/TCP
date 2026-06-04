package classical

type Result int

const (
    Ongoing Result = iota
    WhiteWins
    BlackWins
    Draw
)

var resultNames = map[Result]string{
    Ongoing:   "*",
    WhiteWins: "1-0",
    BlackWins: "0-1",
    Draw:      "1/2-1/2",
}

func (r Result) String() string { return resultNames[r] }

func (r Result) IsOver() bool { return r != Ongoing }

func resultForWinner(c Color) Result {
    if c == White {
        return WhiteWins
    }
    return BlackWins
}

type Termination int

const (
    NoTermination Termination = iota // Ongoing
    Checkmate
    Stalemate
    InsufficientMaterial
    FiftyMoveRule
    ThreefoldRepetition
    Resignation
    TimeForfeit
    DrawAgreement
)

var terminationNames = map[Termination]string{
    NoTermination:        "ongoing",
    Checkmate:            "checkmate",
    Stalemate:            "stalemate",
    InsufficientMaterial: "insufficient material",
    FiftyMoveRule:        "fifty-move rule",
    ThreefoldRepetition:  "threefold repetition",
    Resignation:          "resignation",
    TimeForfeit:          "time forfeit",
    DrawAgreement:        "draw by agreement",
}

func (t Termination) String() string { return terminationNames[t] }

type Outcome struct {
    Result      Result
    Termination Termination
}

func (o Outcome) String() string {
    if o.Termination == NoTermination {
        return o.Result.String()
    }
    return o.Result.String() + " (" + o.Termination.String() + ")"
}

// NewOutcome builds the Outcome for a termination. For decisive terminations
// (Checkmate, Resignation, TimeForfeit) winner names the side that won; it is
// ignored for draw terminations.
func NewOutcome(t Termination, winner Color) Outcome {
    switch t {
    case Checkmate, Resignation, TimeForfeit:
        return Outcome{Result: resultForWinner(winner), Termination: t}
    case Stalemate, InsufficientMaterial, FiftyMoveRule, ThreefoldRepetition, DrawAgreement:
        return Outcome{Result: Draw, Termination: t}
    default: // NoTermination
        return Outcome{}
    }
}
