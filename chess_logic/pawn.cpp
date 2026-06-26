#include "pawn.h"
#include "board.h"
#include <cmath>

bool Pawn::isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const {
    // White pawns start on row 6 and move toward row 0 (decreasing y).
    // Black pawns start on row 1 and move toward row 7 (increasing y).
    int direction = (color == Color::WHITE) ? -1 : 1;
    int startRow = (color == Color::WHITE) ? 6 : 1;

    int diffX = toX - fromX;
    int diffY = toY - fromY;

    ChessPiece* target = board.getPieceAt(toX, toY);

    if (diffX == 0) {
        // Single move
        if (diffY == direction && target == nullptr) {
            return true;
        }

        // Double move
        if (fromY == startRow && diffY == 2 * direction && target == nullptr) {
            ChessPiece* squareInBetween = board.getPieceAt(fromX, fromY + direction);
            if (squareInBetween == nullptr) {
                return true;
            }
        }

        return false;
    }

    // Diagonal capture
    if (std::abs(diffX) == 1 && diffY == direction) {
        // Normal capture
        if (target != nullptr && target->getColor() != this->color) {
            return true;
        }
        if (target == nullptr && toX == board.getEnPassantTargetX() && toY == board.getEnPassantTargetY()) {
            return true;
        }
    }

    return false;
}
