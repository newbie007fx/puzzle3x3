package astar

import (
	"context"
	"puzzle3x3/game"
	"sort"
	"sync"
	"time"
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

type AstarGameSolver struct {
	MapResultTarget map[int]game.Point
	OpenList        []*Node
	OpenMap         map[string]*Node
	ClosedMap       map[string]*Node
	FinalNode       *Node
	sync.RWMutex
}

func CreateGameSolver(mapResultTarget map[int]game.Point) *AstarGameSolver {
	return &AstarGameSolver{
		MapResultTarget: mapResultTarget,
		OpenList:        []*Node{},
		OpenMap:         map[string]*Node{},
		ClosedMap:       map[string]*Node{},
	}
}

func (a *AstarGameSolver) popNode() *Node {
	a.RWMutex.Lock()
	defer a.RWMutex.Unlock()

	if len(a.OpenList) == 0 {
		return nil
	}

	var currentNode *Node
	currentNode, a.OpenList = a.OpenList[0], a.OpenList[1:]
	delete(a.OpenMap, currentNode.Board.ToString())
	a.ClosedMap[currentNode.Board.ToString()] = currentNode

	return currentNode
}

func (a *AstarGameSolver) isNodeKeyClosed(key string) bool {
	a.RWMutex.RLock()
	defer a.RWMutex.RUnlock()

	_, ok := a.ClosedMap[key]

	return ok
}

func (a *AstarGameSolver) getOpenNodeByKey(key string) *Node {
	a.RWMutex.RLock()
	defer a.RWMutex.RUnlock()

	if node, ok := a.ClosedMap[key]; ok {
		return node
	}

	return nil
}

func (a *AstarGameSolver) setOpenNode(key string, node *Node) {
	a.RWMutex.Lock()
	defer a.RWMutex.Unlock()
	a.OpenList = a.appendData(a.OpenList, node)
	a.OpenMap[key] = node
}

func (a *AstarGameSolver) reconstructMoves(goalNode *Node) []game.ACTION {
	moves := []game.ACTION{}

	if goalNode != nil && goalNode.Move != game.MOVE_NONE {
		moves = append([]game.ACTION{goalNode.Move}, moves...)
		nextMoves := a.reconstructMoves(goalNode.Parent)
		moves = append(nextMoves, moves...)
	}

	return moves
}

func (a *AstarGameSolver) FindSteps(board game.Board, basePoint game.Point) []game.ACTION {
	startNode := createNode(board, basePoint, 0, board.CalcaulateDistanceFromTarget(a.MapResultTarget), game.MOVE_NONE, nil)
	a.OpenList = append(a.OpenList, startNode)
	a.OpenMap[board.ToString()] = startNode

	complete := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	for range 2 {
		go a.workerRun(ctx, complete)
	}

	<-complete
	cancel()

	return a.reconstructMoves(a.FinalNode)
}

func (a *AstarGameSolver) workerRun(ctx context.Context, completed chan int) {
	for {
		select {
		case <-ctx.Done():

			return
		default:
			if a.process() {
				completed <- 0
				return
			}
		}
	}
}

func (a *AstarGameSolver) process() bool {
	currentNode := a.popNode()
	if currentNode == nil {
		time.Sleep(time.Second)
		return false
	}

	currentBoard := currentNode.Board
	if currentBoard.CalcaulateDistanceFromTarget(a.MapResultTarget) == 0 {
		a.FinalNode = currentNode
		return true
	}

	for move, neighborPoint := range currentNode.GetValidNeighborPoints() {
		neighborBoard := a.copyBoard(currentBoard)
		neighborBoard.SwapPoint(currentNode.BasePoint, neighborPoint)
		neighborBoardString := neighborBoard.ToString()
		if a.isNodeKeyClosed(neighborBoardString) {
			continue
		}

		tentativeG := currentNode.G + 1
		neighbor := a.getOpenNodeByKey(neighborBoardString)
		if neighbor == nil {
			neighbor = createNode(neighborBoard, neighborPoint, tentativeG, neighborBoard.CalcaulateDistanceFromTarget(a.MapResultTarget), move, currentNode)
			a.setOpenNode(neighborBoardString, neighbor)
		} else if tentativeG < neighbor.G {
			neighbor.Move = move
			neighbor.G = tentativeG
			neighbor.F = tentativeG + neighbor.H
			neighbor.Parent = currentNode
		}
	}

	return false
}

func (a *AstarGameSolver) appendData(data []*Node, item *Node) []*Node {
	i := sort.Search(len(data), func(i int) bool { return data[i].F >= item.F })
	data = append(data, &Node{})
	copy(data[i+1:], data[i:])
	data[i] = item
	return data
}

func (a *AstarGameSolver) copyBoard(board game.Board) game.Board {
	tmpBoard := make(game.Board, len(board))
	for i := range board {
		tmpBoard[i] = make([]int, len(board[i]))
		copy(tmpBoard[i], board[i])
	}

	return tmpBoard
}
