#pragma once
#include "piece.h"

class Knight : public ChessPiece {
public:
    Knight(Color c) : ChessPiece(c, PieceType::KNIGHT) {}

    bool isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const override;
};