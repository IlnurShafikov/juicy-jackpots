# Juicy Jackpots

A classic fruit slot written in Go: 5 reels, 3 rows, 5 fixed paylines and
**Jackpot Coins** that are drawn on every spin. The repository contains the game math,
a deterministic round engine, a client data-transfer layer and a high-throughput Monte Carlo
simulator that is used to verify RTP and volatility.

[![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white)](go.mod)
[![golangci-lint](https://img.shields.io/badge/lint-golangci--lint%20v2-brightgreen)](.golangci.yaml)

> Juicy Jackpots is a working title. The game is built on top of the in-house
> [`slot-engine`](https://github.com/IlnurShafikov/slot-engine) library (RNG interfaces,
> reel generator, lines evaluator, win limiter).

## Game at a glance

| | |
|---|---|
| Layout | 5 reels x 3 rows, 5 fixed paylines paying from left to right |
| Symbols | High1, High2, Med1, Med2, Low1, Low2, Low3, Low4 and Wild |
| Wild | Reels 2-4 only, substitutes for every paying symbol, has no payout of its own |
| Lines | Only the longest match on a line pays; Low4 pays from 2 reels, other symbols from 3 |
| Jackpot Coins | Drawn on **every spin**, independent of the reels and lines (see below) |
| Features | No scatter and no free spins: the whole game is the base spin plus the coins |
| Max win | **1000x** the total bet, applied once to the sum of the line win and the jackpot prize |
| RTP versions | 90, 92, 94 and 96 (same strips, paytable and prizes, different weights) |

### Jackpot Coins

On every spin 0-5 coins are drawn by weights (`jp_coin_count.csv`). Their reels and rows are random and
only visual: coins never replace symbols and never change a line win. The prize depends on the number of
coins only:

| Coins | Prize | Total bets |
|:-:|:-:|:-:|
| 0-1 | none | 0 |
| 2 | MINI | 10x |
| 3 | MINOR | 20x |
| 4 | MAJOR | 100x |
| 5 | GRAND | 1000x |

At most one prize is paid per spin. The only interaction with the lines is the win cap: a spin with GRAND
always pays exactly 1000x, so its line win is lost to the cap.

The full specification (reel selection, stop positions, line evaluation, caps) is in
[`docs/Juicy Jackpots algorithm.md`](docs/Juicy%20Jackpots%20algorithm.md); the paytable and the math sheet are in
[`docs/Juicy Jackpots Par Sheet.xlsx`](docs/Juicy%20Jackpots%20Par%20Sheet.xlsx).

## Simulation results

Every RTP version was verified with **1 000 000 000 rounds** (4 workers). Full reports with win-range tables
are stored in [`docs/simulations`](docs/simulations).

| RTP version | Expected RTP | Simulated RTP | Delta | By lines | By coins | Hit rate | Volatility index | Max win hits |
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
| 90 | 90.35% | 90.35% | 0.00 pp | 75.53% | 14.82% | 26.93% | 8.89 | 8 772 |
| 92 | 92.35% | 92.34% | -0.01 pp | 76.30% | 16.04% | 27.16% | 9.11 | 9 378 |
| 94 | 94.35% | 94.36% | +0.01 pp | 77.74% | 16.62% | 27.49% | 9.29 | 9 824 |
| 96 | 96.35% | 96.39% | +0.04 pp | 79.40% | 16.99% | 27.87% | 9.41 | 10 017 |

The coins contribute roughly 15-17% of the RTP, so the game combines frequent small line wins with a
jackpot layer. One billion rounds take about 13-16 minutes with 4 workers (roughly 1.1-1.3 million rounds per
second). The expected RTP comes from `internal/configs/csv/<rtp>/exp_rtp.csv`.

## Architecture

```
pkg/manager            public entry point: GameManager (enter game, play a round)
pkg/datatransfer       converts an internal round into client steps with events
internal/game          game mod: runs the base phase under one win limiter
  phase/base           base spin: board, line wins, jackpot coins
internal/components    filler (jackpot coins), manager (coin win), collectors, limiter
internal/entities      results, coins, events, game rules and constants
internal/configs       CSV configs (embedded into the binary) and their readers
internal/simulator     Monte Carlo simulator: runners, statistics, report writers
```

Request flow:

```
GameManager.GetRoundResult(mod, rtp, seed, bet)
   |  reseed RNG with the round seed
   v
game.Mod.Spin ──> base.SpinPhase.Spin
                     |-- pick reel group by weights, generate the board
                     |-- filler.JpCoins.Fill        draw coins, resolve prize by coin count
                     |-- calculators.Lines          5 paylines, wild substitution
                     '-- manager.JpCoinWin          add the prize through the same WinLimiter
   v
datatransfer.Transfer ──> client step: symbols, line wins event, jackpot win event
```

Design decisions worth mentioning:

- **Deterministic rounds.** A round is fully defined by `(rtp version, mod, seed)`, so every payout is
  reproducible and verifiable.
- **Configs are data, not code.** Strips, weights, paytable, win lines, prizes and caps are CSV files embedded
  with `go:embed`; RTP versions differ only by weights. No file system access at runtime.
- **One win limiter per round.** Line wins and the jackpot prize go through the same limiter, so the 1000x cap
  is enforced consistently and the order of the two does not matter.
- **Independent random events.** Lines and coins use separate draws, which makes each part of the RTP
  separately measurable (see the "By lines" and "By coins" columns).
- **Client DTOs are separated from the engine.** Payouts in events are scaled by the bet only in
  `datatransfer`; the engine always works in base units.

## Getting started

Requirements: Go 1.27+ and access to the `slot-engine` module
(`GOPRIVATE=github.com/IlnurShafikov/*` while the module is private).

```bash
git clone https://github.com/IlnurShafikov/juicy-jackpots.git
cd juicy-jackpots
go test ./...
```

### Using the game manager

```go
gm, err := manager.NewGameManager()
if err != nil {
    return err
}

// Play one round in the simple mode, RTP 96, bet 1, with an explicit seed.
res, err := gm.GetRoundResult(ctx, rules.SimpleModID, rules.EmptyFeatureID, 96, seed, 1)
if err != nil {
    return err
}

for _, s := range res.Steps {
    fmt.Println(s.Type, s.StepWin, s.TotalWin, s.MaxWinReached)
}
```

### Running the simulator

```bash
# 10 million rounds, RTP 96 version, 4 parallel workers
go run ./internal/simulator/cmd -rtp 96 -spins 10000000 -workers 4

# step-by-step interactive mode: press Enter to run every spin and see the boards
go run ./internal/simulator/cmd -play -spins 20

# player balance simulation: how many spins players survive with given balances
go run ./internal/simulator/cmd -playerMod -balance 100,500,1000 -playersCount 100 -spins 100000
```

Reports are written into the `out/` directory (ignored by git).

| Flag | Default | Description |
|---|---|---|
| `-rtp` | 90 | RTP version to simulate (90, 92, 94, 96) |
| `-spins` | 1 000 000 000 | Number of rounds |
| `-workers` | 4 | Parallel workers; `1` runs sequentially |
| `-seed` | fixed | Main seed. Results are reproducible for the same seed and worker count |
| `-mod` / `-feature` | `simple` / empty | Game mode and optional feature |
| `-rng` | `mt` | RNG backend (Mersenne Twister) |
| `-play` | false | Interactive mode |
| `-playerMod` | false | Player balance simulation (`-balance`, `-playersCount`) |
| `-bet`, `-costMultiplier` | 1, 5 | Bet value and cost multiplier of the simple mode |
| `-description` | empty | Free-form label stored in the report |

Each report contains the total RTP split into line wins and coin wins, the number of jackpot triggers,
standard deviation, volatility index, hit rate, max win statistics and a win-range histogram (count,
frequency and RTP share).

## Code quality

- `golangci-lint` v2 with a strict configuration ([`.golangci.yaml`](.golangci.yaml)): bugs, errors,
  style, complexity and performance linters, zero open findings.
- Every package, type, function and struct field has a godoc comment (`Parameters` / `Returns` blocks
  for non-trivial functions).
- Errors are wrapped with context and compared with `errors.Is`; no panics in the game path.
- Safe numeric conversions (`internal/tools/cast`) are covered by table-driven tests.
- Integration test: rounds are played through the manager and the total win sent to the client is
  compared with the internal result.

```bash
golangci-lint run ./...
go test -race ./...
```

## Project layout

```
.
├── docs/                  algorithm, math sheet, simulation reports
├── internal/
│   ├── components/        filler (coins), manager (coin win), collectors, limiter
│   ├── configs/           embedded CSV configs and readers
│   ├── entities/          results, coins, events, rules
│   ├── game/              game mod and base phase
│   ├── simulator/         CLI simulator, runners, reports
│   └── tools/             cast, path builder, test RNG
└── pkg/
    ├── datatransfer/      client DTOs and transformation
    └── manager/           public game manager
```
