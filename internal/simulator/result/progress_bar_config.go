package result

import "time"

// Configs holds the settings of the simulation progress bar.
type Configs struct {
	// updateMask gates renders: the bar is redrawn only when current&updateMask == 0.
	updateMask int
	// refreshRate is the minimum time between two renders.
	refreshRate time.Duration
	// barWidth is the width of the bar in characters.
	barWidth int
	// total is the total number of iterations to track.
	total int
	// showProgressBar enables the graphical bar; otherwise only the text is printed.
	showProgressBar bool
}

// DefaultConfigs returns the default progress bar settings for the given total.
//
// Parameters:
//   - total: total number of iterations to track
//
// Returns:
//   - progress bar settings with the default mask, refresh rate and width
func DefaultConfigs(total int) Configs {
	const (
		defaultMask  = 1023
		defaultRate  = time.Millisecond * 200
		defaultWidth = 50
		defaultShow  = false
	)

	return NewConfigs(defaultMask, defaultWidth, total, defaultRate, defaultShow)
}

// NewConfigs creates progress bar settings.
//
// Parameters:
//   - updateMask: bitmask that gates renders (render when current&updateMask == 0)
//   - barWidth: width of the bar in characters
//   - total: total number of iterations to track
//   - refreshRate: minimum time between two renders
//   - showProgressBar: whether to draw the graphical bar
//
// Returns:
//   - progress bar settings
func NewConfigs(
	updateMask, barWidth, total int,
	refreshRate time.Duration,
	showProgressBar bool,
) Configs {
	return Configs{
		updateMask:      updateMask,
		refreshRate:     refreshRate,
		barWidth:        barWidth,
		total:           total,
		showProgressBar: showProgressBar,
	}
}
