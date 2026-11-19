package tools

import (
	"path"
)

// PathBuilder represents a hierarchical path to game configuration files.
// It provides a fluent API for constructing paths by appending successive
// path segments (RTP version, mod type, phase, etc.).
type PathBuilder struct {
	// rootCfgPath is the path built so far.
	rootCfgPath string
}

// NewPathBuilder creates a new PathBuilder rooted at the given path.
//
// Parameters:
//   - rootPath: the base directory from which all paths are built
//
// Returns:
//   - pointer to a new PathBuilder instance
func NewPathBuilder(rootPath string) *PathBuilder {
	return &PathBuilder{rootCfgPath: rootPath}
}

// String returns the current path as a plain string.
func (gp *PathBuilder) String() string {
	return gp.rootCfgPath
}

// Last returns only the final segment of the current path (the base name).
func (gp *PathBuilder) Last() string {
	_, last := path.Split(gp.rootCfgPath)

	return last
}

// RTP appends the given RTP version string to the current path
// and returns a new PathBuilder pointing to the resulting directory.
//
// Parameters:
//   - rtpVersion: string representation of the RTP version, e.g. "96"
func (gp *PathBuilder) RTP(rtpVersion string) *PathBuilder {
	return &PathBuilder{rootCfgPath: path.Join(gp.rootCfgPath, rtpVersion)}
}

// Mod appends the given mod type to the current path
// and returns a new PathBuilder pointing to the resulting directory.
//
// Parameters:
//   - modType: mod directory name, e.g. "simple"
func (gp *PathBuilder) Mod(modType string) *PathBuilder {
	return &PathBuilder{rootCfgPath: path.Join(gp.rootCfgPath, modType)}
}

// Phase appends the given game phase name to the current path
// and returns a new PathBuilder pointing to the resulting directory.
//
// Parameters:
//   - phase: phase directory name, e.g. "base" or "free"
func (gp *PathBuilder) Phase(phase string) *PathBuilder {
	return &PathBuilder{rootCfgPath: path.Join(gp.rootCfgPath, phase)}
}

// CfgDir appends the given configuration directory name to the current path
// and returns a new PathBuilder pointing to the resulting directory.
//
// Parameters:
//   - cfgDir: configuration sub-directory name
func (gp *PathBuilder) CfgDir(cfgDir string) *PathBuilder {
	return &PathBuilder{rootCfgPath: path.Join(gp.rootCfgPath, cfgDir)}
}

// ReelSet appends the given reel-set (scenario) name to the current path
// and returns a new PathBuilder pointing to the resulting directory.
//
// Parameters:
//   - reelSetName: reel set directory name
func (gp *PathBuilder) ReelSet(reelSetName string) *PathBuilder {
	return &PathBuilder{rootCfgPath: path.Join(gp.rootCfgPath, reelSetName)}
}

// Dir appends the given directory path segment to the current path
// and returns a new PathBuilder pointing to the resulting directory.
//
// Parameters:
//   - dirPath: directory name or relative path segment to append
func (gp *PathBuilder) Dir(dirPath string) *PathBuilder {
	return &PathBuilder{rootCfgPath: path.Join(gp.rootCfgPath, dirPath)}
}

// File returns the full path to a specific configuration file
// by joining the current path with the given file name.
//
// Parameters:
//   - configName: the configuration file name (without directory prefix)
//
// Returns:
//   - full path string suitable for passing to an fs.FS reader
func (gp *PathBuilder) File(configName string) string {
	return path.Join(gp.rootCfgPath, configName)
}
