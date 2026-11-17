package main

import (
	"log"

	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator"
	"github.com/IlnurShafikov/juicy-jackpots/internal/simulator/runners"
)

func main() {
	cfg := simulator.ParseConfigs()

	simManager := runners.NewSimulationManager()
	if err := simManager.RunSimulation(cfg); err != nil {
		//nolint:forbidigo // it's ok for simulator
		log.Fatal(err)
	}
}
