#include "ChessNetwork.hpp"

#include <boost/beast/core.hpp>
#include <boost/beast/version.hpp>

#include <iostream>
#include <sstream>

ChessNetwork::ChessNetwork()
  : resolver_(io_context_),
    websocket_(io_context_),
    pending_move_(false),
    is_connected_(false) {
}

ChessNetwork::~ChessNetwork() {
  disconnect();
}

bool ChessNetwork::connect(const std::string& ip_address, std::uint16_t port, const EventHandler& handle) {
  try {
    if (is_connected_) {
      return true;
    }

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
    if (is_connected_) {
        is_connected_ = false;

        beast::error_code ec;
        websocket_.next_layer().shutdown(tcp::socket::shutdown_both, ec);
        websocket_.next_layer().close(ec);
    }

    // Always join: the receive thread may have already cleared is_connected_
    // itself on a read error (server drop, match end, socket teardown). The
    // std::thread is still joinable and must be joined before destruction,
    // otherwise its destructor calls std::terminate().
    if (receive_thread_.joinable() &&
        receive_thread_.get_id() != std::this_thread::get_id()) {
        receive_thread_.join();
    }
}

bool ChessNetwork::send_move(const std::string from, const std::string to, char promotion) {
  if(!is_connected_ || has_pending_move()) {
    return false;
  }
  
  try {
    pending_move_ = true;

    std::string messageType = "MOVE";
    std::string message = from.substr(0, 2) + " " + to.substr(0, 2) + " " + promotion;
    std::cout << messageType << " " << message << std::endl;
    websocket_.write(net::buffer(messageType + " " + message));

    return true;
  } catch (const std::exception& e) {
    pending_move_ = false;
    std::cerr << "Failed to send move: " << e.what() << std::endl;
    return false;
  }
}

bool ChessNetwork::resign() {
  if (!is_connected_) {
    return false;
  }

  try {
    websocket_.write(net::buffer(std::string("RSGN")));
    return true;
  } catch (const std::exception& e) {
    std::cerr << "Failed to resign: " << e.what() << std::endl;
    return false;
  }
}

void ChessNetwork::receive_loop() {
  try {
    while (is_connected_) {
      beast::flat_buffer buffer;
      beast::error_code ec;
      websocket_.read(buffer, ec);

      if (ec) {
        is_connected_ = false;
        pending_move_ = false;

        Event e;
        e.type = DISCONNECTED;

        handler_(e);

        break;
      }

      std::string message = beast::buffers_to_string(buffer.data());
      std::string parsed_message;
      std::string messageType;

      Event e;
      MessageType type = OTHER;

      if (message.length() >= 4) {
        messageType = message.substr(0, 4);
        parsed_message = "";

        if (message.length() > 5) {
          parsed_message = message.substr(5);
        }

        if (messageType == "MOVE") {
          parse_move(parsed_message, e);
          parse_timers(parsed_message, e);
          type = OPPONENT_MOVE;
        } 
        else if (messageType == "ACPT") {
          parse_timers(parsed_message, e);
          type = MOVE_ACCEPTED;
          pending_move_ = false;
        }
        else if (messageType == "RJCT") {
          parse_timers(parsed_message, e);
          type = MOVE_REJECTED;
          pending_move_ = false;
        }
        else if (messageType == "INVL") {
          type = INVALID;
          pending_move_ = false;
        }
        else if (messageType == "WHTE") {
          parse_timers(parsed_message, e);
          type = WHITE;
        }
        else if (messageType == "BLCK") {
          parse_timers(parsed_message, e);
          type = BLACK;
        }
        else if (messageType == "ENDM") {
          e.reason = parsed_message;
          type = MATCH_ENDED;
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
  char promotion = '-';
  ss >> from >> to >> promotion;
  e.from = from;
  e.to = to;
  e.promotion = promotion;
}

void ChessNetwork::parse_timers(const std::string& message, Event& e) {
  std::stringstream ss(message);
  std::string from, to, promotion;

  if (ss >> from >> to >> promotion >> e.white_time >> e.black_time) {
    return;
  }

  ss.clear();
  ss.str(message);
  ss >> e.white_time >> e.black_time;
}

bool ChessNetwork::has_pending_move() const {
  return pending_move_;
}
