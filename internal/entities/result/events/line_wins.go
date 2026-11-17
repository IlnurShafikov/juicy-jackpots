package events

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/IlnurShafikov/slot-engine/pkg/calculators"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
)

// LineWinsEvent represents a win event containing all winning lines from a single spin.
// It is serialised and sent to the client as part of the spin result.
type LineWinsEvent struct {
	// Event is the type identifier for this event (always WinsType).
	Event Type `json:"type"`
	// Positions is the deduplicated, sorted list of all winning symbol positions.
	Positions model.Positions `json:"positions"`
	// Lines contains the individual winning line details.
	Lines []Line `json:"lines,omitempty"`
}

// NewLineWinsEvent constructs a LineWinsEvent from the calculated line wins.
// It deduplicates and sorts winning positions, then builds the Line slice
// with position indices referencing the Positions array.
//
// Parameters:
//   - winLines: the line win results produced by the calculator
//
// Returns:
//   - pointer to the constructed LineWinsEvent
//   - error if any winning position cannot be resolved (ErrPositionNotFound)
func NewLineWinsEvent(winLines calculators.LineWins) (*LineWinsEvent, error) {
	pos, indicesMap := collectWinsPositions(winLines)

	lines, err := buildLines(winLines, indicesMap)
	if err != nil {
		return nil, fmt.Errorf("build win lines: %w", err)
	}

	return &LineWinsEvent{
		Event:     WinsType,
		Positions: pos,
		Lines:     lines,
	}, nil
}

// EventType returns the event type identifier, satisfying the Event interface.
func (c *LineWinsEvent) EventType() Type {
	return c.Event
}

// Line represents a single winning payline within a LineWinsEvent.
type Line struct {
	// PosIndices contains the indices into the parent LineWinsEvent.Positions slice
	// that identify which symbol positions contributed to this win.
	PosIndices []int `json:"posIndices"`
	// Line is the payline number.
	Line int `json:"line"`
	// Symbol is the symbol identifier that formed the win.
	Symbol int `json:"symbol"`
	// Payout is the win amount for this line.
	Payout uint64 `json:"payout"`
}

// collectWinsPositions extracts all unique winning positions from line wins.
// It creates a sorted list of positions and a map for quick index lookups.
//
// Parameters:
//   - lineWins: calculated line wins containing positions
//
// Returns:
//   - sorted slice of unique positions
//   - map from position to its index in the positions slice
func collectWinsPositions(
	lineWins calculators.LineWins,
) (positions []model.Position, indicesMap map[model.Position]int) {
	posMap := make(map[model.Position]struct{})

	for _, win := range lineWins {
		for _, pos := range win.Positions {
			posMap[pos] = struct{}{}
		}
	}

	positions = make([]model.Position, 0, len(posMap))
	for pos := range posMap {
		positions = append(positions, pos)
	}

	slices.SortFunc(positions, comparePos)

	indicesMap = make(map[model.Position]int, len(positions))
	for i, pos := range positions {
		indicesMap[pos] = i
	}

	return positions, indicesMap
}

// buildLines constructs Line objects from line wins and position indices.
// Each line contains references to positions via indices for efficient serialization.
//
// Parameters:
//   - lineWins: calculated line wins
//   - indexMap: map from position to its index in the positions array
//
// Returns:
//   - slice of Line objects
//   - error if a position is not found in the index map
func buildLines(
	lineWins calculators.LineWins,
	indexMap map[model.Position]int,
) (lines []Line, err error) {
	lines = make([]Line, 0, len(lineWins))

	for _, lineWin := range lineWins {
		linePayout := lineWin.Payout

		posIndices := make([]int, 0, len(lineWin.Positions))

		for _, pos := range lineWin.Positions {
			if idx, ok := indexMap[pos]; ok {
				posIndices = append(posIndices, idx)
			} else {
				return nil, fmt.Errorf("line %d, pos %v: %w",
					lineWin.Number, pos, ErrPositionNotFound)
			}
		}

		lines = append(lines, Line{
			PosIndices: posIndices,
			Line:       lineWin.Number,
			Symbol:     lineWin.Symbol,
			Payout:     linePayout,
		})
	}

	return lines, nil
}

// comparePos compares two positions for sorting.
// Positions are sorted first by reel index, then by row index.
//
// Parameters:
//   - posA: first position to compare
//   - posB: second position to compare
//
// Returns:
//   - negative if posA < posB, zero if equal, positive if posA > posB
func comparePos(posA, posB model.Position) int {
	if posA.ReelIndex != posB.ReelIndex {
		return cmp.Compare(posA.ReelIndex, posB.ReelIndex)
	}

	return cmp.Compare(posA.RowIndex, posB.RowIndex)
}
