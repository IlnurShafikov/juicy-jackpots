// Command simulator runs the slot game simulation and writes the statistics reports.
package main

import (
	"log"

	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/runners"
)

// main parses the CLI flags and runs the simulation.
// A failed simulation terminates the program with a non-zero exit code.
func main() {
	cfg := simulator.ParseConfigs()

	simManager := runners.NewSimulationManager()
	if err := simManager.RunSimulation(cfg); err != nil {
		//nolint:forbidigo // it's ok for simulator
		log.Fatal(err)
	}
}
