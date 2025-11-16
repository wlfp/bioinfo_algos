package alignment

import (
	"fmt"
	"strings"
)

type gridEntry struct {
	costToReachSquare int
	backpointers      [][2]int // Use grid indices rather than real pointers.
}

type alignmentGrid struct {
	entries    []gridEntry
	numRows    int
	numColumns int
}

func NewGrid(numRows, numColumns int) *alignmentGrid {
	grid := alignmentGrid{numRows: numRows, numColumns: numColumns}
	grid.entries = make([]gridEntry, numRows*numColumns)
	return &grid
}

func (grid *alignmentGrid) elementIndex(row, column int) int {
	return column + grid.numColumns*row
}

func (grid *alignmentGrid) SetElement(row, column int, value gridEntry) {
	grid.entries[grid.elementIndex(row, column)] = value
}

func (grid *alignmentGrid) GetElement(row, column int) gridEntry {
	return grid.entries[grid.elementIndex(row, column)]
}

func (grid alignmentGrid) String() string {
	builder := strings.Builder{}
	for entryIndex, entry := range grid.entries {
		if entryIndex%grid.numColumns == 0 {
			builder.WriteRune('\n')
		}
		builder.WriteString(fmt.Sprint(entry))
	}
	return builder.String()
}
