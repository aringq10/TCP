const remote = "ws://192.168.88.14:6767/ws"

const lightColor = "#ffffff"
const darkColor = "#69923e"
const selectedColor = "red"

const boardSize = 400
const squareSize = boardSize / 8

const ASSET_PATH = "./assets/";

const layout = [
  "rnbqkbnr",
  "pppppppp",
  "........",
  "........",
  "........",
  "........",
  "PPPPPPPP",
  "RNBQKBNR"
];

const imageFiles = {
  K: "wk.png", Q: "wq.png", R: "wr.png",
  B: "wb.png", N: "wn.png", P: "wp.png",
  k: "bk.png", q: "bq.png", r: "br.png",
  b: "bb.png", n: "bn.png", p: "bp.png"
};

const game = {
  ongoing: false,
  ended: false,
  images: {},
  board: null,
  pendingMove: null,
  myColor: null
};

function start() {
  const boardEl = document.getElementById("board")
  const stateEl = document.getElementById("game-state")
  boardEl.setAttribute("width", boardSize)
  boardEl.setAttribute("height", boardSize)
  const ctx = boardEl.getContext("2d")

  stateEl.innerText = "Waiting for Match"

  drawBoard()
  drawPieces()

  const socket = new WebSocket(remote)
  socket.onmessage = (event) => {
    if (typeof event.data === "string") {
      if (event.data.length < 4) {
        console.log("received unknown message from remote:", event.data)
      }
      if (event.data === "WHTE") {
        game.ongoing = true
        game.myColor = "white"
        console.log("match started, you're", event.data)
        stateEl.innerText = `Match Ongoing, you're ${game.myColor}`
      } else if (event.data === "BLCK") {
        game.ongoing = true
        game.myColor = "black"
        console.log("match started, you're", event.data)
        stateEl.innerText = `Match Ongoing, you're ${game.myColor}`
      } else if (event.data === "RJCT") {
        game.pendingMove = null
        console.log("move rejected")
        window.alert("Move Rejected")
      } else if (event.data === "ACPT") {
        execMove(stringToCoords(game.pendingMove))
        drawBoard()
        drawPieces()
        game.pendingMove = null
        console.log("move accepted")
      } else if (event.data.startsWith("MOVE ")) {
        if (event.data.length < 10) {
          console.log("received unknown message from remote:", event.data)
        }
        execMove(stringToCoords(event.data.slice(5, 10)))
        drawBoard()
        drawPieces()
        console.log("received move from opponent:", event.data.slice(5, 10))
      }
    } else {
      console.log("received non-text message from remote:", event.data);
    }
  };
  socket.onclose = (event) => {
    let message = "Match Ended: "
    if (event.reason.startsWith("EOM ") && event.reason.length >= 5) {
      message += event.reason.slice(4)
    } else {
      message += "unrecognized reason"
    }
    game.ended = true
    console.log(message)
    stateEl.innerText = message
  };

  let selR = null, selC = null

  boardEl.addEventListener("click", function (e) {
    if (!game.ongoing || game.ended) {
      return
    }
    const x = e.clientX - this.offsetLeft - this.clientLeft
    const y = e.clientY - this.offsetTop - this.clientTop
    const r = Math.floor(y / squareSize)
    const c = Math.floor(x / squareSize)
    if (r >= 0 && r <= 7 && c >= 0 && c <= 7) {
      if (selR === null || selC === null) {
        const square = game.board[r][c]
        if (!square || square.color !== game.myColor) {
          return
        }
        selR = r
        selC = c
        ctx.fillStyle = selectedColor
        drawSquare(selR, selC)
        drawPieces()
      } else if (selR === r && selC === c) {
        ctx.fillStyle = squareColor(selR, selC)
        drawSquare(selR, selC)
        drawPieces()
        selR = selC = null
      } else {
        const square = game.board[r][c]
        if (square && square.color === game.myColor) {
          ctx.fillStyle = squareColor(selR, selC)
          drawSquare(selR, selC)
          selR = r
          selC = c
          ctx.fillStyle = selectedColor
          drawSquare(selR, selC)
          drawPieces()
        } else {
          const from = coordsToString(selR, selC),
                  to = coordsToString(r, c)
          ctx.fillStyle = squareColor(selR, selC)
          drawSquare(selR, selC)
          drawPieces()
          selR = selC = null
          game.pendingMove = `${from} ${to}`
          socket.send(`MOVE ${from} ${to}`)
        }
      }
    }
  })

  function drawBoard() {
    for (let row = 0; row < 8; row++) {
      for (let col = 0; col < 8; col++) {
        const x = col * squareSize
        const y = row * squareSize
        ctx.fillStyle = squareColor(row, col)
        ctx.fillRect(x, y, squareSize, squareSize)
      }
    }
  }
  function drawPieces() {
    for (let row = 0; row < 8; row++) {
      for (let col = 0; col < 8; col++) {
        const piece = game.board[row][col];

        if (piece) {
          ctx.drawImage(
            piece.img,
            col * squareSize,
            row * squareSize,
            squareSize,
            squareSize
          );
        }
      }
    }
  }
  function drawSquare(r, c) {
    ctx.fillRect(c * squareSize, r * squareSize, squareSize, squareSize)
  }
  function execMove({ fromRow, fromCol, toRow, toCol }) {
    if (
      fromRow < 0 || fromRow > 7 ||
      fromCol < 0 || fromCol > 7 ||
      toRow < 0 || toRow > 7 ||
      toCol < 0 || toCol > 7
    ) return;

    const piece = game.board[fromRow][fromCol];
    if (!piece) return;

    game.board[toRow][toCol] = piece;
    game.board[fromRow][fromCol] = null;
  }
}

function coordsToString(r, c) {
  return `${String.fromCharCode(97 + c)}${8 - r}`;
}

function stringToCoords(coords) {
  const [from, to] = coords.split(" ");

  function parse(square) {
    const col = square.charCodeAt(0) - 97;
    const row = 8 - Number(square[1]);

    return { row, col };
  }

  const fromCoords = parse(from);
  const toCoords = parse(to);

  return {
    fromRow: fromCoords.row,
    fromCol: fromCoords.col,
    toRow: toCoords.row,
    toCol: toCoords.col
  };
}

function squareColor(row, col) {
  return (row + col) % 2 ? darkColor : lightColor
}

async function init() {
  // preload sprites
  await Promise.all(
    Object.entries(imageFiles).map(([key, file]) => {
      return new Promise((resolve, reject) => {
        const img = new Image();
        img.onload = () => {
          game.images[key] = img;
          resolve();
        };
        img.onerror = reject;
        img.src = ASSET_PATH + file;
      });
    })
  );

  // setup board
  game.board = layout.map(row =>
    [...row].map(ch =>
      ch === "."
        ? null
        : {
            type: ch.toLowerCase(),
            color: ch === ch.toUpperCase() ? "white" : "black",
            img: game.images[ch]
          }
    )
  );

  start();
}

document.addEventListener("DOMContentLoaded", init)
