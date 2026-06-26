CXX      := g++
CXXFLAGS := -std=c++17 -Wall -O2
SFML     := $(shell pkg-config --cflags --libs sfml-all)

SRC := main.cpp \
       client_ui/UI.cpp \
       chess_logic/board.cpp \
       client_net/ChessNetwork.cpp \
       chess_logic/knight.cpp \
       chess_logic/pawn.cpp \
       chess_logic/rook.cpp \
       chess_logic/bishop.cpp \
       chess_logic/queen.cpp \
       chess_logic/king.cpp \

BIN := client

$(BIN): $(SRC)
	$(CXX) $(CXXFLAGS) $(SRC) -o $(BIN) $(SFML) $(LDLIBS)

clean:
	rm -f $(BIN)

.PHONY: clean
