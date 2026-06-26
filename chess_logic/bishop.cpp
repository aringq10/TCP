#include "bishop.h"
#include "board.h"
#include <cmath>

bool Bishop::isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const {
    int diffX = toX - fromX;
    int diffY = toY - fromY;

    // Must move diagonally: |dx| == |dy| and non-zero
    if (std::abs(diffX) != std::abs(diffY) || diffX == 0) return false;

    int stepX = (diffX > 0) ? 1 : -1;
    int stepY = (diffY > 0) ? 1 : -1;

    int x = fromX + stepX, y = fromY + stepY;
    while (x != toX || y != toY) {
        if (board.getPieceAt(x, y) != nullptr) return false; // blocked
        x += stepX;
        y += stepY;
    }

    ChessPiece* target = board.getPieceAt(toX, toY);
    if (target != nullptr && target->getColor() == this->color) return false;
    return true;
}
