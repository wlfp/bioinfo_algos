package alignment

import (
	"fmt"
)

type scoreMatrix struct {
	insertionAmount int
	deletionAmount  int
	matchAmount     int
	mismatchAmount  int
}

func initialiseGrid(grid *alignmentGrid, scoreMatrix scoreMatrix) {
	grid.SetElement(0, 0, gridEntry{costToReachSquare: 0, backpointers: nil})

	for rowNum := range grid.numRows {
		if rowNum == 0 {
			continue // Already set first.
		}
		grid.SetElement(rowNum, 0, gridEntry{
			costToReachSquare: scoreMatrix.insertionAmount * rowNum,
			backpointers:      [][2]int{{rowNum - 1, 0}}})
	}

	for colNum := range grid.numColumns {
		if colNum == 0 {
			continue
		}
		grid.SetElement(0, colNum, gridEntry{
			costToReachSquare: scoreMatrix.deletionAmount * colNum,
			backpointers:      [][2]int{{0, colNum - 1}}})
	}
}

/*
In this implementation, u is considered the reference against a query v.
*/
func computeAlignmentGrid(u, v string, scoreMatrix scoreMatrix) *alignmentGrid {
	grid := NewGrid(len(u), len(v))
	initialiseGrid(grid, scoreMatrix)

	for rowNumber := 1; rowNumber < grid.numRows; rowNumber++ {
		for columnNumber := 1; columnNumber < grid.numColumns; columnNumber++ {
			isAMatch := u[rowNumber] == v[columnNumber]
			minimumCostAndPath := grid.findOptimalMove(rowNumber, columnNumber, isAMatch, scoreMatrix)
			grid.SetElement(rowNumber, columnNumber, minimumCostAndPath)
		}
	}

	return grid
}

func (grid *alignmentGrid) findOptimalMove(rowNumber, columnNumber int, isAMatch bool, scoreMatrix scoreMatrix) gridEntry {
	type previousSquare struct {
		costToReachSquare int
		index             [2]int
	}
	moves := [3]previousSquare{}

	west := grid.GetElement(rowNumber, columnNumber-1)
	north := grid.GetElement(rowNumber-1, columnNumber)
	moves[0] = previousSquare{
		costToReachSquare: west.costToReachSquare + scoreMatrix.insertionAmount,
		index:             [2]int{rowNumber, columnNumber - 1},
	}
	moves[1] = previousSquare{
		costToReachSquare: north.costToReachSquare + scoreMatrix.deletionAmount,
		index:             [2]int{rowNumber - 1, columnNumber},
	}

	northWest := grid.GetElement(rowNumber-1, columnNumber-1)
	northWestMoveCost := scoreMatrix.mismatchAmount
	if isAMatch {
		northWestMoveCost = scoreMatrix.matchAmount
	}
	moves[2] = previousSquare{
		costToReachSquare: northWest.costToReachSquare + northWestMoveCost,
		index:             [2]int{rowNumber - 1, columnNumber - 1},
	}

	bestCost := max(moves[0].costToReachSquare, max(moves[1].costToReachSquare, moves[2].costToReachSquare))
	var backpointers [][2]int
	for _, move := range moves {
		if move.costToReachSquare == bestCost {
			backpointers = append(backpointers, move.index)
		}
	}

	return gridEntry{costToReachSquare: bestCost, backpointers: backpointers}
}

func AlignmentExample() {
	grid := computeAlignmentGrid("CGTGAA", "GACTTAC", scoreMatrix{insertionAmount: -4, deletionAmount: -4, matchAmount: 5, mismatchAmount: -3})
	fmt.Println(grid)
}
