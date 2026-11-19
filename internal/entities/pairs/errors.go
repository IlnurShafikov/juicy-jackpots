package pairs

import "errors"

// ErrInvalidData is returned when the source rows cannot be parsed into pairs.
var ErrInvalidData = errors.New("invalid data")
