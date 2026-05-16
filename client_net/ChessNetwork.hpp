
#include <boost/asio/io_context.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <boost/beast/websocket.hpp>

#include <cstdint>
#include <functional>
#include <string>
#include <thread>
#include <atomic>

namespace beast = boost::beast;
namespace websocket = beast::websocket;
namespace net = boost::asio;

enum MessageType {
    OPPONENT_MOVE,
    MOVE_ACCEPTED,
    MOVE_REJECTED,
    INVALID,
    WHITE,
    BLACK,
    OTHER
};

struct Event {
    MessageType type;
    std::string received_message;
    std::string from;
    std::string to;
};

using EventHandler = std::function<void(Event)>;

class ChessNetwork {
public:
  ChessNetwork();
  ~ChessNetwork();

  bool connect(const std::string& ip_address, std::uint16_t port, const EventHandler& handle);
  bool send_move(const char from[2], const char to[2]);
  void disconnect();

private:
  using tcp = boost::asio::ip::tcp;

  void receive_loop();
  void parse_move(const std::string& message, Event& e);
  bool has_pending_move() const;

  net::io_context io_context_;
  tcp::resolver resolver_;
  websocket::stream<tcp::socket> websocket_;
  std::atomic<bool> pending_move_;

  std::atomic<bool> is_connected_;
  std::thread receive_thread_;
  EventHandler handler_;
};
