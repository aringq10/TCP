#pragma once
#include <SFML/Graphics.hpp>
#include <string>
#include <vector>
#include "../chess_logic/board.h"
#include "../client_net/ChessNetwork.hpp"

class ChessBoardUI {
public:
    ChessBoardUI(float tileSize);

    void handleEvent(const sf::Event& event, const sf::RenderWindow& window, Board& board, ChessNetwork& network);

    void draw(sf::RenderTarget& target, const Board& board) const;

private:
    float m_tileSize;
    int m_selectedX;
    int m_selectedY;
    std::vector<sf::Vector2i> m_validMoves;

    sf::Font m_font;
    bool m_hasfont;

    std::string gridToNotation(int x, int y) const;
    void calculateValidMoves(Board& board);

    sf::Vector2f getBoardOffset(const sf::RenderTarget& target) const;
};
