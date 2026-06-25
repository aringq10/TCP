#pragma once
#include "piece.h"

class King : public ChessPiece {
public:
    King(Color c) : ChessPiece(c, PieceType::KING) {}
    bool isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const override;
};
