package astar

import (
	"puzzle3x3/game"
	"sort"
)

type Node struct {
	Move      game.ACTION
	Board     game.Board
	BasePoint game.Point
	G         int
	H         int
	F         int
	Parent    *Node
}

func createNode(board game.Board, point game.Point, G, H int, move game.ACTION, parent *Node) *Node {
	return &Node{
		Board:     board,
		BasePoint: point,
		G:         G,
		H:         H,
		F:         G + H,
		Parent:    parent,
		Move:      move,
	}
}

func (n Node) GetValidNeighborPoints() map[game.ACTION]game.Point {
	possibleMoves := []game.ACTION{game.MOVE_DOWN, game.MOVE_UP, game.MOVE_RIGHT, game.MOVE_LEFT}

	result := map[game.ACTION]game.Point{}
	for _, move := range possibleMoves {
		newPoint := n.BasePoint.Move(move)
		if n.Board.ValidatePoint(newPoint) {
			result[move] = newPoint
		}
	}

	return result
}

func reconstructMoves(goalNode *Node) []game.ACTION {
	moves := []game.ACTION{}

	if goalNode != nil && goalNode.Move != game.MOVE_NONE {
		moves = append([]game.ACTION{goalNode.Move}, moves...)
		nextMoves := reconstructMoves(goalNode.Parent)
		moves = append(nextMoves, moves...)
	}

	return moves
}

func FindSteps(board game.Board, basePoint game.Point, mapResultTarget map[int]game.Point) []game.ACTION {
	startNode := createNode(board, basePoint, 0, board.CalcaulateDistanceFromTarget(mapResultTarget), game.MOVE_NONE, nil)
	openList := []*Node{startNode}
	openMap := map[string]*Node{board.ToString(): startNode}
	closedMap := map[string]*Node{}

	var currentNode *Node
	for len(openList) > 0 {
		currentNode, openList = openList[0], openList[1:]
		currentBoard := currentNode.Board

		if currentBoard.CalcaulateDistanceFromTarget(mapResultTarget) == 0 {
			return reconstructMoves(currentNode)
		}

		closedMap[currentBoard.ToString()] = currentNode
		for move, neighborPoint := range currentNode.GetValidNeighborPoints() {
			neighborBoard := copyBoard(currentBoard)
			neighborBoard.SwapPoint(currentNode.BasePoint, neighborPoint)
			neighborBoardString := neighborBoard.ToString()
			if _, ok := closedMap[neighborBoardString]; ok {
				continue
			}

			tentativeG := currentNode.G + 1

			if _, ok := openMap[neighborBoardString]; !ok {
				neighbor := createNode(neighborBoard, neighborPoint, tentativeG, neighborBoard.CalcaulateDistanceFromTarget(mapResultTarget), move, currentNode)
				openList = appendData(openList, neighbor)
				openMap[neighborBoardString] = neighbor
			} else if tentativeG < openMap[neighborBoardString].G {
				neighbor := openMap[neighborBoardString]
				neighbor.Move = move
				neighbor.G = tentativeG
				neighbor.F = tentativeG + neighbor.H
				neighbor.Parent = currentNode
			}
		}
	}

	return []game.ACTION{}
}

func appendData(data []*Node, item *Node) []*Node {
	i := sort.Search(len(data), func(i int) bool { return data[i].F >= item.F })
	data = append(data, &Node{})
	copy(data[i+1:], data[i:])
	data[i] = item
	return data
}

func copyBoard(board game.Board) game.Board {
	tmpBoard := make(game.Board, len(board))
	for i := range board {
		tmpBoard[i] = make([]int, len(board[i]))
		copy(tmpBoard[i], board[i])
	}

	return tmpBoard
}
