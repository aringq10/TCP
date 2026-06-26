#pragma once
#include <string>
#include "piece.h"
#include "pawn.h"

class Board {
private:
    ChessPiece* board[8][8]; 
    
    mutable Piece uiBoardCache[8][8];

    Color myColor;
    Color whoseTurn;

    int lastFromX, lastFromY, lastToX, lastToY;
    ChessPiece* capturedPieceHistory;
    bool hasMoveToUndo;

    bool isInsideBoard(int x, int y) const;
    bool parseCoordinate(std::string coord, int &x, int &y);

public:
    Board();
    ~Board(); 

    bool isValidMove(int fromX, int fromY, int toX, int toY);
    bool makeMove(std::string from, std::string to);
    bool makeMove(int fromX, int fromY, int toX, int toY);
    bool makeOppMove(std::string from, std::string to);
    bool makeOppMove(int fromX, int fromY, int toX, int toY);
    
    void setColor(Color c);
    bool undoLastMove();
    
    ChessPiece* getPieceAt(int x, int y) const; 

    const Piece (&getBoard() const)[8][8]; 
};