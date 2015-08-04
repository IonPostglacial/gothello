package othello

import "bytes"

const BOARD_SIZE = 8

type Color int

const (
	BLACK = Color(iota)
	WHITE
	EMPTY
)

func (color Color) Opponent() Color {
	return (color + 1) % EMPTY
}

var boardRun [][]Position

type Board struct {
	cells [2][BOARD_SIZE * BOARD_SIZE]bool
}

func (board *Board) Cell(x, y int) Color {
	switch {
	case board.cells[BLACK][x+y*BOARD_SIZE]:
		return BLACK
	case board.cells[WHITE][x+y*BOARD_SIZE]:
		return WHITE
	default:
		return EMPTY
	}
}

func (board *Board) SetCell(x, y int, color Color) {
	board.cells[color][x+y*BOARD_SIZE] = true
	board.cells[color.Opponent()][x+y*BOARD_SIZE] = false
}

func (board *Board) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("----------------\n")
	defer buffer.WriteString("----------------\n")
	for y := 0; y < BOARD_SIZE; y++ {
		for x := 0; x < BOARD_SIZE; x++ {
			switch board.Cell(x, y) {
			case EMPTY:
				buffer.WriteString(" .")
			case BLACK:
				buffer.WriteString(" X")
			case WHITE:
				buffer.WriteString(" O")
			}
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

type Position struct {
	X, Y int
}

/**
 * A PositionIterator is a function iterating on x and y
 * returned parameters are :
 * x, y int
 * lineStart bool: whether we are starting a "line"
 * run bool: whether the iteration should go on or not
 */
type PositionIterator func() (int, int, bool, bool)

func iterStraight(positions *[][]Position, horizontal bool) {
	x, y := 0, 0
	rx, ry := &x, &y
	if !horizontal {
		rx, ry = ry, rx
	}
	for y = 0; y < BOARD_SIZE; y++ {
		line := make([]Position, BOARD_SIZE)
		for x = 0; x < BOARD_SIZE; x++ {
			line[x] = Position{*rx, *ry}
		}
		*positions = append(*positions, line)
	}
}

func iterDiag(positions *[][]Position, mirror bool) {
	x, y, n := 0, 0, 0
	rx, ry := &x, &y
	var line []Position
	if !mirror {
		rx, ry = ry, rx
	}
	for n < 2*BOARD_SIZE {
		diagonalLength := n % BOARD_SIZE
		startLine := x == diagonalLength
		if startLine {
			n++
			if line != nil {
				*positions = append(*positions, line)
			}
			line = make([]Position, 0)
			x = 0
		} else {
			x++
		}
		y = diagonalLength - x
		if n < BOARD_SIZE {
			line = append(line, Position{*rx, *ry})
		} else {
			line = append(line, Position{BOARD_SIZE - 1 - *rx, BOARD_SIZE - 1 - *ry})
		}
	}
	*positions = append(*positions, line)
}

func findPossibleMoves(moves *[][]Position, positions [][]Position, board *Board, color Color) {
	opponent := color.Opponent()
	var move []Position
	for _, line := range positions {
		prev := Position{0, 0}
		isCrossingOpponent := false
		for _, position := range line {
			x, y := position.X, position.Y
			switch {
			case board.Cell(x, y) == opponent && board.Cell(prev.X, prev.Y) != opponent:
				isCrossingOpponent = true
				move = make([]Position, 0)
				move = append(move, Position{prev.X, prev.Y})
			case isCrossingOpponent && board.Cell(x, y) != opponent:
				isCrossingOpponent = false
				if board.Cell(move[0].X, move[0].Y) == board.Cell(x, y) {
					move = nil
				} else {
					move = append(move, Position{x, y})
					*moves = append(*moves, move)
				}
			}
			if isCrossingOpponent {
				move = append(move, Position{x, y})
			}
			prev = Position{x, y}
		}
	}
}

func PossibleMoves(board *Board, color Color) [][]Position {
	moves := make([][]Position, 0)
	findPossibleMoves(&moves, boardRun, board, color)
	return moves
}

func TestBoardRun() [][]Position {
	return boardRun
}

func init() {
	boardRun = make([][]Position, 0)
	iterStraight(&boardRun, true)
	iterStraight(&boardRun, false)
	iterDiag(&boardRun, false)
	iterDiag(&boardRun, true)
}
