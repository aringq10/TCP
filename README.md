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
- **Direction** — who sends the message to whom.
- **Message** — the wire format of the message.

**`PREGAME`**

| Direction         | Message |
|-------------------|---------|
| Server → Opponent | `WHTE`  |
| Server → Opponent | `BLCK`  |

**`IN GAME`**

| Direction         | Message            |
|-------------------|--------------------|
| Client → Server   | `MOVE <from> <to>` |
| Server → Client   | `ACPT`             |
| Server → Client   | `RJCT`             |
| Server → Opponent | `MOVE <from> <to>` |
| Server → Both     | `EOM: <reason>`    |

> `EOM: <reason>` is not a normal data message — it is the reason string of the WebSocket close frame the server sends when ending the match.
> The `<reason>` payload is the string form of [`classical.Outcome`](server/chess/classical/outcome.go) — a `Result` and an optional `Termination`.

> Will be expanded in the course of TP3.

## License

MIT License.
