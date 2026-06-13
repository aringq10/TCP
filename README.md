# TCP — The Chess Project

A networked chess game client in C++ and remote server in Go.

> **Stack:** SFML 3.0.2 Golang ASIO/Boost

## TP2 goals

| Module            | Owner | Responsibilities                                                                                    |
| ----------------- | ----- | --------------------------------------------------------------------------------------------------- |
| UI                | GF    | Draw board with a single piece. Make the piece interactive; call logic and network modules on move. |
| Chess logic       | JK    | Create `Board` class. Expose `makeMove` and `makeOppMove` from `Board` to the UI.                   |
| Client networking | AP    | Expose `send_move`, `connect` and `disconnect` to the UI.                                           |
| Server networking | AC    | Expose protocol commands to the client network module (see below), server side move validation.     |

## TP3 goals

| Module                     | Owner  | Responsibilities                                                                                  |
| -------------------------- | ------ | ------------------------------------------------------------------------------------------------- |
| UI                         | GF     | Pieces represented by PNGs. Main menu. Buttons for Resign/Leave, Join Queue etc.                  |
| Chess logic                | JK     | Create a `Piece` class instead of enums. Add movement logic for all other classical chess pieces. |
| Client & Server Networking | AP, AC | Come up with remaining protocol messages. Implement in code.                                      |

### Server protocol

- **Game State (the table title)** — the state the game is in.
- **Direction** — who sends the message to whom. S - server, C - client.
- **Message** — the wire format of the message.

`<wt>` - how much time white has left in seconds. Floating point, 3 decimal places of precision, e.g. `6.123`.
`<bt>` - same as `<wt>` but for black.
`<from>`, `<to>` - squares from where to where you want to move. Consists of one letter (a-h) and one number (1-8), e.g. `a2`, `e5`.
`<promotion>` - expresses pawn promotion and can take one of the following values: `-` (none), `Q` (Queen), `R` (Rook), `B` (Bishop), `N` (Knight).
`<reason>` - end of match reason, the string form of [`classical.Outcome`](server/chess/classical/outcome.go).

**`PREGAME`**

| Direction | Message          | Meaning                      |
| --------- | ---------------- | ---------------------------- |
| S→C       | `WHTE <wt> <bt>` | Match started, you're white. |
| S→C       | `BLCK <wt> <bt>` | Match started, you're black. |

**`IN GAME`**

| Direction | Message                                  | Meaning                      |
| --------- | ---------------------------------------- | ---------------------------- |
| C→S       | `MOVE <from> <to> <promotion>`           | Your move.                   |
| S→C       | `ACPT <wt> <bt>`                         | Your move was valid.         |
| S→C       | `RJCT <wt> <bt>`                         | Your move was invalid.       |
| S→C       | `MOVE <from> <to> <promotion> <wt> <bt>` | Opponent's move.             |
| S→C       | `INVL`                                   | Unrecognized message format. |
| C→S       | `RSGN`                                   | Resign.                      |
| S→C       | `ENDM <reason>`                          | Match ended.                 |

> The server currently closes the connection for the following reasons:
> Receives a message larger than 64 bytes.
> Receives too many messages to which it responded with `INVL`. 

## License

MIT License.
