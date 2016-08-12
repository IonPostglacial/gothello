package othello

const MAX_DEPTH = 5

func evaluateMove(move []Position, board Board, player Color, currentPlayer Color, depth int) int {
	modifiedBoard := &board
	for _, position := range move {
		modifiedBoard.SetCell(position.H, position.V, currentPlayer)
	}
	if depth == 0 {
		score := len(move)
		if player != currentPlayer {
			return score
		} else {
			return -score
		}
	} else {
		return minMax(modifiedBoard, player, currentPlayer.Opponent(), depth-1)
	}
}

func minMax(board *Board, player Color, currentPlayer Color, depth int) int {
	possibleMoves := PossibleMoves(board, currentPlayer)
	globalScore := 0

	for _, move := range possibleMoves {
		score := evaluateMove(move, *board, player, currentPlayer, depth)
		if score > globalScore {
			globalScore = score
		}
	}
	return globalScore
}

func EvaluateMove(move []Position, board *Board, player Color, currentPlayer Color) int {
	return evaluateMove(move, *board, player, currentPlayer, MAX_DEPTH)
}
