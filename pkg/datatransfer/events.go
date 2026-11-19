package datatransfer

import "github.com/IlnurShafikov/juicy-jackpots/internal/entities/result/events"

// stepEventProvider is implemented by any game step that exposes its events.
type stepEventProvider interface {
	Events() []events.Event
}

// ExtractAndTransformEvents retrieves events from a step and scales monetary values by bet.
func ExtractAndTransformEvents(step stepEventProvider, betValue uint64) ([]events.Event, error) {
	eventList := step.Events()

	for _, curEvent := range eventList {
		if err := transformSingleEvent(curEvent, betValue); err != nil {
			return nil, err
		}
	}

	return eventList, nil
}

// transformSingleEvent routes a single event to its appropriate transformation handler.
//
// Parameters:
//   - curEvent: the event to transform
//   - betValue: the player's bet value
//
// Returns:
//   - error if the event type is unknown or transformation fails
func transformSingleEvent(curEvent events.Event, betValue uint64) error {
	switch curEvent.EventType() {
	case events.WinsType:
		return transformWinEvent(curEvent, betValue)
	case events.JackpotWinType:
		return transformJackpotWinEvent(curEvent, betValue)
	default:
		return ErrInvalidEventType
	}
}

// transformWinEvent transforms a line wins event by multiplying payouts by bet value.
//
// Parameters:
//   - curEvent: the event to transform (must be *events.LineWinsEvent)
//   - betValue: the player's bet value
//
// Returns:
//   - error if the event is not the correct type
func transformWinEvent(curEvent events.Event, betValue uint64) error {
	winEvent, ok := curEvent.(*events.LineWinsEvent)
	if !ok {
		return ErrInvalidEventType
	}

	for i := range winEvent.Lines {
		winEvent.Lines[i].Payout *= betValue
	}

	return nil
}

// transformJackpotWinEvent scales the jackpot coin value by bet.
//
// Parameters:
//   - curEvent: the event to transform (must be *events.JackpotWinEvent)
//   - betValue: the player's bet value used to scale the jackpot payout
//
// Returns:
//   - [ErrInvalidEventType] if curEvent is not a *events.JackpotWinEvent
func transformJackpotWinEvent(curEvent events.Event, betValue uint64) error {
	jpWinEvent, ok := curEvent.(*events.JackpotWinEvent)
	if !ok {
		return ErrInvalidEventType
	}

	jpWinEvent.Target.Value *= betValue

	return nil
}
