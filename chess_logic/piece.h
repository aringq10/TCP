#pragma once
#include <vector>

enum Piece {
    EMPTY = 0,
    WHITE_PAWN = 1, WHITE_ROOK = 2, WHITE_KNIGHT = 3, WHITE_BISHOP = 4, WHITE_QUEEN = 5, WHITE_KING = 6,
    BLACK_PAWN = 7, BLACK_ROOK = 8, BLACK_KNIGHT = 9, BLACK_BISHOP = 10, BLACK_QUEEN = 11, BLACK_KING = 12
};

enum class Color { WHITE, BLACK };
enum class PieceType { PAWN, ROOK, KNIGHT, BISHOP, QUEEN, KING, EMPTY };

class Board; 

class ChessPiece {
protected:
    Color color;
    PieceType type;

public:
    ChessPiece(Color c, PieceType t) : color(c), type(t) {}
    virtual ~ChessPiece() = default;

    Color getColor() const { return color; }
    PieceType getType() const { return type; }

    Piece toEnum() const {
        if (type == PieceType::PAWN) {
            return (color == Color::WHITE) ? WHITE_PAWN : BLACK_PAWN;
        }
        if (type == PieceType::KNIGHT) {
            return (color == Color::WHITE) ? WHITE_KNIGHT : BLACK_KNIGHT;
        }
        if (type == PieceType::ROOK) {
            return (color == Color::WHITE) ? WHITE_ROOK : BLACK_ROOK;
        }
        if (type == PieceType::BISHOP) {
            return (color == Color::WHITE) ? WHITE_BISHOP : BLACK_BISHOP;
        }
        if (type == PieceType::QUEEN) {
            return (color == Color::WHITE) ? WHITE_QUEEN : BLACK_QUEEN;
        }
        if (type == PieceType::KING) {
            return (color == Color::WHITE) ? WHITE_KING : BLACK_KING;
        }
        return EMPTY;
    }

    virtual bool isValidMove(int fromX, int fromY, int toX, int toY, const Board& board) const = 0;
};