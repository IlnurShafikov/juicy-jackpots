package common

// ParentStep represents the parent base phase step data.
// Contains information about the phase that triggered the respins (base or free spins).
type ParentStep struct {
	// Payload contains the parent phase data
	Payload ParentPayload `json:"payload"`
}

type ParentPayload struct {
	// Counter indicates the number of free spins or respins triggered (if any)
	Counter *Counter `json:"counter,omitempty"`

	// Symbols represent the initial symbol board before any transformations
	Symbols []*Symbol `json:"symbols"`
}
