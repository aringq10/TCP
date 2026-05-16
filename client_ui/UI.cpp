#include <iostream>
#include "UI.hpp"

ChessBoardUI::ChessBoardUI(float tileSize) : m_tileSize(tileSize), m_selectedX(-1), m_selectedY(-1) {}

std::string ChessBoardUI::gridToNotation(int x, int y) const {
    char file = 'a' + x;
    char rank = '8' - y;
    return std::string{ file, rank };
}

void ChessBoardUI::calculateValidMoves(Board& board) {
    m_validMoves.clear();
    // Assuming your friend made isValidMove public in Board.h!
    for (int y = 0; y < 8; ++y) {
        for (int x = 0; x < 8; ++x) {
            if (board.isValidMove(m_selectedX, m_selectedY, x, y)) {
                m_validMoves.push_back({ x, y });
            }
        }
    }
}

void ChessBoardUI::handleEvent(const sf::Event& event, const sf::RenderWindow& window, Board& board, ChessNetwork& network) {
    if (const auto* mousePressed = event.getIf<sf::Event::MouseButtonPressed>()) {
        if (mousePressed->button == sf::Mouse::Button::Left) {

            int gridX = mousePressed->position.x / static_cast<int>(m_tileSize);
            int gridY = mousePressed->position.y / static_cast<int>(m_tileSize);

            if (gridX >= 0 && gridX < 8 && gridY >= 0 && gridY < 8) {

                // If we already have a piece selected, check if we clicked a valid move
                if (m_selectedX != -1 && m_selectedY != -1) {
                    bool isValidMoveClicked = false;
                    for (const auto& move : m_validMoves) {
                        if (move.x == gridX && move.y == gridY) {
                            isValidMoveClicked = true;
                            break;
                        }
                    }

                    if (isValidMoveClicked) {
                        std::string from = gridToNotation(m_selectedX, m_selectedY);
                        std::string to = gridToNotation(gridX, gridY);

                        // Tell your friend's logic to make the move
                        if (board.makeMove(from, to)) {
                            std::cout << "Local Move made: " << from << " to " << to << std::endl;

                            // ==========================================
                            // [SERVER CONNECTION: SENDING MOVES]
                            // When your backend teammates are ready, this is 
                            // where you will call network.sendMove(from, to);
                            // ==========================================

                            if (!network.send_move(from, to)) {
                                std::cout << "Could not send move to server" << std::endl;
                            }
                        }

                        // Deselect piece
                        m_selectedX = -1;
                        m_selectedY = -1;
                        m_validMoves.clear();
                        return;
                    }
                }

                // If we didn't complete a move, handle selecting a piece
                const Piece(&grid)[8][8] = board.getBoard();
                if (grid[gridY][gridX] != EMPTY) {
                    m_selectedX = gridX;
                    m_selectedY = gridY;
                    calculateValidMoves(board);
                }
                else {
                    m_selectedX = -1;
                    m_selectedY = -1;
                    m_validMoves.clear();
                }
            }
        }
    }
}

void ChessBoardUI::draw(sf::RenderTarget& target, const Board& board) const {
    sf::RectangleShape square(sf::Vector2f({ m_tileSize, m_tileSize }));

    for (int x = 0; x < 8; ++x) {
        for (int y = 0; y < 8; ++y) {
            square.setPosition({ x * m_tileSize, y * m_tileSize });
            if ((x + y) % 2 == 0) square.setFillColor(sf::Color(240, 217, 181));
            else square.setFillColor(sf::Color(181, 136, 99));
            target.draw(square);
        }
    }

    if (m_selectedX != -1 && m_selectedY != -1) {
        square.setPosition({ m_selectedX * m_tileSize, m_selectedY * m_tileSize });
        square.setFillColor(sf::Color(255, 255, 51, 150));
        target.draw(square);
    }
    for (const auto& move : m_validMoves) {
        square.setPosition({ static_cast<float>(move.x) * m_tileSize, static_cast<float>(move.y) * m_tileSize });
        square.setFillColor(sf::Color(50, 205, 50, 150));
        target.draw(square);
    }

    const Piece(&grid)[8][8] = board.getBoard();

    for (int y = 0; y < 8; ++y) {
        for (int x = 0; x < 8; ++x) {
            Piece p = grid[y][x];

            if (p != EMPTY) {
                sf::CircleShape pieceShape(m_tileSize / 2.5f);
                pieceShape.setPosition({ x * m_tileSize + (m_tileSize * 0.1f), y * m_tileSize + (m_tileSize * 0.1f) });

                switch (p) {
                case WHITE_PAWN:
                    pieceShape.setFillColor(sf::Color::White);
                    pieceShape.setOutlineThickness(2);
                    pieceShape.setOutlineColor(sf::Color::Black);
                    break;
                case BLACK_PAWN:
                    pieceShape.setFillColor(sf::Color(50, 50, 50));
                    pieceShape.setOutlineThickness(2);
                    pieceShape.setOutlineColor(sf::Color::Black);
                    break;
                case WHITE_ROOK:
                case WHITE_KNIGHT:
                case WHITE_BISHOP:
                case WHITE_QUEEN:
                case WHITE_KING:
                    pieceShape.setFillColor(sf::Color::Blue);
                    break;
                case BLACK_ROOK:
                case BLACK_KNIGHT:
                case BLACK_BISHOP:
                case BLACK_QUEEN:
                case BLACK_KING:
                    pieceShape.setFillColor(sf::Color::Red);
                    break;
                default:
                    break;
                }
                target.draw(pieceShape);
            }
        }
    }
}
