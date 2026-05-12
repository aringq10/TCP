#include "ChessNetwork.hpp"

#include <iostream>

void handle_event(Event e);

int main() {
  ChessNetwork chess_network;

  if (chess_network.connect("193.219.91.103", 6767, handle_event)) {
    std::cout << "Test connection succeeded" << std::endl;
  } else{
    std::cout << "Test connection failed" << std::endl;
  }

  // test sent_move
  chess_network.send_move("e2", "e4");
  std::cin >> std::ws; // Wait before exiting
  return 0;
}

void handle_event(Event e) {
  std::cout << "Handler called" << std::endl;
  switch (e.type) {
    case MOVE_OK:
      std::cout << "Move OK event received" << std::endl;
      break;
    // Handle other event types...
    default:
      std::cout << "Unknown event type received" << std::endl;
      break;
  }
}
