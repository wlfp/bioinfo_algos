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
	grid.SetElement(0, 0, gridEntry{costToReachSqure: 0, backpointers: nil})

	for rowNum := range grid.numRows {
		if rowNum == 0 {
			continue // Already set first.
		}
		grid.SetElement(rowNum, 0, gridEntry{
			costToReachSqure: scoreMatrix.insertionAmount * rowNum,
			backpointers:     [][2]int{{rowNum - 1, 0}}})
	}

	for colNum := range grid.numColumns {
		if colNum == 0 {
			continue
		}
		grid.SetElement(0, colNum, gridEntry{
			costToReachSqure: scoreMatrix.deletionAmount * colNum,
			backpointers:     [][2]int{{0, colNum - 1}}})
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
	west := grid.GetElement(rowNumber, columnNumber-1)
	north := grid.GetElement(rowNumber-1, columnNumber)
	cost := max(west.costToReachSqure+scoreMatrix.insertionAmount, north.costToReachSqure+scoreMatrix.deletionAmount)

	northWest := grid.GetElement(rowNumber-1, columnNumber-1)
	northWestMoveCost := scoreMatrix.mismatchAmount
	if isAMatch {
		northWestMoveCost = scoreMatrix.matchAmount
	}
	cost = max(cost, northWest.costToReachSqure+northWestMoveCost)
	// TODO: Backpointers building along the way.
	return gridEntry{costToReachSqure: cost, backpointers: nil}
}

func AlignmentExample() {
	grid := computeAlignmentGrid("CGTGAA", "GACTTAC", scoreMatrix{insertionAmount: -4, deletionAmount: -4, matchAmount: 5, mismatchAmount: -3})
	fmt.Println(grid)
}
