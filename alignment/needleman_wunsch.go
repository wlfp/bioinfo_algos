package alignment

import (
	"bufio"
	"fmt"
	"os"
)

type scoreMatrix struct {
	insertionAmount int
	deletionAmount  int
	matchAmount     int
	mismatchAmount  int
}

func calculateAlignment(u, v string, scoreMatrix scoreMatrix) {
	grid := computeAlignmentGrid(u, v, scoreMatrix)
	fmt.Printf("The maximal cost alignment between the two sequences is %d.\n", grid.GetElement(len(u), len(v)).costToReachSquare)
	fmt.Println("With alignment sequences:")
	grid.traceAlignmentSequences(u, v)
}

func AlignmentExample() {
	u := "CGTGAA"
	v := "GACTTAC"
	scoreMatrix := scoreMatrix{insertionAmount: -4, deletionAmount: -4, matchAmount: 5, mismatchAmount: -3}
	calculateAlignment(u, v, scoreMatrix)
}

func Alignment() {
	var u, v string
	var scoreMatrix scoreMatrix

	fmt.Printf("Please enter a filename with both sequences, one per line: ")
	var filename string
	fmt.Scanln(&filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	u = scanner.Text()
	scanner.Scan()
	v = scanner.Text()

	fmt.Println("Please enter, as integers separated by a space, the score function values for an insertion, deletion, match, and mismatch:")
	fmt.Scanf("%d %d %d %d\n", &scoreMatrix.insertionAmount, &scoreMatrix.deletionAmount, &scoreMatrix.matchAmount, &scoreMatrix.mismatchAmount)
	calculateAlignment(u, v, scoreMatrix)
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
	type alignmentState struct {
		refOffset    int
		queryOffset  int
		referenceStr []byte
		queryStr     []byte
	}

	var alignmentSequences [][2]string
	stack := []alignmentState{{refOffset: 1, queryOffset: 1}}

	// DFS-style walk.
	for len(stack) > 0 {
		state := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		row := grid.numRows - state.refOffset
		column := grid.numColumns - state.queryOffset

		if row == 0 && column == 0 {
			alignmentSequences = append(alignmentSequences, [2]string{
				reverseByteArrToString(append([]byte{}, state.referenceStr...)),
				reverseByteArrToString(append([]byte{}, state.queryStr...)),
			})
			continue
		}

		backpointers := grid.GetElement(row, column).backpointers
		for _, prevSquare := range backpointers {
			nextState := alignmentState{
				refOffset:    state.refOffset,
				queryOffset:  state.queryOffset,
				referenceStr: append([]byte{}, state.referenceStr...),
				queryStr:     append([]byte{}, state.queryStr...),
			}

			switch prevSquare {
			case North:
				nextState.referenceStr = append(nextState.referenceStr, u[len(u)-nextState.refOffset])
				nextState.refOffset++
				nextState.queryStr = append(nextState.queryStr, '-')
			case West:
				nextState.referenceStr = append(nextState.referenceStr, '-')
				nextState.queryStr = append(nextState.queryStr, v[len(v)-nextState.queryOffset])
				nextState.queryOffset++
			case NorthWest:
				nextState.referenceStr = append(nextState.referenceStr, u[len(u)-nextState.refOffset])
				nextState.refOffset++
				nextState.queryStr = append(nextState.queryStr, v[len(v)-nextState.queryOffset])
				nextState.queryOffset++
			}

			stack = append(stack, nextState)
		}
	}

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