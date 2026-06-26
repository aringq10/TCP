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

    // Square a pawn may move onto to capture en passant (-1 if none),
    // i.e. the square skipped by the last double pawn push.
    int enPassantX, enPassantY;

    // Everything needed to reverse the last local move (used on server reject).
    struct Undo {
        bool valid = false;
        int fromX, fromY, toX, toY;
        ChessPiece* captured = nullptr; // captured piece (nullptr if none)
        int capturedX, capturedY;       // its square (differs from `to` on en passant)
        bool moverMovedBefore = false;  // mover's hasMoved() before the move
        bool isCastle = false;
        int rookFromX, rookFromY, rookToX, rookToY;
        bool rookMovedBefore = false;
        int enPassantXBefore, enPassantYBefore;
    };
    Undo undo;

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

    int getEnPassantTargetX() const { return enPassantX; }
    int getEnPassantTargetY() const { return enPassantY; }

    const Piece (&getBoard() const)[8][8];
};