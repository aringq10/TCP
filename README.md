# TCP — The Chess Project

A networked chess game client in C++ and remote server in Go.

> **Stack:** SFML 3.0.2

## TP2 goals

| Module            | Owner | Responsibilities                                                                                     |
| ----------------- | ----- | ---------------------------------------------------------------------------------------------------- |
| UI                | GF    | Draw board with a single piece. Make the piece interactive; call logic and network modules on move.  |
| Chess logic       | JK    | Create `Board` and `Piece` classes. Expose `make_move` and `get_valid_moves` from `Board` to the UI. |
| Client networking | AP    | Connect/disconnect to remote on app start/close; expose `send_move` to the UI.                       |
| Server networking | AC    | Expose protocol commands to the client network module (see below).                                   |

### Server protocol

- **Game State** — the state the game is in.
- **Client sends** — messages the client is allowed to send while in that state.
- **Server responds** — possible server replies to the corresponding client message.

| Game State | Direction         | Message                     |
|------------|-------------------|-----------------------------|
| `IN GAME`  | Client → Server   | `MOVE <from> <to>`          |
| `IN GAME`  | Server → Client   | `MOVE ACCEPTED`             |
| `IN GAME`  | Server → Client   | `MOVE REJECTED`             |
| `IN GAME`  | Server → Opponent | `OPPONENT_MOVE <from> <to>` |

## License

MIT License.
