package simulator

const (
	// DefaultRTPVersion is the default RTP version used when none is specified via CLI flags.
	DefaultRTPVersion = 96

	// DefaultSpinsCount is the default number of spins in a simulation run.
	DefaultSpinsCount = 100_000_000

	// DefaultPlayersCount is the default number of simulated players per balance group.
	DefaultPlayersCount = 100

	// FilePermissions defines the Unix permission bits used when creating output directories.
	FilePermissions = 0o755

	// DefaultRngType is the default random number generator backend.
	// "mt" uses the Mersenne Twister algorithm
	DefaultRngType = "mt"

	// DefaultBetValue is the default bet amount per spin (in the game's smallest unit).
	DefaultBetValue = 1

	// DefaultSimpleCostMultiplier is the cost multiplier used for the simple game mode.
	// TODO: set a unique value for each game (usually taken from the general config).
	DefaultSimpleCostMultiplier = 5

	// OutputDir is the directory where simulation result files are written.
	OutputDir = "out"

	// DefaultWorkersCount is the number of parallel workers used in simulation.
	// 1 means sequential (no parallelism).
	DefaultWorkersCount = 1
)
