package msg

import (
	"github.com/aringq10/TCP/server/chess/classical"
)

const (
    Black      = "BLCK"
    White      = "WHTE"
    Accept     = "ACPT"
    Reject     = "RJCT"
    Move       = "MOVE"
    Invalid    = "INVL"
    Resign     = "RSGN"
    EndOfMatch = "EOM"
)

var MatchColor = map[classical.Color]string {
    classical.White: White,
    classical.Black: Black,
}
