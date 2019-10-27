package main

import (
	"fmt"
	"math"
	"strings"
)

// A Player is either Black or White
type Player int

// The Set of possible players.
const (
	None = Player(iota)
	Black
	White
)

var playerNames = []string{"None", "Black", "White"}

func (p Player) String() string {
	return playerNames[p]
}

func (p Player) opponent() Player {
	switch p {
	case Black:
		return White
	case White:
		return Black
	default:
		return None
	}
}

func (p Player) rune() rune {
	switch p {
	case Black:
		return 'X'
	case White:
		return '0'
	default:
		return '.'
	}
}

// A Position is a place on the Othello Board
type Position struct {
	x int
	y int
}

func (pos Position) String() string {
	return fmt.Sprintf("(%d, %d)", pos.x, pos.y)
}

func (pos Position) neighbor(dx int, dy int) Position {
	return Position{pos.x + dx, pos.y + dy}
}

func abs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (pos Position) distance(other Position) int {
	return max(abs(pos.x-other.x), abs(pos.y-other.y))
}

func (pos Position) isValid() bool {
	return pos.x >= 0 && pos.x < 8 && pos.y >= 0 && pos.y < 8
}

// The Board for the Othello game.
type Board [8][8]Player

// Init the Othello Board for the start of a game.
func (board *Board) Init() {
	board[3][3] = White
	board[4][4] = White
	board[3][4] = Black
	board[4][3] = Black
}

func (board *Board) playerAt(pos Position) Player {
	return board[pos.y][pos.x]
}

func (board *Board) playerScore(player Player) (score int) {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if board[y][x] == player {
				score++
			}
		}
	}
	return
}

func (board *Board) findBridgeCandidate(pos Position, player Player, dx int, dy int) []Position {
	if !pos.isValid() || board.playerAt(pos) != None {
		return nil
	}
	currentPos := pos.neighbor(dx, dy)

	for currentPos.isValid() && board.playerAt(currentPos) == player.opponent() {
		currentPos = currentPos.neighbor(dx, dy)
	}

	bridgeLen := currentPos.distance(pos)

	if currentPos.isValid() && board.playerAt(currentPos) == player && bridgeLen > 1 {
		bridge := make([]Position, bridgeLen)
		for i := 0; i < bridgeLen; i++ {
			bridge[i] = pos
			pos = pos.neighbor(dx, dy)
		}
		return bridge
	}
	return nil
}

type direction struct {
	dx int
	dy int
}

var playDirections = []direction{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

// PlayAt plays a piecs
func (board *Board) PlayAt(pos Position, player Player) bool {
	validPlay := false

	for _, dir := range playDirections {
		bridge := board.findBridgeCandidate(pos, player, dir.dx, dir.dy)
		if bridge != nil {
			validPlay = true
			for _, position := range bridge {
				board[position.y][position.x] = player
			}
		}
	}
	return validPlay
}

func (board *Board) String() string {
	var str strings.Builder
	for _, line := range board {
		for _, color := range line {
			str.WriteRune(color.rune())
		}
		str.WriteString("\n")
	}
	return str.String()
}

func negamaxAB(board *Board, depth int, alpha int, beta int, player Player) int {
	if depth == 0 {
		return board.playerScore(player)
	}
	terminalNode := true
	score := math.MinInt32
forEachNodes:
	for y, line := range board {
		for x := range line {
			boardCopy := *board
			child := &boardCopy
			if child.PlayAt(Position{x, y}, player) {
				terminalNode = false
				score = max(score, -negamaxAB(child, depth-1, -beta, -alpha, player.opponent()))
				alpha = max(alpha, score)
				if alpha >= beta {
					break forEachNodes
				}
			}
		}
	}
	if terminalNode {
		score = board.playerScore(player)
	}
	return score
}

// Negamax evaluates the situation of a Player on a given Board.
func Negamax(board *Board, depth int, player Player) int {
	return negamaxAB(board, depth, math.MinInt32, math.MaxInt32, player)
}

func main() {
	board := new(Board)
	board.Init()
	board.PlayAt(Position{2, 3}, Black)
	fmt.Println(board)
	fmt.Printf("black negamax : %d\n", Negamax(board, 4, Black))
	board.PlayAt(Position{2, 2}, White)
	fmt.Println(board)
	fmt.Printf("black negamax : %d\n", Negamax(board, 4, Black))
}
