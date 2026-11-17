package common

import (
	"fmt"
	"strconv"

	"github.com/IlnurShafikov/juicy-jackpots/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/reels"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

// ReadReels reads reel strip data for every trigger group from the phase CSV directory.
// Each group ID corresponds to a separate CSV file (group_N) that defines
// the weighted symbol sequence for that distribution group.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//   - phasePath: path builder pointing to the game phase configuration directory
//   - groupIds: list of trigger group IDs whose reel files should be loaded
//
// Returns:
//   - slice of Reels, one entry per group ID in the same order as groupIds
//   - error if any reel file cannot be read or parsed
func ReadReels(
	reader *reader.Adapter,
	phasePath *tools.PathBuilder,
	groupIds []int,
) ([]*reels.Reels, error) {
	const reelDir = "reels"

	reelPath := phasePath.Dir(reelDir)
	fileName := "group_"

	reelsGroup := make([]*reels.Reels, 0, len(groupIds))

	for _, groupID := range groupIds {
		filePrefix := fileName + strconv.Itoa(groupID)

		reel, err := readReelGroup(reader, reelPath.File(filePrefix), groupID)
		if err != nil {
			return nil, fmt.Errorf("read reels: %w", err)
		}

		reelsGroup = append(reelsGroup, reel)
	}

	return reelsGroup, nil
}

// readReelGroup reads a single reel-strip CSV file and parses it into a Reels
// value tagged with the given group ID.
//
// Parameters:
//   - reader: CSV file adapter for reading configuration files
//   - prefix: file path within the embedded FS to the reel group CSV
//   - groupID: numeric identifier assigned to this reel strip group
//
// Returns:
//   - pointer to Reels containing the parsed symbol strips for the group
//   - error if the file cannot be read or the strip data is invalid
func readReelGroup(
	reader *reader.Adapter,
	prefix string,
	groupID int,
) (*reels.Reels, error) {
	data, err := reader.ReadCSV(prefix)
	if err != nil {
		return nil, fmt.Errorf("read group reel: %w", err)
	}

	resp, err := reels.Parse(groupID, data)
	if err != nil {
		return nil, fmt.Errorf("parse reels group id[%d]: %w", groupID, err)
	}

	return resp, nil
}
