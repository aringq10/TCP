#include <SFML/Graphics.hpp>
#include <optional>
#include <iostream>
#include "client_ui/UI.hpp"
#include "chess_logic/board.h"
#include "client_net/ChessNetwork.hpp"

int main(int argc, char** argv) {
    if (argc < 2) {
        std::cout << "Usage: " << argv[0] << " server domain" << std::endl;
        return 1;
    }

    std::string address = argv[1];

    const float tileSize = 80.0f;

    sf::RenderWindow window(sf::VideoMode({ 600, 600 }), "TCP");
    window.setFramerateLimit(60);

    Board board;
    ChessNetwork network;
    ChessBoardUI ui(tileSize);

    Color myColor = Color::WHITE;

    auto handle_event = [&](Event e) {
        switch (e.type) {
        case OPPONENT_MOVE:
            std::cout << "OPPONENT MOVED: " << e.received_message << std::endl;
            board.makeOppMove(e.from, e.to);
            break;

        case MOVE_ACCEPTED:
            std::cout << "MOVE ACCEPTED" << std::endl;
            break;

        case MOVE_REJECTED:
            board.undoLastMove();
            std::cout << "MOVE REJECTED" << std::endl;
            break;

        case INVALID:
            std::cout << "INVALID MOVE" << std::endl;
            break;

        case WHITE:
            std::cout << "You are playing as WHITE" << std::endl;
            myColor = Color::WHITE;
            board.setColor(Color::WHITE);
            break;

        case BLACK:
            std::cout << "You are playing as BLACK" << std::endl;
            myColor = Color::BLACK;
            board.setColor(Color::BLACK);
            ui.setFlipped(true);
            break;

        case MATCH_ENDED: {
            std::cout << "Match has ended: " << e.reason << std::endl;
            bool isWin = (myColor == Color::WHITE && e.reason.rfind("1-0", 0) == 0) ||
                         (myColor == Color::BLACK && e.reason.rfind("0-1", 0) == 0);
            ui.setGameOver(isWin);
            break;
        }

        case DISCONNECTED:
            std::cout << "Connection to server closed" << std::endl;
            break;

        default:
            std::cout << "Unknown event type received" << std::endl;
            break;
        }
        };

    while (window.isOpen()) {
        while (const std::optional event = window.pollEvent()) {
            if (event->is<sf::Event::Closed>()) {
                window.close();
            }
            if (const auto* resized = event->getIf < sf::Event::Resized>()) {
                sf::FloatRect visibleArea({ 0.f,0.f }, { static_cast<float>(resized->size.x),static_cast<float>(resized->size.y) });
                window.setView(sf::View(visibleArea));
            }
            UIAction action = ui.handleEvent(*event, window, board, network);

            if (action == UIAction::JoinMatch) {
                if (!network.connect(address, 6767, handle_event)) {
                    std::cout << "Could not connect to server" << std::endl;
                    ui.setGameState(GameState::MainMenu);
                }
            }
            else if (action == UIAction::CloseWindow) {
                window.close();
            }
        }

        window.clear();
        ui.draw(window, board);
        window.display();
    }

    network.disconnect();
    return 0;
}
