package collectors

import "github.com/IlnurShafikov/juicy-jackpots/internal/entities/result/events"

// defaultCapacity defines the initial capacity for the events slice to optimize memory allocation.
const defaultCapacity = 2

// EventsCollector is responsible for collecting game events during a spin.
// It stores events such as wins, triggers, and other game mechanics that occur during gameplay.
type EventsCollector struct {
	// events stores the collected events
	events []events.Event
	// cap defines the initial capacity for the events slice
	cap int
}

// NewEvents creates a new EventsCollector with default capacity.
//
// Returns:
//   - Pointer to a new EventsCollector instance
func NewEvents() *EventsCollector {
	return &EventsCollector{
		events: nil,
		cap:    defaultCapacity,
	}
}

// Add appends a new event to the collector.
// Initializes the events slice on the first addition with predefined capacity.
//
// Parameters:
//   - Event: the event to be added to the collection
func (e *EventsCollector) Add(event events.Event) {
	if len(e.events) == 0 {
		e.events = make([]events.Event, 0, e.cap)
	}

	e.events = append(e.events, event)
}

// Get returns all collected events.
//
// Returns:
//   - slice containing all collected events
func (e *EventsCollector) Get() []events.Event {
	return e.events
}
