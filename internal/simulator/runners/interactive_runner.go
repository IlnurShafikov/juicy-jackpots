//nolint:forbidigo // Prints it's ok for interactive mode
package runners

import (
	"bufio"
	"fmt"
	"os"

	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result"
	"github.com/IlnurShafikov/juicy-jackpots/internal/game"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator"
	"github.com/IlnurShafikov/juicy-jackpots/internal/tools/rngtest"
)

const (
	InputDelimiter = '\n'
)

// Interactive messages
const (
	WelcomeMessageFormat = "Spin limit: %d. Press Enter to run each spin.\n"
	SpinPromptFormat     = "Spin %d/%d. Press Enter to start..."
	CompletionMessage    = "\nAll spins completed!"

	BaseSpin            = "BASE SPIN:"
	FreeSpin            = "FREE SPIN: %d:"
	BoardStartCascade   = "BOARD START CASCADE CYCLE:"
	BoardEndCascade     = "BOARD END OF CASCADE CYCLE:"
	WinsForStepCascades = "WINS FOR STEP CASCADES: %d"
)

// InteractiveRunner handles interactive simulation execution
type InteractiveRunner struct {
	reader *bufio.Reader
}

// NewInteractiveRunner creates a new InteractiveRunner instance
func NewInteractiveRunner() *InteractiveRunner {
	return &InteractiveRunner{
		reader: bufio.NewReader(os.Stdin),
	}
}

// runInteractiveSimulation executes interactive simulation with user input
func (ir *InteractiveRunner) runInteractiveSimulation(
	gameMod *game.Mod,
	config *simulator.SimulationConfig,
) error {
	rng := ir.initializeRNG(config.MainSeed)

	ir.displayWelcomeMessage(config.SpinsCount)

	if err := ir.executeSpinLoop(gameMod, rng, config.SpinsCount); err != nil {
		return fmt.Errorf("run interactive mode: %w",
			simulator.NewSimulationError("execute spin loop", err))
	}

	ir.displayCompletionMessage()

	return nil
}

// initializeRNG creates and seeds the random number generator
func (ir *InteractiveRunner) initializeRNG(mainSeed int64) *rngtest.Mt19937 {
	rng := rngtest.New()
	rng.Reseed(mainSeed)

	return rng
}

// displayWelcomeMessage shows initial instructions to user
func (ir *InteractiveRunner) displayWelcomeMessage(spinsCount int) {
	fmt.Printf(WelcomeMessageFormat, spinsCount)
}

// displayCompletionMessage shows completion message
func (ir *InteractiveRunner) displayCompletionMessage() {
	fmt.Println(CompletionMessage)
}

// executeSpinLoop runs the main spin loop
func (ir *InteractiveRunner) executeSpinLoop(
	gameMod *game.Mod,
	rng *rngtest.Mt19937,
	totalSpins int,
) error {
	for i := range totalSpins {
		if err := ir.processInteractiveSpin(gameMod, rng, i+1, totalSpins); err != nil {
			return fmt.Errorf("process spin %d: %w", i+1, err)
		}
	}

	return nil
}

// processInteractiveSpin handles a single interactive spin
func (ir *InteractiveRunner) processInteractiveSpin(
	gameMod *game.Mod,
	rng *rngtest.Mt19937,
	currentSpin, totalSpins int,
) error {
	if err := ir.waitForUserInput(currentSpin, totalSpins); err != nil {
		return fmt.Errorf("wait for user input: %w", err)
	}

	spinResult, err := ir.executeSpin(gameMod, rng)
	if err != nil {
		return fmt.Errorf("execute spin: %w", err)
	}

	ir.displaySpinResults(spinResult)

	return nil
}

// waitForUserInput prompts user and waits for Enter key
func (ir *InteractiveRunner) waitForUserInput(currentSpin, totalSpins int) error {
	fmt.Printf(SpinPromptFormat, currentSpin, totalSpins)

	_, err := ir.reader.ReadString(InputDelimiter)
	if err != nil {
		return fmt.Errorf("read user input: %w", err)
	}

	return nil
}

// executeSpin performs the actual game spin
func (ir *InteractiveRunner) executeSpin(
	gameMod *game.Mod, rng *rngtest.Mt19937,
) (*result.Round, error) {
	spinResult, err := gameMod.Spin(rng)
	if err != nil {
		return nil, fmt.Errorf("game mod spin: %w", err)
	}

	return spinResult, nil
}

// displaySpinResults prints the results of a spin
func (ir *InteractiveRunner) displaySpinResults(res *result.Round) {
	ir.displayBaseSpinResults(res)

	// if len(res.FreePhase) > 0 {
	// 	ir.displayFreeSpinResults(res)
	// }

	fmt.Println()
}

// displayBaseSpinResults displays base game spin results
func (ir *InteractiveRunner) displayBaseSpinResults(res *result.Round) {
	// TODO implement board printer here depending on your types

	// fmt.Println(BaseSpin)
	//
	// fmt.Println(BoardStartCascade)
	// res.BasePhase.BeforeBoard().PrintSymbols()
	//
	// fmt.Printf(WinsForStepCascades+"\n", res.BasePhase.TotalWin())
}

// displayFreeSpinResults displays free spin results
func (ir *InteractiveRunner) displayFreeSpinResults(res *result.Round) { //nolint:unused // TODO
	// TODO implement board printer here depending on your types

	// for i, freeSpin := range res.FreePhase {
	// 	fmt.Printf(FreeSpin, i+1)
	//
	// 	fmt.Println(BoardStartCascade)
	// 	freeSpin.BeforeBoard().PrintSymbols()
	//
	// 	fmt.Printf(WinsForStepCascades+"\n", freeSpin.StepWin())
	// }
}
