// Package parser provides the RoundParser which processes raw game round results
// and accumulates them into a simulation Result for subsequent statistical analysis.
package parser

import (
	resp "github.com/IlnurShafikov/juicy-jackpots/internal/entities/result"
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result/events"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/result"
)

// RoundParser processes round results and accumulates simulation statistics
type RoundParser struct {
	res *result.Result
}

// NewRoundParser creates a new RoundParser with initialized result container
func NewRoundParser(roundCost uint64, simpleCost int) *RoundParser {
	return &RoundParser{
		res: result.New(roundCost, simpleCost),
	}
}

// MergeResult merges a worker's accumulated result into this parser's result.
// Used by the parallel runner after all workers complete.
func (rp *RoundParser) MergeResult(other *result.Result) {
	rp.res.Merge(other)
}

// Parse processes a complete round result and updates simulation statistics
func (rp *RoundParser) Parse(roundResult *resp.Round) error {
	totalBaseWin := rp.parseBasePhase(roundResult.BasePhase)

	roundTotalWin := totalBaseWin

	rp.updateRoundStatistics(roundTotalWin)

	return nil
}

// parseBasePhase processes base game phase results
func (rp *RoundParser) parseBasePhase(baseResult *resp.Base) uint64 {
	baseTotalWin := baseResult.TotalWin()

	rp.res.Base().AddWins(baseTotalWin)
	rp.res.BaseWinRanges.AddWinRange(baseTotalWin)

	for _, eve := range baseResult.Events() {
		switch eve.EventType() {
		case events.JackpotWinType:
			coinWins := parseJackpotWinEvent(eve)
			rp.res.Base().AddCoinWins(coinWins)
		case events.WinsType:
			lineWins := parseLineWinEvent(eve)
			rp.res.Base().AddLineWins(lineWins)
		}
	}

	if baseResult.MaxWinReached() {
		rp.res.AddMaxWinReached()
	}

	return baseTotalWin
}

// updateRoundStatistics updates overall round statistics
func (rp *RoundParser) updateRoundStatistics(totalWin uint64) {
	if totalWin > 0 {
		rp.res.Xi().Add(totalWin)
		rp.res.AddMaxTotalWin(totalWin)
		rp.res.PositiveBets().Add()
	}
}

// GetResult returns the accumulated simulation results
func (rp *RoundParser) GetResult() *result.Result {
	return rp.res
}

// parseLineWinEvent extracts the total payout from a LineWinsEvent.
// Returns 0 if the event is not a *events.LineWinsEvent.
//
// Parameters:
//   - curEvent: game event to inspect
//
// Returns:
//   - sum of all line payouts in the event, or 0 if type assertion fails
func parseLineWinEvent(curEvent events.Event) uint64 {
	winEvent, ok := curEvent.(*events.LineWinsEvent)
	if !ok {
		return 0
	}

	var lineWins uint64

	for _, lineWin := range winEvent.Lines {
		lineWins += lineWin.Payout
	}

	return lineWins
}

func parseJackpotWinEvent(curEvent events.Event) uint64 {
	coinWinEvent, ok := curEvent.(*events.JackpotWinEvent)
	if !ok {
		return 0
	}

	coinWins := coinWinEvent.Target.Value

	return coinWins
}
