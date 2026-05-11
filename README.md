# TCP - The Chess Project

Currently using SFML 3.0.2

## TP2 goals

UI(GF) - Draw board with a single piece. Make the piece interactive, call logic and network modules on move.
Chess logic(JK) - Create Board and Piece classes. Expose make_move and get_valid_moves method from Board to UI.
Client Networking(AP) - Connect/disconnect to remote on app start/close, expose send_move method to UI.
Server Networking(AP) - Expose protocol commands to client network module:
```
  MOVE SEND <a-h><1-8> <a-h><1-8>
  MOVE RECEIVED <a-h><1-8> <a-h><1-8>
```

MIT License
