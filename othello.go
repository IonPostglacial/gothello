package othello

import "bytes"

const BOARD_SIZE = 8

type Color int8

const (
	EMPTY Color = iota
	BLACK
	WHITE
)

func (color Color) Opponent() Color {
	switch color {
	case BLACK:
		return WHITE
	case WHITE:
		return BLACK
	default:
		return EMPTY
	}
}

func (color Color) String() string {
	switch color {
	case BLACK:
		return "BLACK"
	case WHITE:
		return "WHITE"
	default:
		return "EMPTY"
	}
}

type Board struct {
	cells [BOARD_SIZE * BOARD_SIZE]Color
}

func (board *Board) Cell(x, y int) Color {
	return board.cells[x+y*BOARD_SIZE]
}

func (board *Board) SetCell(x, y int, color Color) {
	board.cells[x+y*BOARD_SIZE] = color
}

func (board *Board) Count(color Color) int {
	count := 0
	for y := 0; y < BOARD_SIZE; y++ {
		for x := 0; x < BOARD_SIZE; x++ {
			if board.Cell(x, y) == color {
				count++
			}
		}
	}
	return count
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
	H, V int
}

type movePart struct {
	move   Position
	tokens []Position
}

func sendPossibleMoves(board *Board, player Color, move []Position, h, v int, output chan<- movePart) []Position {
	switch {
	case board.Cell(h, v) == player:
		move = nil
	case board.Cell(h, v) == player.Opponent():
		move = append(move, Position{h, v})
	case board.Cell(h, v) == EMPTY && len(move) > 0:
		move = append(move, Position{h, v})
		output <- movePart{
			move:   Position{h, v},
			tokens: move,
		}
		move = nil
	default:
		move = nil
	}
	return move
}

func PossibleMoves(board *Board, color Color) map[Position][]Position {
	moves := make(map[Position][]Position)
	for v := 0; v < BOARD_SIZE; v++ {
		for h := 0; h < BOARD_SIZE; h++ {
			if board.Cell(h, v) == color {
				movesChan := make(chan movePart)
				go func() {
					var move []Position
					for rh := h + 1; rh < BOARD_SIZE; rh++ { // EAST
						if move = sendPossibleMoves(board, color, move, rh, v, movesChan); move == nil {
							break
						}
					}
					move = nil
					for rv := v + 1; rv < BOARD_SIZE; rv++ { // SOUTH
						if move = sendPossibleMoves(board, color, move, h, rv, movesChan); move == nil {
							break
						}
					}
					move = nil
					for rh := h - 1; rh >= 0; rh-- { // WEST
						if move = sendPossibleMoves(board, color, move, rh, v, movesChan); move == nil {
							break
						}
					}
					move = nil
					for rv := v - 1; rv >= 0; rv-- { // NORTH
						if move = sendPossibleMoves(board, color, move, h, rv, movesChan); move == nil {
							break
						}
					}
					move = nil
					for rv, rh := v+1, h+1; rv < BOARD_SIZE && rh < BOARD_SIZE; rv, rh = rv+1, rh+1 { // SOUTH EAST
						if move = sendPossibleMoves(board, color, move, rh, rv, movesChan); move == nil {
							break
						}
					}
					move = nil
					for rv, rh := v+1, h-1; rv < BOARD_SIZE && rh >= 0; rv, rh = rv+1, rh-1 { // SOUTH WEST
						if move = sendPossibleMoves(board, color, move, rh, rv, movesChan); move == nil {
							break
						}
					}
					move = nil
					for rv, rh := v-1, h+1; rv >= 0 && rh < BOARD_SIZE; rv, rh = rv-1, rh+1 { // NORTH EAST
						if move = sendPossibleMoves(board, color, move, rh, rv, movesChan); move == nil {
							break
						}
					}
					move = nil
					for rv, rh := v-1, h-1; rv >= 0 && rh >= 0; rv, rh = rv-1, rh-1 { // NORTH WEST
						if move = sendPossibleMoves(board, color, move, rh, rv, movesChan); move == nil {
							break
						}
					}
					close(movesChan)
				}()
				for movePart := range movesChan {
					move := movePart.move
					for _, pos := range movePart.tokens {
						moves[move] = append(moves[move], pos)
					}
				}
			}
		}
	}
	return moves
}
