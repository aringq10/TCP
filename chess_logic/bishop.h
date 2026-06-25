#pragma once
#include "piece.h"

class Bishop : public ChessPiece {
public:
    Bishop(Color c) : ChessPiece(c, PieceType::BISHOP) {}
    bool isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const override;
};
