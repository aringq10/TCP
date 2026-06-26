#pragma once
#include <SFML/Graphics.hpp>
#include <string>
#include <vector>
#include <optional>
#include "../chess_logic/board.h"
#include "../client_net/ChessNetwork.hpp"

enum class GameState {
    MainMenu,
    Playing,
    Promotion,
    GameOver
};

enum class UIAction {
    None,
    JoinMatch,
    LeaveGame,
    Resign,
    CloseWindow
};

class Button {
public:
    sf::RectangleShape rect;
    std::optional<sf::Text> text;

    Button() = default;
    Button(const sf::Font& font, const std::string& str, sf::Vector2f size);

    void setCenterPosition(sf::Vector2f center);
    bool isClicked(sf::Vector2f mousePos) const;
    void draw(sf::RenderTarget& target) const;
};

class ChessBoardUI {
public:
    ChessBoardUI(float tileSize);

    void setFlipped(bool isFlipped);
    void setGameState(GameState state);
    GameState getGameState() const;
    void setGameOver(bool isWin);

    UIAction handleEvent(const sf::Event& event, const sf::RenderWindow& window, Board& board, ChessNetwork& network);

    void draw(sf::RenderTarget& target, const Board& board) const;

private:
    float m_tileSize;
    int m_selectedX;
    int m_selectedY;
    std::vector<sf::Vector2i> m_validMoves;

    sf::Font m_font;
    bool m_hasfont;

    sf::Texture m_bgTexture;
    sf::Sprite m_bgSprite;
    bool m_hasBackground;

    sf::Texture m_piece_textures[13];
    bool m_has_piece_texture[13] = { false };

    bool m_isFlipped;
    GameState m_currentState;
    bool m_isWin;

    mutable Button m_joinButton;
    mutable Button m_quitButton;
    mutable Button m_resignButton;
    mutable Button m_okButton;

    std::string gridToNotation(int x, int y) const;
    void calculateValidMoves(Board& board);

    sf::Vector2f getBoardOffset(const sf::RenderTarget& target) const;
    void updateLayout(sf::Vector2f viewSize) const;

    void drawMainMenu(sf::RenderTarget& target) const;
    void drawBoard(sf::RenderTarget& target, const Board& board) const;
    void drawPromotion(sf::RenderTarget& target) const;
    void drawGameOver(sf::RenderTarget& target) const;
};