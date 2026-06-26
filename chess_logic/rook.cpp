#include "rook.h"
#include "board.h"

bool Rook::isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const {
    bool sameRow = (fromY == toY);
    bool sameCol = (fromX == toX);
    if (!sameRow && !sameCol) return false;

    int stepX = 0, stepY = 0;
    if (sameRow) stepX = (toX > fromX) ? 1 : -1;
    else         stepY = (toY > fromY) ? 1 : -1;

    int x = fromX + stepX;
    int y = fromY + stepY;
    while (x != toX || y != toY) {
        if (board.getPieceAt(x, y) != nullptr) return false; // blocked
        x += stepX;
        y += stepY;
    }

    ChessPiece* target = board.getPieceAt(toX, toY);
    if (target != nullptr && target->getColor() == this->color) return false;

    return true;
}
