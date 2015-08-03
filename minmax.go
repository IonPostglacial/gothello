package othello

const MAX_DEPTH = 1

func evaluateMove(move []Position, board Board, player int, currentPlayer int, depth int) int {
  modifiedBoard := &board
  for _, position := range move {
    modifiedBoard.SetCell(position.X, position.Y, currentPlayer)
  }
  if depth == 0 {
    score := len(move)
    if player != currentPlayer {
      return score
    } else {
      return -score
    }
  } else {
    return minMax(modifiedBoard, player, OpponentColor(currentPlayer), depth - 1)
  }
}

func minMax(board *Board, player int, currentPlayer int, depth int) int {
  isMinNode := player != currentPlayer
  possibleMoves := PossibleMoves(board, currentPlayer)
  globalScore := 0

  for _, move := range possibleMoves {
    score := evaluateMove(move, *board, player, currentPlayer, depth)
    if score < globalScore != isMinNode {
      globalScore = score
    }
  }
  return globalScore
}

func EvaluateMove(move []Position, board *Board, player int, currentPlayer int) int {
  return evaluateMove(move, *board, player, currentPlayer, MAX_DEPTH)
}
