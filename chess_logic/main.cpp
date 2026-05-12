#include <iostream>
using namespace std;

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

class Board {
private:
    Piece board[8][8];

public:
    Board() {
        Piece startingBoard[8][8] = {
            {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
            {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
            {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
            {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
            {BLACK_KNIGHT, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
            {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
            {WHITE_PAWN, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
            {EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY}
        };

        for(int y = 0; y < 8; y++) {
            for(int x = 0; x < 8; x++) {
                board[y][x] = startingBoard[y][x];
            }
        }
    }
    bool isInsideBoard(int x, int y){
        
    return x >= 0 && x < 8 &&
           y >= 0 && y < 8;
    }

    bool isValidPawnMove(int fromX, int fromY, int toX, int toY, Piece piece) {

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
    bool isValidMove(int fromX, int fromY, int toX, int toY){

        if(!isInsideBoard(fromX, fromY)||!isInsideBoard(toX, toY))
        {
            return false;
        }

        Piece piece = board[fromY][fromX];

        if(piece == WHITE_PAWN || piece == BLACK_PAWN){
        return isValidPawnMove(fromX, fromY, toX, toY, piece);
        }
        return false;
    }

    void printBoard() { //kazka gal padaryt reikes cia
        cout << board[6][0] << endl;
    }
};

int main() {
    Board game;

    game.printBoard();

    cout << game.isValidMove(0, 6, 0, 4);

    return 0;
}