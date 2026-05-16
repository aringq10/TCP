#include "ChessNetwork.hpp"

#include <boost/beast/core.hpp>
#include <boost/beast/version.hpp>

#include <iostream>

ChessNetwork::ChessNetwork()
  : resolver_(io_context_),
    websocket_(io_context_),
    is_connected_(false),
    pending_move_(false) {
}

ChessNetwork::~ChessNetwork() {
  disconnect();
}

bool ChessNetwork::connect(const std::string& ip_address, std::uint16_t port, const EventHandler& handle) {
  try {
    if (is_connected_) {
      return true;
    }
    // Event e {OTHER};
    // handle(e);

    const auto endpoints = resolver_.resolve(ip_address, std::to_string(port));
    const auto endpoint = net::connect(websocket_.next_layer(), endpoints);

    const std::string host_header = ip_address + ":" + std::to_string(endpoint.port());

    websocket_.handshake(host_header, "/ws");
    is_connected_ = true;

    std::cout << "Connected to WebSocket server at " << host_header << std::endl;

    handler_ = handle;
    
    receive_thread_ = std::thread([this]()-> void {
      receive_loop();
    });

    return true;
  } catch (const std::exception& exception) {
    std::cerr << "WebSocket connection failed: " << exception.what() << std::endl;
    return false;
  }
}

void ChessNetwork::disconnect() {
  if (!is_connected_) {
    return;
  }

  try {
    websocket_.close(websocket::close_code::normal);
    is_connected_ = false;
    std::cout << "Disconnected successfully" << std::endl;
  } catch (const std::exception&) {
  }
}

bool ChessNetwork::send_move(const char from[2], const char to[2]) {
  if(!is_connected_ || has_pending_move() || from == nullptr || to == nullptr) {
    return false;
  }
  
  try {
    pending_move_ = true;

    std::string messageType = "MOVE";
    std::string message = std::string(from, 2) + " " + std::string(to, 2);
    websocket_.write(net::buffer(messageType + " " + message));

    return true;
  } catch (const std::exception& e) {
    pending_move_ = false;
    std::cerr << "Failed to send move: " << e.what() << std::endl;
    return false;
  }
}

void ChessNetwork::receive_loop() {
  try {
    while (is_connected_) {
      beast::flat_buffer buffer;
      websocket_.read(buffer);

      std::string message = beast::buffers_to_string(buffer.data());
      std::string parsed_message;
      std::string messageType;

      Event e;
      MessageType type = OTHER;

      std::cout << "Received message from server: " << message << std::endl;

      if (message.length() >= 4) {
        messageType = message.substr(0, 4);
        parsed_message = "";

        if (message.length() > 5) {
          parsed_message = message.substr(5);
        }

        if (messageType == "MOVE") {
          parse_move(parsed_message, e);
          type = OPPONENT_MOVE;
        } 
        else if (messageType == "ACPT") {
          type = MOVE_ACCEPTED;
          pending_move_ = false;
        }
        else if (messageType == "RJCT") {
          type = MOVE_REJECTED;
          pending_move_ = false;
        }
        else if (messageType == "INVL") {
          type = INVALID;
          pending_move_ = false;
        }
      }
    
      e.type = type;
      e.received_message = parsed_message;
      handler_(e);
    }
  } catch (const std::exception& e) {
    std::cerr << "Error in receive loop: " << e.what() << std::endl;
  }
}

void ChessNetwork::parse_move(const std::string& message, Event& e) {
  std::stringstream ss(message);
  std::string from, to;
  ss >> from >> to;
  e.from = from;
  e.to = to;
}

bool ChessNetwork::has_pending_move() const {
  return pending_move_;
}