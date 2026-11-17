// Package reader provides functionality for reading and parsing game configuration files.
// It handles loading configurations from CSV files and converting them into usable game structures.
package readers

import (
	"errors"
)

// ErrReadGameConfigs is returned when an error occurs while reading game configuration files.
var ErrReadGameConfigs = errors.New("error reading game configs")
