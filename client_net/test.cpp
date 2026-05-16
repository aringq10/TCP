#include "ChessNetwork.hpp"
#include "../chess_logic/board.h"

#include <iostream>

void handle_event(Event e);
char piece_to_char(Piece piece);
void draw_board(const Board& board);

Board board;

int main() {
  ChessNetwork chess_network;

  if (chess_network.connect("aringas.dev", 6767, handle_event)) {
    std::cout << "Test connection succeeded" << std::endl;
  } else{
    std::cout << "Test connection failed" << std::endl;
  }
  std::string from, to;
  draw_board(board);
  while(std::cin >> from >> to) {
    bool valid = board.makeMove(from, to);
    if (valid && chess_network.send_move(from.c_str(), to.c_str())) {
      std::cout << "Move sent" << std::endl;
    } else {
      std::cout << "Failed to send move" << std::endl;
    }
  }

  std::cin >> std::ws; // Wait before exiting
  chess_network.disconnect();
  return 0;
} 

char piece_to_char(Piece piece) {
  switch (piece) {
    case WHITE_PAWN: return 'P';
    case WHITE_ROOK: return 'R';
    case WHITE_KNIGHT: return 'N';
    case WHITE_BISHOP: return 'B';
    case WHITE_QUEEN: return 'Q';
    case WHITE_KING: return 'K';
    case BLACK_PAWN: return 'p';
    case BLACK_ROOK: return 'r';
    case BLACK_KNIGHT: return 'n';
    case BLACK_BISHOP: return 'b';
    case BLACK_QUEEN: return 'q';
    case BLACK_KING: return 'k';
    case EMPTY:
    default:
      return '.';
  }
}


// temp function to draw the board in the console
void draw_board(const Board& board) {
  const Piece (&position)[8][8] = board.getBoard();

  std::cout << "\n  a b c d e f g h\n";

  for (int y = 0; y < 8; y++) {
    std::cout << 8 - y << ' ';

    for (int x = 0; x < 8; x++) {
      std::cout << piece_to_char(position[y][x]) << ' ';
    }

    std::cout << 8 - y << '\n';
  }

  std::cout << "  a b c d e f g h\n\n";
}

void handle_event(Event e) {
  std::cout << "Handler called" << std::endl;
  switch (e.type) {
    case OPPONENT_MOVE:
      std::cout << "Opponent move event received" << std::endl;
      std::cout << e.received_message << std::endl;
      board.makeOppMove(e.from, e.to);
      draw_board(board);
      break;

    case MOVE_ACCEPTED:
      std::cout << "Move accepted" << std::endl;
      draw_board(board);
      break;

    case MOVE_REJECTED:
      std::cout << "MOVE REJECTED" << std::endl;
      break;

    case INVALID:
      std::cout << "INVALID MOVE" << std::endl;
      break;

    case WHITE:
      std::cout << "You are playing as WHITE" << std::endl;
      board.setColor(Color::WHITE);
      break;

    case BLACK:
      std::cout << "You are playing as BLACK" << std::endl;
      board.setColor(Color::BLACK);
      break;

    // Handle other event types...
    default:
      std::cout << "Unknown event type received" << std::endl;
      break;
  }
}
