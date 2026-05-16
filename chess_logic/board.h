#pragma once
#include <string>

enum Piece {
    EMPTY = 0,

    WHITE_PAWN = 1,
    WHITE_ROOK = 2,
    WHITE_KNIGHT = 3,
    WHITE_BISHOP = 4,
    WHITE_QUEEN = 5,
    WHITE_KING = 6,

    BLACK_PAWN = 7,
    BLACK_ROOK = 8,
    BLACK_KNIGHT = 9,
    BLACK_BISHOP = 10,
    BLACK_QUEEN = 11,
    BLACK_KING = 12
};

enum class Color {
    WHITE,
    BLACK
};

class Board {
private:
    Piece board[8][8];
    Color myColor;
    Color whoseTurn;

    int lastFromX;
    int lastFromY;
    int lastToX;
    int lastToY;

    Piece movedPiece;
    Piece capturedPiece;

    bool hasMoveToUndo;
    bool isInsideBoard(int x, int y);
    bool isValidPawnMove(int fromX, int fromY, int toX, int toY, Piece piece);
    bool parseCoordinate(std::string coord, int &x, int &y);

public:
    Board();

    bool isValidMove(int fromX, int fromY, int toX, int toY);
    bool makeMove(std::string from, std::string to);
    bool makeMove(int fromX, int fromY, int toX, int toY);
    bool makeOppMove(std::string from, std::string to);
    bool makeOppMove(int fromX, int fromY, int toX, int toY);
    void setColor(Color c);
    bool undoLastMove();
    const Piece (&getBoard() const)[8][8];
};
