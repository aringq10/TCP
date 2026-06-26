#pragma once
#include "piece.h"

class Rook : public ChessPiece {
public:
    Rook(Color c) : ChessPiece(c, PieceType::ROOK) {}

    bool isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const override;
};
