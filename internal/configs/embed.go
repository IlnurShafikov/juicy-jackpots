// Package configs provides access to embedded game configuration files.
// The package uses Go's embed feature to include CSV configuration files directly in the binary.
package configs

import (
	"embed"
)

//go:embed csv
var ConfigFS embed.FS

// RootPath defines the root directory path for configuration files within the embedded filesystem.
var RootPath = "csv"
