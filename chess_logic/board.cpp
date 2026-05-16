#include "board.h"
#include <string>

Board::Board() {
    Piece startingBoard[8][8] = {
        {EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY},
        {BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN},
        {EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY},
        {EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY},
        {EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY},
        {EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY},
        {WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN},
        {EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY,      EMPTY}
    };

    whoseTurn = Color::WHITE;

    for(int y = 0; y < 8; y++)
        for(int x = 0; x < 8; x++)
            board[y][x] = startingBoard[y][x];
    
}

bool Board::isInsideBoard(int x, int y) {
    return x >= 0 && x < 8 &&
        y >= 0 && y < 8;
}

bool Board::isValidPawnMove(int fromX, int fromY, int toX, int toY, Piece piece) {
    Piece targetPiece = board[toY][toX];

    if (piece == EMPTY) {
        return false;
    }
    if (piece > 6 && targetPiece > 6) {
        return false;
    }
    if (piece < 7 && targetPiece < 7 && targetPiece != EMPTY) {
        return false;
    }

    int dir = 0;
    int startRow = 0;

    if(piece == WHITE_PAWN) {
        dir = -1;
        startRow = 6;
    }
    else if(piece == BLACK_PAWN) {
        dir = +1;
        startRow = 1;
    }

    int diffY = toY - fromY;
    int diffX = toX - fromX;

    if(board[toY][toX] != EMPTY) {
        if ((diffX == -1 || diffX == 1) && diffY == dir) {
            return true;
        }
        return false;
    }

    if(diffY == dir && diffX == 0) return true;

    if(fromY == startRow && diffY == 2 * dir && diffX == 0) {
        int midY = fromY + dir;

        if(board[midY][fromX] == EMPTY)
            return true;
    }
    return false;
}

bool Board::isValidMove(int fromX, int fromY, int toX, int toY) {
    if (myColor != whoseTurn) {
        return false;
    }

    if(!isInsideBoard(fromX, fromY)||!isInsideBoard(toX, toY)) {
        return false;
    }

    Piece piece = board[fromY][fromX];

    if (piece == EMPTY || (myColor == Color::WHITE && piece > 6) || (myColor == Color::BLACK && piece < 7)) {
        return false;
    }

    if(piece == WHITE_PAWN || piece == BLACK_PAWN) {
        return isValidPawnMove(fromX, fromY, toX, toY, piece);
    }
    return false;
}

bool Board::parseCoordinate(std::string coord, int &x, int &y) {
    if(coord.length() != 2) return false;

    char column = coord[0];
    char row = coord[1];

    if(column < 'a' || column > 'h') return false;

    if(row < '1' || row > '8') return false;

    x = column - 'a';
    y = 8 - (row - '0');

    return true;
}

bool Board::makeMove(std::string from, std::string to) {
    int fromX, fromY;
    int toX, toY;

    if(!parseCoordinate(from, fromX, fromY)) return false;

    if(!parseCoordinate(to, toX, toY)) return false;
    return makeMove(fromX, fromY, toX, toY);
}

bool Board::makeMove(int fromX, int fromY, int toX, int toY) {
    if(!isValidMove(fromX, fromY, toX, toY)){
        return false;
    }

    whoseTurn = myColor == Color::WHITE ? Color::BLACK : Color::WHITE;

    movedPiece = board[fromY][fromX];
    capturedPiece = board[toY][toX];

    lastFromX = fromX;
    lastFromY = fromY;
    lastToX = toX;
    lastToY = toY;
    hasMoveToUndo = true;

    board[toY][toX]=board[fromY][fromX];
    board[fromY][fromX]=EMPTY;
    return true;
}

bool Board::makeOppMove(std::string from, std::string to) {
    int fromX, fromY;
    int toX, toY;

    whoseTurn = myColor;

    if(!parseCoordinate(from, fromX, fromY)) return false;

    if(!parseCoordinate(to, toX, toY)) return false;

    return makeOppMove(fromX, fromY, toX, toY);
}

bool Board::makeOppMove(int fromX, int fromY, int toX, int toY) {

    board[toY][toX]=board[fromY][fromX];
    board[fromY][fromX]=EMPTY;
    return true;
}

void Board::setColor(Color c) {
    myColor = c;
}

bool Board::undoLastMove() {
    if(!hasMoveToUndo)
        return false;

    board[lastFromY][lastFromX] = movedPiece;
    board[lastToY][lastToX] = capturedPiece;

    hasMoveToUndo = false;

    return true;
}

const Piece (&Board::getBoard() const)[8][8] {
    return board;
}
