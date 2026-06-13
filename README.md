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

**`PREGAME`**

| Direction | Message                       | Meaning                      | Notes                                                                                                                                                                                |
| --------- | ----------------------------- | ---------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| S→C       | `WHTE <your_time> <opp_time>` | Match started, you're white. | `<your_time>` and `<opp_time>` show how much time you and your opponent have left in seconds. The format is a floating-point value with 3 decimal places of precision, e.g. `6.123`. |
| S→C       | `BLCK <your_time> <opp_time>` | Match started, you're black. | See `WHTE`.                                                                                                                                                                          |

**`IN GAME`**

| Direction | Message                                               | Meaning                      | Notes                                                                                                                                                                                                                                         |
| --------- | ----------------------------------------------------- | ---------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| C→S       | `MOVE <from> <to> <promotion>`                        | Your move.                   | `<from>` and `<to>` both consist of one letter (a-h) and one number (1-8), e.g. `a2`, `e5`. `<promotion>` expresses pawn promotion and can take one of the following values: `-` (none), `Q` (Queen), `R` (Rook), `B` (Bishop), `N` (Knight). |
| S→C       | `ACPT`                                                | Your move was valid.         |                                                                                                                                                                                                                                               |
| S→C       | `RJCT`                                                | Your move was invalid.       |                                                                                                                                                                                                                                               |
| S→C       | `MOVE <from> <to> <promotion> <your_time> <opp_time>` | Opponent's move.             | See `MOVE` (C→S) and `WHTE`.                                                                                                                                                                                                                  |
| S→C       | `INVL`                                                | Unrecognized message format. |                                                                                                                                                                                                                                               |
| C→S       | `RSGN`                                                | Resign.                      |                                                                                                                                                                                                                                               |
| S→C       | `EOM <reason>`                                        | Match ended.                 | This is not a normal data message - it is the reason string of the WebSocket close frame.  The `<reason>` payload is the string form of [`classical.Outcome`](server/chess/classical/outcome.go), or a custom string if outcome is unknown.   |

> The server currently closes the connection upon receiving a message larger than 64 bytes.

## License

MIT License.
