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
  disconnect();
}

bool ChessNetwork::connect(const std::string& ip_address, std::uint16_t port, const EventHandler& handle) {
  try {
    if (is_connected_) {
      return true;
    }
    Event e {OTHER};
    handle(e);

    const auto endpoints = resolver_.resolve(ip_address, std::to_string(port));
    const auto endpoint = net::connect(websocket_.next_layer(), endpoints);

    const std::string host_header = ip_address + ":" + std::to_string(endpoint.port());

    websocket_.handshake(host_header, "/ws");
    is_connected_ = true;

    std::cout << "Connected to WebSocket server at " << host_header << std::endl;

    handler_ = handle;
    
    receive_thread_ = std::thread([this]() {
      receive_loop();
    });


    /*
      while (true) {
        data = websocket_.read(); // Read incomming messages from ws conn
        Event e = { from data }   // Construct Event
        handle(e);                // Call client-passed handler function
      }
    */

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
    is_connected_ = false;
    websocket_.close(websocket::close_code::normal);
  } catch (const std::exception&) {
  }
}

void ChessNetwork::send_move(const char from[2], const char to[2]) {
  websocket_.write(net::buffer(std::string(from, 2) + std::string(to, 2)));
}

void ChessNetwork::receive_loop() {
  try {
    while (is_connected_) {
      beast::flat_buffer buffer;
      websocket_.read(buffer);

      std::string message = beast::buffers_to_string(buffer.data());

      std::cout << "Received message: " << message << std::endl;
    }
  } catch (const std::exception& e) {
    std::cerr << "Error in receive loop: " << e.what() << std::endl;
  }
}