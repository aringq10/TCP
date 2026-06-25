#include "king.h"
#include "board.h"
#include <cmath>

bool King::isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const {
    int diffX = std::abs(toX - fromX);
    int diffY = std::abs(toY - fromY);

    // King moves exactly one square in any direction
    if (diffX > 1 || diffY > 1 || (diffX == 0 && diffY == 0)) return false;

    ChessPiece* target = board.getPieceAt(toX, toY);
    if (target != nullptr && target->getColor() == this->color) return false;
    return true;
}
