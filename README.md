# TCP — The Chess Project

A networked chess game client in C++ and remote server in Go.

> **Stack:** SFML 3.0.2

## TP2 goals

| Module            | Owner | Responsibilities                                                                                     |
| ----------------- | ----- | ---------------------------------------------------------------------------------------------------- |
| UI                | GF    | Draw board with a single piece. Make the piece interactive; call logic and network modules on move.  |
| Chess logic       | JK    | Create `Board` class. Expose `makeMove` and `makeOppMove` from `Board` to the UI.                    |
| Client networking | AP    | Expose `send_move`, `connect` and `disconnect` to the UI.                                            |
| Server networking | AC    | Expose protocol commands to the client network module (see below), server side move validation.      |

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

| Direction         | Message                     |
|-------------------|-----------------------------|
| Client → Server   | `MOVE <from> <to>`          |
| Server → Client   | `MOVE ACCEPTED`             |
| Server → Client   | `MOVE REJECTED`             |
| Server → Opponent | `OPPONENT_MOVE <from> <to>` |

## License

MIT License.
