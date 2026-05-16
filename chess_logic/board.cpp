#include "board.h"
#include <string>

// OPPONENT_MOVE
// makeMove(opponent's move)

Board::Board() {
    Piece startingBoard[8][8] = {
        {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
        {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
        {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
        {WHITE_PAWN, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
        {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
        {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
        {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
        {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY}
    };

    for(int y = 0; y < 8; y++)
        for(int x = 0; x < 8; x++)
            board[y][x] = startingBoard[y][x];
}
bool Board::isInsideBoard(int x, int y){
            
        return x >= 0 && x < 8 &&
            y >= 0 && y < 8;
}
bool Board::isValidPawnMove(int fromX, int fromY, int toX, int toY, Piece piece) {
    if(fromX != toX)
        return false;

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

    if(board[toY][toX] != EMPTY)
        return false;

    int diff = toY - fromY;

    if(diff == dir) return true;

    if(fromY == startRow && diff == 2 * dir) {
        int midY = fromY + dir;

        if(board[midY][fromX] == EMPTY)
            return true;
    }
    return false;
}
bool Board::isValidMove(int fromX, int fromY, int toX, int toY){

        if(!isInsideBoard(fromX, fromY)||!isInsideBoard(toX, toY))
        {
            return false;
        }

        Piece piece = board[fromY][fromX];

        if (piece == EMPTY || (myColor == WHITE && piece > 6) || (myColor == BLACK && piece < 7)) {
            return false;
        }
        if(piece == WHITE_PAWN || piece == BLACK_PAWN){
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
bool Board::makeMove(int fromX, int fromY, int toX, int toY){

            if(!isValidMove(fromX, fromY, toX, toY)){
                return false;
            }
            board[toY][toX]=board[fromY][fromX];
            board[fromY][fromX]=EMPTY;
            return true;
}
bool Board::makeOppMove(std::string from, std::string to) {
        int fromX, fromY;
        int toX, toY;

        if(!parseCoordinate(from, fromX, fromY)) return false;

        if(!parseCoordinate(to, toX, toY)) return false;
        return makeOppMove(fromX, fromY, toX, toY);
}
bool Board::makeOppMove(int fromX, int fromY, int toX, int toY){
            board[toY][toX]=board[fromY][fromX];
            board[fromY][fromX]=EMPTY;
            return true;
}
void Board::setColor(Color c) {
    myColor = c;
}
const Piece (&Board::getBoard() const)[8][8] { //getteris
    return board;
}
