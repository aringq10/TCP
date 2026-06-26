#include "king.h"
#include "board.h"
#include <cmath>

bool King::isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const {
    int dX = toX - fromX;
    int diffX = std::abs(dX);
    int diffY = std::abs(toY - fromY);

    // Castling: the king steps two squares sideways along its own rank.
    // Legality wrt. check / passing through attacked squares is enforced by
    // the authoritative server; here we only verify the geometric conditions.
    if (diffY == 0 && diffX == 2) {
        if (this->hasMoved()) return false;

        bool kingside = dX > 0;
        int rookX = kingside ? 7 : 0;
        ChessPiece* rook = board.getPieceAt(rookX, fromY);
        if (rook == nullptr || rook->getType() != PieceType::ROOK ||
            rook->getColor() != this->color || rook->hasMoved()) {
            return false;
        }

        // Every square between the king and the rook must be empty.
        int step = kingside ? 1 : -1;
        for (int x = fromX + step; x != rookX; x += step) {
            if (board.getPieceAt(x, fromY) != nullptr) return false;
        }
        return true;
    }

    // King moves exactly one square in any direction
    if (diffX > 1 || diffY > 1 || (diffX == 0 && diffY == 0)) return false;

    ChessPiece* target = board.getPieceAt(toX, toY);
    if (target != nullptr && target->getColor() == this->color) return false;
    return true;
}
