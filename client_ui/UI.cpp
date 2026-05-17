#include <iostream>
#include "UI.hpp"

ChessBoardUI::ChessBoardUI(float tileSize) : m_tileSize(tileSize), m_selectedX(-1), m_selectedY(-1), m_isFlipped(false), m_bgSprite(m_bgTexture) {
    m_hasfont = m_font.openFromFile("font.ttf");

    m_hasBackground = m_bgTexture.loadFromFile("background.jpg");
    if (m_hasBackground) {
        m_bgSprite.setTexture(m_bgTexture);
    }
}

void ChessBoardUI::setFlipped(bool isFlipped) {
    m_isFlipped = isFlipped;
}

std::string ChessBoardUI::gridToNotation(int x, int y) const {
    char file = 'a' + x;
    char rank = '8' - y;
    return std::string{ file, rank };
}

void ChessBoardUI::calculateValidMoves(Board& board) {
    m_validMoves.clear();
    for (int y = 0; y < 8; ++y) {
        for (int x = 0; x < 8; ++x) {
            if (board.isValidMove(m_selectedX, m_selectedY, x, y)) {
                m_validMoves.push_back({ x, y });
            }
        }
    }
}

//Render board in the middle
sf::Vector2f ChessBoardUI::getBoardOffset(const sf::RenderTarget& target) const {
    float boardSize = m_tileSize * 8.0f;
    sf::Vector2f viewSize = target.getView().getSize();
    float offsetX = (viewSize.x - boardSize) / 2.0f;
    float offsetY = (viewSize.y - boardSize) / 2.0f;
    return { offsetX,offsetY };
}

void ChessBoardUI::handleEvent(const sf::Event& event, const sf::RenderWindow& window, Board& board, ChessNetwork& network) {
    if (const auto* mousePressed = event.getIf<sf::Event::MouseButtonPressed>()) {
        if (mousePressed->button == sf::Mouse::Button::Left) {
            sf::Vector2f worldPos = window.mapPixelToCoords(mousePressed->position);

            sf::Vector2f offset = getBoardOffset(window);
            float relativeX = worldPos.x - offset.x;
            float relativeY = worldPos.y - offset.y;

            if (relativeX < 0 || relativeY < 0) return;

            int gridX = static_cast<int>(relativeX / m_tileSize);
            int gridY = static_cast<int>(relativeY / m_tileSize);

            if (m_isFlipped) {
                gridX = 7 - gridX;
                gridY = 7 - gridY;
            }

            if (gridX >= 0 && gridX < 8 && gridY >= 0 && gridY < 8) {

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

                        if (board.makeMove(from, to)) {
                            std::cout << "Local Move made: " << from << " to " << to << std::endl;

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
    if (m_hasBackground) {
        //stretch to fit
        sf::Vector2f targetSize = target.getView().getSize();
        sf::Vector2f textureSize(m_bgTexture.getSize());
        const_cast<sf::Sprite&>(m_bgSprite).setScale({ targetSize.x / textureSize.x, targetSize.y / textureSize.y });
        target.draw(m_bgSprite);
    }
    sf::RenderStates states;
    sf::Vector2f offset = getBoardOffset(target);
    states.transform.translate({ offset.x,offset.y });

    sf::RectangleShape square(sf::Vector2f({ m_tileSize, m_tileSize }));

    for (int x = 0; x < 8; ++x) {
        for (int y = 0; y < 8; ++y) {
            int visX = m_isFlipped ? 7 - x : x;
            int visY = m_isFlipped ? 7 - y : y;
            square.setPosition({ visX * m_tileSize, visY * m_tileSize });
            if ((x + y) % 2 == 0) square.setFillColor(sf::Color(240, 217, 181));
            else square.setFillColor(sf::Color(181, 136, 99));
            target.draw(square,states);
        }
    }

    if (m_selectedX != -1 && m_selectedY != -1) {
        int visX = m_isFlipped ? 7 - m_selectedX : m_selectedX;
        int visY = m_isFlipped ? 7 - m_selectedY : m_selectedY;
        square.setPosition({ visX * m_tileSize, visY * m_tileSize });
        square.setFillColor(sf::Color(255, 255, 51, 150));
        target.draw(square,states);
    }
    for (const auto& move : m_validMoves) {
        int visX = m_isFlipped ? 7 - move.x : move.x;
        int visY = m_isFlipped ? 7 - move.y : move.y;
        square.setPosition({ static_cast<float>(visX) * m_tileSize, static_cast<float>(visY) * m_tileSize });
        square.setFillColor(sf::Color(50, 205, 50, 150));
        target.draw(square,states);
    }

    const Piece(&grid)[8][8] = board.getBoard();

    for (int y = 0; y < 8; ++y) {
        for (int x = 0; x < 8; ++x) {
            Piece p = grid[y][x];

            if (p != EMPTY) {
                int visX = m_isFlipped ? 7 - x : x;
                int visY = m_isFlipped ? 7 - y : y;
                sf::CircleShape pieceShape(m_tileSize / 2.5f);
                pieceShape.setPosition({ visX * m_tileSize + (m_tileSize * 0.1f), visY * m_tileSize + (m_tileSize * 0.1f) });

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
                target.draw(pieceShape,states);
            }
        }
    }

    if (m_hasfont && target.getView().getSize().x > (m_tileSize * 8.0f) + 200.f) {
        sf::Text leftText(m_font);
        leftText.setString("ALIO");
        leftText.setCharacterSize(20);
        leftText.setFillColor(sf::Color::White);
        leftText.setPosition({ 20.f,20.f });
        target.draw(leftText);

        sf::Text leftText(m_font);
        leftText.setString("VALIO");
        leftText.setCharacterSize(20);
        leftText.setFillColor(sf::Color::White);
        leftText.setPosition({ target.getView().getSize().x - 180.f,20.f });
        target.draw(leftText);
    }
}
