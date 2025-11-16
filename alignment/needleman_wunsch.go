package main

import (
	"fmt"
	"strings"
)

type gridEntry struct {
	costToReachSqure int
	backpointers     []int // Use grid indices rather than real pointers.
}

type alignmentGrid [][]gridEntry

type scoreMatrix struct {
	indelAmount    int
	matchAmount    int
	mismatchAmount int
}

func computeAlignmentGrid(u, v string, scoreMatrix scoreMatrix) alignmentGrid {
	grid := make(alignmentGrid, len(u))
	for rowIndex := range grid {
		grid[rowIndex] = make([]gridEntry, len(v))
	}

	grid[0][0] = gridEntry{costToReachSqure: 0, backpointers: nil}
	for emptyFirstRowSquare := range grid[0][1:] {
		firstRowElementIndex := emptyFirstRowSquare + 1 // Already populated (0,0).
		grid[0][firstRowElementIndex] = gridEntry{
			costToReachSqure: scoreMatrix.indelAmount * firstRowElementIndex,
			backpointers:     []int{firstRowElementIndex - 1},
		}
	}
	for emptyFirstColumn := range grid[1:] {
		columnIndex := emptyFirstColumn + 1
		grid[columnIndex][0] = gridEntry{
			costToReachSqure: scoreMatrix.indelAmount * columnIndex,
			backpointers:     []int{columnIndex - 1},
		}
	}

	return grid
}

func main() {
	grid := computeAlignmentGrid("CGTGAA", "GACTTAC", scoreMatrix{indelAmount: -4, matchAmount: 5, mismatchAmount: -3})
	fmt.Println(grid)
}

func (grid alignmentGrid) String() string {
	builder := strings.Builder{}
	for rowIndex, row := range grid {
		if rowIndex != 0 {
			builder.WriteRune('\n')
		}
		builder.WriteString(fmt.Sprint(row))
	}
	return builder.String()
}
