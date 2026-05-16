CXX      := g++
CXXFLAGS := -std=c++17 -Wall -O2
SFML     := $(shell pkg-config --cflags --libs sfml-all)

SRC := client_ui/main.cpp \
       client_ui/UI.cpp \
       chess_logic/board.cpp \
       client_net/ChessNetwork.cpp

BIN := client

$(BIN): $(SRC)
	$(CXX) $(CXXFLAGS) $(SRC) -o $(BIN) $(SFML) $(LDLIBS)

clean:
	rm -f $(BIN)

.PHONY: clean
