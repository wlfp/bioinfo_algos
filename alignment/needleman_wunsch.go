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
			backpointers:      []moveType{North}})
	}

	for colNum := range grid.numColumns {
		if colNum == 0 {
			continue
		}
		grid.SetElement(0, colNum, gridEntry{
			costToReachSquare: scoreMatrix.deletionAmount * colNum,
			backpointers:      []moveType{West}})
	}
}

/*
In this implementation, u is considered the reference against a query v.
*/
func computeAlignmentGrid(u, v string, scoreMatrix scoreMatrix) *alignmentGrid {
	grid := NewGrid(len(u)+1, len(v)+1)
	initialiseGrid(grid, scoreMatrix)

	for rowNumber := 1; rowNumber < grid.numRows; rowNumber++ {
		for columnNumber := 1; columnNumber < grid.numColumns; columnNumber++ {
			isAMatch := u[rowNumber-1] == v[columnNumber-1]
			minimumCostAndPath := grid.findOptimalMove(rowNumber, columnNumber, isAMatch, scoreMatrix)
			grid.SetElement(rowNumber, columnNumber, minimumCostAndPath)
		}
	}

	return grid
}

func (grid *alignmentGrid) findOptimalMove(rowNumber, columnNumber int, isAMatch bool, scoreMatrix scoreMatrix) gridEntry {
	type previousSquare struct {
		costToReachSquare int
		index             moveType
	}
	moves := [3]previousSquare{}

	west := grid.GetElement(rowNumber, columnNumber-1)
	north := grid.GetElement(rowNumber-1, columnNumber)
	moves[0] = previousSquare{
		costToReachSquare: west.costToReachSquare + scoreMatrix.insertionAmount,
		index:             West,
	}
	moves[1] = previousSquare{
		costToReachSquare: north.costToReachSquare + scoreMatrix.deletionAmount,
		index:             North,
	}

	northWest := grid.GetElement(rowNumber-1, columnNumber-1)
	northWestMoveCost := scoreMatrix.mismatchAmount
	if isAMatch {
		northWestMoveCost = scoreMatrix.matchAmount
	}
	moves[2] = previousSquare{
		costToReachSquare: northWest.costToReachSquare + northWestMoveCost,
		index:             NorthWest,
	}

	bestCost := max(moves[0].costToReachSquare, max(moves[1].costToReachSquare, moves[2].costToReachSquare))
	var backpointers []moveType
	for _, move := range moves {
		if move.costToReachSquare == bestCost {
			backpointers = append(backpointers, move.index)
		}
	}

	return gridEntry{costToReachSquare: bestCost, backpointers: backpointers}
}

func (grid *alignmentGrid) traceAlignmentSequences(u, v string) {
	alignmentSequences := make([][2]string, 1)
	var refOffset, queryOffset int = 1, 1
	backpointers := grid.GetElement(grid.numRows-refOffset, grid.numColumns-queryOffset).backpointers
	var referenceString, queryString []byte
	for backpointers != nil {
		// TODO: Split logic here for the rest of the strings.
		prevSquare := backpointers[0]
		switch prevSquare {
		case North:
			referenceString = append(referenceString, u[len(u)-refOffset])
			refOffset++
			queryString = append(queryString, '-')
		case West:
			referenceString = append(referenceString, '-')
			queryString = append(queryString, v[len(v)-queryOffset])
			queryOffset++
		case NorthWest:
			referenceString = append(referenceString, u[len(u)-refOffset])
			refOffset++
			queryString = append(queryString, v[len(v)-queryOffset])
			queryOffset++
		}

		backpointers = grid.GetElement(grid.numRows-refOffset, grid.numColumns-queryOffset).backpointers
	}

	alignmentSequences[0] = [2]string{reverseByteArrToString(referenceString), reverseByteArrToString(queryString)}

	for seqNum, sequence := range alignmentSequences {
		if seqNum != 0 {
			fmt.Println()
		}
		fmt.Printf("%s\n%s\n", sequence[0], sequence[1])
	}
}

func reverseByteArrToString(buf []byte) string {
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf)
}

func AlignmentExample() {
	u := "CGTGAA"
	v := "GACTTAC"
	grid := computeAlignmentGrid(u, v, scoreMatrix{insertionAmount: -4, deletionAmount: -4, matchAmount: 5, mismatchAmount: -3})
	fmt.Printf("The maximal cost alignment between the two sequences is %d.\n", grid.GetElement(len(u), len(v)).costToReachSquare)
	fmt.Println("With alignment sequences:")
	grid.traceAlignmentSequences(u, v)
}
