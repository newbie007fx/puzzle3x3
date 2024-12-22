package game

import "fmt"

type Board [][]int

func (b Board) ToString() string {
	return fmt.Sprintf("%+v", b)
}

func (b Board) SwapPoint(basePoint Point, targetPoint Point) bool {
	if b.ValidatePoint(targetPoint) {
		b[basePoint.PointX][basePoint.PointY] = b[targetPoint.PointX][targetPoint.PointY]
		b[targetPoint.PointX][targetPoint.PointY] = 0
		return true
	}

	return false
}

func (b Board) ValidatePoint(pos Point) bool {
	return pos.PointX >= 0 && pos.PointX < len(b) && pos.PointY >= 0 && pos.PointY < len(b)
}

func (b Board) CalcaulateDistanceFromTarget(mapResultTarget map[int]Point) int {
	totalDistance := 0

	number := 1
	for row := range b {
		for col := range b[row] {
			if b[row][col] != number && b[row][col] != 0 {
				point := mapResultTarget[b[row][col]]
				totalDistance += b.calculateDistance(row, col, point)
			}
			number++
		}
	}

	return totalDistance
}

func (b Board) calculateDistance(x, y int, point Point) int {
	xDistance := x - point.PointX
	if xDistance < 0 {
		xDistance *= -1
	}

	yDistance := y - point.PointY
	if yDistance < 0 {
		yDistance *= -1
	}

	return xDistance + yDistance
}

func (b Board) Print() {
	for row := range b {
		fmt.Println("-------")
		fmt.Print("|")
		for col := range b[row] {
			if b[row][col] == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print(b[row][col])
			}
			fmt.Print("|")
		}
		fmt.Println("")
	}
	fmt.Println("-------")
}
