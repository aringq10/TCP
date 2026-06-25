#pragma once
#include "piece.h"

class Pawn : public ChessPiece {
public:
    Pawn(Color c) : ChessPiece(c, PieceType::PAWN) {}

    bool isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const override;
};
