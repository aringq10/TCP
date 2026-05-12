
#include <boost/asio/io_context.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <boost/beast/websocket.hpp>

#include <cstdint>
#include <string>

namespace beast = boost::beast;
namespace websocket = beast::websocket;
namespace net = boost::asio;

enum MessageType {
    MOVE_OK
};

struct Event {
    MessageType type;
    // ...
};

using EventHandler = std::function<void(Event)>;

class ChessNetwork
{
public:
  ChessNetwork();
  ~ChessNetwork();

  bool connect(const std::string& ip_address, std::uint16_t port, const EventHandler& handle);
  void send_move(const char from[2], const char to[2]);
  void disconnect();

private:
  using tcp = boost::asio::ip::tcp;

  net::io_context io_context_;
  tcp::resolver resolver_;
  websocket::stream<tcp::socket> websocket_;
  bool is_connected_;
};
