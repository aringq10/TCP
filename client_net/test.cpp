#include "ChessNetwork.hpp"
#include "Board.h"

#include <iostream>

void handle_event(Event e);
char piece_to_char(Piece piece);
void draw_board(const Board& board);

Board board;

int main() {
  ChessNetwork chess_network;

  if (chess_network.connect("193.219.91.103", 6767, handle_event)) {
    std::cout << "Test connection succeeded" << std::endl;
  } else{
    std::cout << "Test connection failed" << std::endl;
  }
  // draw_board(board);

  // // test sent_move
  // if(board.makeMove("a5", "a6")) {
  //   draw_board(board);
  //   chess_network.send_move("a5", "a6");
  // }

  // if (!chess_network.send_move("a6", "a7")) {
  //   std::cout << "Failed to send move" << std::endl;
  // }
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
      std::cout << e.type <<  " " << e.received_message << std::endl;
      board.makeMove(e.from, e.to);
      break;

    case MOVE_ACCEPTED:
      std::cout << "Move accepted" << std::endl;
      break;

    case MOVE_REJECTED:
      std::cout << "MOVE REJECTED" << std::endl;
      break;

    // Handle other event types...
    default:
      std::cout << "Unknown event type received" << std::endl;
      break;
  }
}
