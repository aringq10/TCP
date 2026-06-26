#pragma once
#include "piece.h"

class Queen : public ChessPiece {
public:
    Queen(Color c) : ChessPiece(c, PieceType::QUEEN) {}
    bool isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const override;
};
