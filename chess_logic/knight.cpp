#include "knight.h"
#include "board.h"
#include <cmath>

bool Knight::isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const {
    int diffX = std::abs(toX - fromX);
    int diffY = std::abs(toY - fromY);

    // L shape
    bool isValidLShape = (diffX == 1 && diffY == 2) || (diffX == 2 && diffY == 1);
    if (!isValidLShape) {
        return false;
    }

    ChessPiece* target = board.getPieceAt(toX, toY);
    if (target != nullptr) {
        if (target->getColor() == this->color) {
            return false;
        }
    }
    return true;
}