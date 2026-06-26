#include "queen.h"
#include "board.h"
#include <cmath>

bool Queen::isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const {
    int diffX = toX - fromX;
    int diffY = toY - fromY;

    bool isDiagonal = (std::abs(diffX) == std::abs(diffY)) && diffX != 0;
    bool isStraight  = (diffX == 0) != (diffY == 0);

    if (!isDiagonal && !isStraight) return false;

    int stepX = (diffX == 0) ? 0 : (diffX > 0 ? 1 : -1);
    int stepY = (diffY == 0) ? 0 : (diffY > 0 ? 1 : -1);

    int x = fromX + stepX, y = fromY + stepY;
    while (x != toX || y != toY) {
        if (board.getPieceAt(x, y) != nullptr) return false;
        x += stepX;
        y += stepY;
    }

    ChessPiece* target = board.getPieceAt(toX, toY);
    if (target != nullptr && target->getColor() == this->color) return false;
    return true;
}
