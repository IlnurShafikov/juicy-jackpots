package events

// Event is the interface that all game events must implement.
// Events represent occurrences during gameplay that need to be communicated to the client.
type Event interface {
	// EventType returns the type identifier for this event
	EventType() Type
}

// Type represents the category of game event.
type (
	Type string
)

// Event type constants define all possible event types in the game.
// TODO: Add other event types here if needed
const (
	// WinsType represents a win event containing line wins or scatter wins
	WinsType        Type = "wins"
	WinType         Type = "win"
	JpTransformType Type = "jpTransform"
	DoubleWins      Type = "doubleWins"
)
