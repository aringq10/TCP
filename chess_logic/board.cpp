#include "board.h"
#include "knight.h"
#include "pawn.h"
#include "rook.h"
#include "bishop.h"
#include "queen.h"
#include "king.h"

Board::Board() : capturedPieceHistory(nullptr), hasMoveToUndo(false), whoseTurn(Color::WHITE), myColor(Color::WHITE) {
    // empty
    for (int y = 0; y < 8; y++) {
        for (int x = 0; x < 8; x++) {
            board[y][x] = nullptr;
        }
    }

    // pawns
    for (int x = 0; x < 8; x++) {
        board[1][x] = new Pawn(Color::BLACK);
    }
    for (int x = 0; x < 8; x++) {
        board[6][x] = new Pawn(Color::WHITE);
    }
    
    // knights
    board[7][1] = new Knight(Color::WHITE); // b1
    board[7][6] = new Knight(Color::WHITE); // g1
    board[0][1] = new Knight(Color::BLACK); // b8
    board[0][6] = new Knight(Color::BLACK); // g8

    // rooks
    board[0][0] = new Rook(Color::BLACK); // a8
    board[0][7] = new Rook(Color::BLACK); // h8
    board[7][0] = new Rook(Color::WHITE); // a1
    board[7][7] = new Rook(Color::WHITE); // h1

    // bishops
    board[0][2] = new Bishop(Color::BLACK); // c8
    board[0][5] = new Bishop(Color::BLACK); // f8
    board[7][2] = new Bishop(Color::WHITE); // c1
    board[7][5] = new Bishop(Color::WHITE); // f1

    // queens
    board[0][3] = new Queen(Color::BLACK);  // d8
    board[7][3] = new Queen(Color::WHITE);  // d1

    // kings
    board[0][4] = new King(Color::BLACK);   // e8
    board[7][4] = new King(Color::WHITE);   // e1
}

Board::~Board() {
    for (int y = 0; y < 8; y++) {
        for (int x = 0; x < 8; x++) {
            delete board[y][x];
        }
    }
    delete capturedPieceHistory;
}

// into enum
const Piece (&Board::getBoard() const)[8][8] {
    for (int y = 0; y < 8; ++y) {
        for (int x = 0; x < 8; ++x) {
            if (board[y][x] == nullptr) {
                uiBoardCache[y][x] = EMPTY;
            } else {
                uiBoardCache[y][x] = board[y][x]->toEnum();
            }
        }
    }
    return uiBoardCache;
}

ChessPiece* Board::getPieceAt(int x, int y) const {
    if (!isInsideBoard(x, y)) return nullptr;
    return board[y][x];
}

bool Board::isInsideBoard(int x, int y) const {
    return x >= 0 && x < 8 && y >= 0 && y < 8;
}

bool Board::isValidMove(int fromX, int fromY, int toX, int toY) {
    //out of bounds safety
    if (whoseTurn != myColor) return false; 
    if (!isInsideBoard(fromX, fromY) || !isInsideBoard(toX, toY)) return false;

    ChessPiece* movingPiece = board[fromY][fromX];
    if (movingPiece == nullptr || movingPiece->getColor() != whoseTurn) return false;

    return movingPiece->isValidMove(fromX, fromY, toX, toY, *this);
}

bool Board::makeMove(int fromX, int fromY, int toX, int toY) {
    if (!isValidMove(fromX, fromY, toX, toY)) return false;

    if (capturedPieceHistory) {
        delete capturedPieceHistory;
        capturedPieceHistory = nullptr;
    }

    capturedPieceHistory = board[toY][toX]; 
    lastFromX = fromX; lastFromY = fromY;
    lastToX = toX; lastToY = toY;
    hasMoveToUndo = true;

    board[toY][toX] = board[fromY][fromX];
    board[fromY][fromX] = nullptr; 
    
    // alternate turn
    whoseTurn = (whoseTurn == Color::WHITE) ? Color::BLACK : Color::WHITE;
    return true;
}

bool Board::makeOppMove(int fromX, int fromY, int toX, int toY) {
    if (!isInsideBoard(fromX, fromY) || !isInsideBoard(toX, toY)) return false;
    
    delete board[toY][toX]; 

    board[toY][toX] = board[fromY][fromX];
    board[fromY][fromX] = nullptr;

    whoseTurn = myColor; // reset turn focus back to player
    return true;
}

bool Board::undoLastMove() {
    if (!hasMoveToUndo) return false;

    board[lastFromY][lastFromX] = board[lastToY][lastToX];
    board[lastToY][lastToX] = capturedPieceHistory; 

    capturedPieceHistory = nullptr; 
    hasMoveToUndo = false;
    whoseTurn = (whoseTurn == Color::WHITE) ? Color::BLACK : Color::WHITE;
    return true;
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
    int fromX, fromY, toX, toY;
    if(!parseCoordinate(from, fromX, fromY)) return false;
    if(!parseCoordinate(to, toX, toY)) return false;
    return makeMove(fromX, fromY, toX, toY);
}

bool Board::makeOppMove(std::string from, std::string to) {
    int fromX, fromY, toX, toY;
    if(!parseCoordinate(from, fromX, fromY)) return false;
    if(!parseCoordinate(to, toX, toY)) return false;
    return makeOppMove(fromX, fromY, toX, toY);
}

void Board::setColor(Color c) { 
    myColor = c; 
    whoseTurn = Color::WHITE; // start with white
}