#include "ChessNetwork.hpp"

#include <boost/beast/core.hpp>
#include <boost/beast/version.hpp>

#include <iostream>

ChessNetwork::ChessNetwork()
  : resolver_(io_context_),
    websocket_(io_context_),
    is_connected_(false) {
}

ChessNetwork::~ChessNetwork() {
  if (!is_connected_) {
    return;
  }

  try {
    websocket_.close(websocket::close_code::normal);
  } catch (const std::exception&) {
  }
}

bool ChessNetwork::connect(const std::string& ip_address, std::uint16_t port) {
  try {
    if (is_connected_) {
      websocket_.close(websocket::close_code::normal);
      is_connected_ = false;
    }

    const auto endpoints = resolver_.resolve(ip_address, std::to_string(port));
    const auto endpoint = net::connect(websocket_.next_layer(), endpoints);

    const std::string host_header = ip_address + ":" + std::to_string(endpoint.port());

    websocket_.handshake(host_header, "/ws");
    is_connected_ = true;

    std::cout << "Connected to WebSocket server at " << host_header << std::endl;
    return true;
  } catch (const std::exception& exception) {
    std::cerr << "WebSocket connection failed: " << exception.what() << std::endl;
    return false;
  }
}
