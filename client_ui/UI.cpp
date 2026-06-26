#include <iostream>
#include "UI.hpp"

Button::Button(const sf::Font& font, const std::string& str, sf::Vector2f size) {
    rect.setSize(size);
    rect.setFillColor(sf::Color(50, 50, 50, 220));
    rect.setOutlineThickness(2.f);
    rect.setOutlineColor(sf::Color::White);

    text.emplace(font, str, 24);
    text->setFillColor(sf::Color::White);
}

void Button::setCenterPosition(sf::Vector2f center) {
    if (text) {
        sf::FloatRect textBounds = text->getLocalBounds();
        text->setOrigin({ textBounds.position.x + textBounds.size.x / 2.0f, textBounds.position.y + textBounds.size.y / 2.0f });
        text->setPosition(center);
    }

    rect.setOrigin({ rect.getSize().x / 2.0f, rect.getSize().y / 2.0f });
    rect.setPosition(center);
}

bool Button::isClicked(sf::Vector2f mousePos) const {
    return rect.getGlobalBounds().contains(mousePos);
}

void Button::draw(sf::RenderTarget& target) const {
    target.draw(rect);
    if (text) {
        target.draw(*text);
    }
}

ChessBoardUI::ChessBoardUI(float tileSize) : m_tileSize(tileSize), m_selectedX(-1), m_selectedY(-1),
m_isFlipped(false), m_currentState(GameState::MainMenu), m_isWin(false), m_bgSprite(m_bgTexture) {

    m_hasfont = m_font.openFromFile("../client_ui/font.ttf");

    if (m_hasfont) {
        m_joinButton = Button(m_font, "Join Match", { 200.f, 60.f });
        m_quitButton = Button(m_font, "Leave Game", { 200.f, 60.f });
        m_resignButton = Button(m_font, "Resign", { 150.f, 50.f });
        m_okButton = Button(m_font, "OK", { 150.f, 50.f });
    }

    m_hasBackground = m_bgTexture.loadFromFile("../client_ui/background.jpg");
    if (m_hasBackground) {
        m_bgSprite.setTexture(m_bgTexture, true);
    }

    std::string pieceFiles[13] = {
        "",
        "../client_ui/w_pawn.png", "../client_ui/w_rook.png", "../client_ui/w_knight.png", "../client_ui/w_bishop.png", "../client_ui/w_queen.png", "../client_ui/w_king.png",
        "../client_ui/b_pawn.png", "../client_ui/b_rook.png", "../client_ui/b_knight.png", "../client_ui/b_bishop.png", "../client_ui/b_queen.png", "../client_ui/b_king.png"
    };

    for (int i = 1; i <= 12; ++i) {
        m_has_piece_texture[i] = m_piece_textures[i].loadFromFile(pieceFiles[i]);
        if (m_has_piece_texture[i]) {
            m_piece_textures[i].setSmooth(true);
        }
        else {
            std::cout << "Warning: Could not load " << pieceFiles[i] << std::endl;
        }
    }
}
void ChessBoardUI::setFlipped(bool isFlipped) { m_isFlipped = isFlipped; }
void ChessBoardUI::setGameState(GameState state) { m_currentState = state; }
GameState ChessBoardUI::getGameState() const { return m_currentState; }
void ChessBoardUI::setGameOver(bool isWin) {
    m_isWin = isWin;
    m_currentState = GameState::GameOver;
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

sf::Vector2f ChessBoardUI::getBoardOffset(const sf::RenderTarget& target) const {
    float boardSize = m_tileSize * 8.0f;
    sf::Vector2f viewSize = target.getView().getSize();
    float offsetX = std::max(0.0f, (viewSize.x - 200.f - boardSize) / 2.0f);
    float offsetY = (viewSize.y - boardSize) / 2.0f;

    if (m_currentState == GameState::MainMenu) {
        offsetX = (viewSize.x - boardSize) / 2.0f; 
    }
    return { offsetX, offsetY };
}

void ChessBoardUI::updateLayout(sf::Vector2f viewSize) const {
    m_joinButton.setCenterPosition({ viewSize.x / 2.f, viewSize.y / 2.f - 40.f });
    m_quitButton.setCenterPosition({ viewSize.x / 2.f, viewSize.y / 2.f + 40.f });
    m_okButton.setCenterPosition({ viewSize.x / 2.f, viewSize.y / 2.f + 80.f });
    m_resignButton.setCenterPosition({ viewSize.x - 120.f, viewSize.y / 2.f });
}

UIAction ChessBoardUI::handleEvent(const sf::Event& event, const sf::RenderWindow& window, Board& board, ChessNetwork& network) {
    updateLayout(window.getView().getSize());

    if (const auto* mousePressed = event.getIf<sf::Event::MouseButtonPressed>()) {
        if (mousePressed->button == sf::Mouse::Button::Left) {
            sf::Vector2f worldPos = window.mapPixelToCoords(mousePressed->position);

            if (m_currentState == GameState::MainMenu) {
                if (m_joinButton.isClicked(worldPos)) {
                    m_currentState = GameState::Playing;
                    return UIAction::JoinMatch;
                }
                if (m_quitButton.isClicked(worldPos)) return UIAction::CloseWindow;
            }
            else if (m_currentState == GameState::GameOver) {
                if (m_okButton.isClicked(worldPos)) {
                    m_currentState = GameState::MainMenu;
                    return UIAction::LeaveGame;
                }
            }
            else if (m_currentState == GameState::Promotion) {
                m_currentState = GameState::Playing;
            }
            else if (m_currentState == GameState::Playing) {
                if (m_resignButton.isClicked(worldPos)) {
                    m_currentState = GameState::GameOver;
                    m_isWin = false;
                    return UIAction::Resign;
                }

                sf::Vector2f offset = getBoardOffset(window);
                float relativeX = worldPos.x - offset.x;
                float relativeY = worldPos.y - offset.y;
                if (relativeX < 0 || relativeY < 0) return UIAction::None;

                int gridX = static_cast<int>(relativeX / m_tileSize);
                int gridY = static_cast<int>(relativeY / m_tileSize);
                if (m_isFlipped) { gridX = 7 - gridX; gridY = 7 - gridY; }

                if (gridX >= 0 && gridX < 8 && gridY >= 0 && gridY < 8) {
                    if (m_selectedX != -1 && m_selectedY != -1) {
                        bool isValidMoveClicked = false;
                        for (const auto& move : m_validMoves) {
                            if (move.x == gridX && move.y == gridY) { isValidMoveClicked = true; break; }
                        }

                        if (isValidMoveClicked) {
                            std::string from = gridToNotation(m_selectedX, m_selectedY);
                            std::string to = gridToNotation(gridX, gridY);

                            if (board.makeMove(from, to)) {
                                std::cout << "Local Move: " << from << " to " << to << std::endl;
                                if (!network.send_move(from, to)) {
                                    std::cout << "Could not send move to server" << std::endl;
                                }

                                // Prototype pawn promotion
                                const Piece(&grid)[8][8] = board.getBoard();
                                Piece p = grid[gridY][gridX];
                                if ((p == WHITE_PAWN && gridY == 0) || (p == BLACK_PAWN && gridY == 7)) {
                                    m_currentState = GameState::Promotion;
                                }
                            }
                            m_selectedX = -1;
                            m_selectedY = -1;
                            m_validMoves.clear();
                            return UIAction::None;
                        }
                    }

                    const Piece(&grid)[8][8] = board.getBoard();
                    if (grid[gridY][gridX] != EMPTY) {
                        m_selectedX = gridX;
                        m_selectedY = gridY;
                        calculateValidMoves(board);
                    }
                    else {
                        m_selectedX = -1; m_selectedY = -1; m_validMoves.clear();
                    }
                }
            }
        }
    }
    return UIAction::None;
}

void ChessBoardUI::drawMainMenu(sf::RenderTarget& target) const {
    sf::Vector2f viewSize = target.getView().getSize();

    if (m_hasfont) {
        sf::Text title(m_font);
        title.setString("TCP");
        title.setCharacterSize(50);
        title.setFillColor(sf::Color::White);
        sf::FloatRect bounds = title.getLocalBounds();
        title.setOrigin({ bounds.position.x + bounds.size.x / 2.f, bounds.position.y + bounds.size.y / 2.f });
        title.setPosition({ viewSize.x / 2.f, viewSize.y / 4.f });
        target.draw(title);
    }
    m_joinButton.draw(target);
    m_quitButton.draw(target);
}

void ChessBoardUI::drawGameOver(sf::RenderTarget& target) const {
    sf::Vector2f viewSize = target.getView().getSize();

    sf::RectangleShape overlay({ viewSize.x, viewSize.y });
    overlay.setFillColor(sf::Color(0, 0, 0, 180));
    target.draw(overlay);

    if (m_hasfont) {
        sf::Text text(m_font);
        text.setString(m_isWin ? "YOU WIN!" : "YOU LOSE!");
        text.setCharacterSize(60);
        text.setFillColor(m_isWin ? sf::Color::Green : sf::Color::Red);
        sf::FloatRect bounds = text.getLocalBounds();
        text.setOrigin({ bounds.position.x + bounds.size.x / 2.f, bounds.position.y + bounds.size.y / 2.f });
        text.setPosition({ viewSize.x / 2.f, viewSize.y / 2.f - 50.f });
        target.draw(text);
    }
    m_okButton.draw(target);
}

void ChessBoardUI::drawPromotion(sf::RenderTarget& target) const {
    sf::Vector2f viewSize = target.getView().getSize();

    sf::RectangleShape overlay({ viewSize.x, viewSize.y });
    overlay.setFillColor(sf::Color(0, 0, 0, 150));
    target.draw(overlay);

    sf::RectangleShape popup({ 400.f, 150.f });
    popup.setFillColor(sf::Color(220, 220, 220));
    popup.setOrigin({ 200.f, 75.f });
    popup.setPosition({ viewSize.x / 2.f, viewSize.y / 2.f });
    target.draw(popup);

    if (m_hasfont) {
        sf::Text text(m_font);
        text.setString("Promote Pawn");
        text.setCharacterSize(24);
        text.setFillColor(sf::Color::Black);
        sf::FloatRect bounds = text.getLocalBounds();
        text.setOrigin({ bounds.position.x + bounds.size.x / 2.f, bounds.position.y + bounds.size.y / 2.f });
        text.setPosition({ viewSize.x / 2.f, viewSize.y / 2.f - 40.f });
        target.draw(text);
    }

    for (int i = 0; i < 4; i++) {
        sf::RectangleShape btn({ 60.f, 60.f });
        btn.setOrigin({ 30.f, 30.f });
        btn.setFillColor(sf::Color(100, 100, 100));
        btn.setOutlineThickness(2.f);
        btn.setOutlineColor(sf::Color::Black);
        btn.setPosition({ viewSize.x / 2.f - 120.f + i * 80.f, viewSize.y / 2.f + 20.f });
        target.draw(btn);
    }
}

void ChessBoardUI::drawBoard(sf::RenderTarget& target, const Board& board) const {
    sf::RenderStates states;
    sf::Vector2f offset = getBoardOffset(target);
    states.transform.translate({ offset.x, offset.y });

    sf::RectangleShape square(sf::Vector2f({ m_tileSize, m_tileSize }));

    for (int x = 0; x < 8; ++x) {
        for (int y = 0; y < 8; ++y) {
            int visX = m_isFlipped ? 7 - x : x;
            int visY = m_isFlipped ? 7 - y : y;
            square.setPosition({ visX * m_tileSize, visY * m_tileSize });
            if ((x + y) % 2 == 0) square.setFillColor(sf::Color(240, 217, 181));
            else square.setFillColor(sf::Color(181, 136, 99));
            target.draw(square, states);
        }
    }

    if (m_selectedX != -1 && m_selectedY != -1) {
        int visX = m_isFlipped ? 7 - m_selectedX : m_selectedX;
        int visY = m_isFlipped ? 7 - m_selectedY : m_selectedY;
        square.setPosition({ visX * m_tileSize, visY * m_tileSize });
        square.setFillColor(sf::Color(255, 255, 51, 150));
        target.draw(square, states);
    }
    for (const auto& move : m_validMoves) {
        int visX = m_isFlipped ? 7 - move.x : move.x;
        int visY = m_isFlipped ? 7 - move.y : move.y;
        square.setPosition({ static_cast<float>(visX) * m_tileSize, static_cast<float>(visY) * m_tileSize });
        square.setFillColor(sf::Color(50, 205, 50, 150));
        target.draw(square, states);
    }

    const Piece(&grid)[8][8] = board.getBoard();
    for (int y = 0; y < 8; ++y) {
        for (int x = 0; x < 8; ++x) {
            Piece p = grid[y][x];

            if (p != EMPTY) {
                int visX = m_isFlipped ? 7 - x : x;
                int visY = m_isFlipped ? 7 - y : y;

                if (p >= 1 && p <= 12 && m_has_piece_texture[p]) {
                    sf::Sprite pieceSprite(m_piece_textures[p]);

                    sf::Vector2u texSize = m_piece_textures[p].getSize();
                    pieceSprite.setScale({ m_tileSize / static_cast<float>(texSize.x), m_tileSize / static_cast<float>(texSize.y) });

                    pieceSprite.setPosition({ visX * m_tileSize, visY * m_tileSize });
                    target.draw(pieceSprite, states);
                }
                else {
                    sf::CircleShape pieceShape(m_tileSize / 2.5f);
                    pieceShape.setPosition({ visX * m_tileSize + (m_tileSize * 0.1f), visY * m_tileSize + (m_tileSize * 0.1f) });

                    switch (p) {
                    case WHITE_PAWN: pieceShape.setFillColor(sf::Color::White); pieceShape.setOutlineThickness(2); pieceShape.setOutlineColor(sf::Color::Black); break;
                    case BLACK_PAWN: pieceShape.setFillColor(sf::Color(50, 50, 50)); pieceShape.setOutlineThickness(2); pieceShape.setOutlineColor(sf::Color::Black); break;
                    case WHITE_ROOK: case WHITE_KNIGHT: case WHITE_BISHOP: case WHITE_QUEEN: case WHITE_KING:
                        pieceShape.setFillColor(sf::Color::Blue); break;
                    case BLACK_ROOK: case BLACK_KNIGHT: case BLACK_BISHOP: case BLACK_QUEEN: case BLACK_KING:
                        pieceShape.setFillColor(sf::Color::Red); break;
                    default: break;
                    }
                    target.draw(pieceShape, states);
                }
            }
        }
    }
}

void ChessBoardUI::draw(sf::RenderTarget& target, const Board& board) const {
    updateLayout(target.getView().getSize());

    if (m_hasBackground) {
        sf::Vector2f targetSize = target.getView().getSize();
        sf::Vector2f textureSize(m_bgTexture.getSize());
        const_cast<sf::Sprite&>(m_bgSprite).setScale({ targetSize.x / textureSize.x, targetSize.y / textureSize.y });
        target.draw(m_bgSprite);
    }

    if (m_currentState == GameState::MainMenu) {
        drawMainMenu(target);
    }
    else {
        drawBoard(target, board);
        m_resignButton.draw(target);

        if (m_currentState == GameState::Promotion) {
            drawPromotion(target);
        }
        else if (m_currentState == GameState::GameOver) {
            drawGameOver(target);
        }
    }
}