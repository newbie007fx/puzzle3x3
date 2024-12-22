package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"puzzle3x3/astar"
	"puzzle3x3/game"
	"time"
)

type BoardGame struct {
	board           game.Board
	point           game.Point
	mapResultTarget map[int]game.Point
}

func createBoardGame(length int) *BoardGame {
	boardGame := BoardGame{
		board:           make([][]int, length),
		mapResultTarget: map[int]game.Point{},
	}

	dataSlice := generateSliceData(length * length)

	index := 0
	for row := range boardGame.board {
		boardGame.board[row] = make([]int, length)
		for col := range length {
			boardGame.board[row][col] = dataSlice[index]
			point := game.Point{
				PointX: row,
				PointY: col,
			}
			if dataSlice[index] == 0 {
				boardGame.point = point
			}
			if index < (length*length)-1 {
				boardGame.mapResultTarget[index+1] = point
			}
			index++
		}
	}

	return &boardGame
}

func generateSliceData(length int) []int {
	dataSlice := rand.Perm(length)

	var countInversion, lastI, lastJ int
	for i := range len(dataSlice) - 1 {
		for j := i + 1; j < len(dataSlice); j++ {
			if dataSlice[i] > 0 && dataSlice[j] > 0 && dataSlice[i] > dataSlice[j] {
				countInversion++
				lastI = i
				lastJ = j
			}
		}
	}

	if countInversion%2 == 1 {
		dataSlice[lastI], dataSlice[lastJ] = dataSlice[lastJ], dataSlice[lastI]
	}

	return dataSlice
}

func (bg *BoardGame) moveBasePoint(point game.Point, action game.ACTION) {
	targetPoint := point.Move(action)

	if bg.board.SwapPoint(point, targetPoint) {
		bg.point = targetPoint
	}
}

func (bg *BoardGame) runGame(action game.ACTION) bool {
	bg.moveBasePoint(bg.point, action)
	bg.board.Print()
	countResult := bg.board.CalcaulateDistanceFromTarget(bg.mapResultTarget)
	if countResult == 0 {
		fmt.Println("Congratulations, you have completed the game")
		return true
	}

	return false
}

func (bg *BoardGame) runAutoSolve() {
	fmt.Println("running auto solve")
	steps := astar.FindSteps(bg.board, bg.point, bg.mapResultTarget)
	if len(steps) > 0 {
		fmt.Printf("found %d steps, starting steps \n", len(steps))
		for _, step := range steps {
			fmt.Println("move = ", step.GetActionString())
			bg.runGame(step)
			time.Sleep(1 * time.Second)
		}
	} else {
		fmt.Println("No path found!")
	}
}

func main() {
	defaultSize := 3
	boardGame := createBoardGame(defaultSize)
	boardGame.board.Print()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("insert action: ")
		scanner.Scan()
		actionName := scanner.Text()

		if actionName == "auto solve" {
			boardGame.runAutoSolve()
			break
		}

		action, err := game.LoadActionFromString(actionName)
		if err != nil {
			fmt.Println("invalid action, try again")
			continue
		}

		if finish := boardGame.runGame(action); finish {
			break
		}
	}
}
