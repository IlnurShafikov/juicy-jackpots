package result

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

// Progress prints the simulation progress, speed and remaining time to stdout.
type Progress struct {
	// progressInformation holds the timing data of the bar.
	progressInformation
	// Configs holds the bar settings.
	Configs
	// current is the number of completed iterations.
	current int
	// speed is the number of iterations per second.
	speed float64
}

// progressInformation holds the timing data of the progress bar.
type progressInformation struct {
	// startTime is the moment the bar was created.
	startTime time.Time
	// lastUpdate is the moment of the last render.
	lastUpdate time.Time
	// nextCheckTime is the earliest moment of the next render.
	nextCheckTime time.Time
	// elapsedTime is the time since startTime at the last render.
	elapsedTime time.Duration
	// remainingTime is the estimated time until completion.
	remainingTime time.Duration
}

// NewProgressBar creates a Progress that starts counting time immediately.
//
// Parameters:
//   - cfg: progress bar settings
//
// Returns:
//   - pointer to a new Progress
func NewProgressBar(cfg Configs) *Progress {
	now := time.Now()

	return &Progress{
		Configs: cfg,
		progressInformation: progressInformation{
			startTime:     now,
			lastUpdate:    now,
			nextCheckTime: now.Add(cfg.refreshRate),
			elapsedTime:   0,
			remainingTime: 0,
		},
		current: 0,
		speed:   0,
	}
}

// Update records the current iteration and redraws the bar when the update mask
// and the refresh rate allow it. The last iteration is always drawn.
//
// Parameters:
//   - current: number of completed iterations
func (pb *Progress) Update(current int) {
	pb.current = current

	now := time.Now()
	if (current&pb.updateMask) == 0 || current == pb.total {
		if now.After(pb.nextCheckTime) || current == pb.total {
			pb.updateStats(now)
			pb.render()
			pb.lastUpdate = now
			pb.nextCheckTime = now.Add(pb.refreshRate)
		}
	}
}

// updateStats recalculates the elapsed time, speed and remaining time.
//
// Parameters:
//   - now: current time
func (pb *Progress) updateStats(now time.Time) {
	pb.elapsedTime = now.Sub(pb.startTime)

	pb.speed = float64(pb.current) / pb.elapsedTime.Seconds()
	if pb.speed > 0 {
		remainingCount := pb.total - pb.current
		pb.remainingTime = time.Duration(float64(remainingCount) / pb.speed * float64(time.Second))
	}
}

// render prints the progress line to stdout, overwriting the previous one.
func (pb *Progress) render() {
	width := getTerminalWidth()
	ratio := float64(pb.current) / float64(pb.total)
	percent := int(ratio * 100) //nolint:mnd // percentage conversion: 100 is a domain constant

	infoText := fmt.Sprintf("| %3d%% | %d/%d | %.2f it/s | %v remaining",
		percent, pb.current, pb.total, pb.speed, pb.remainingTime.Round(time.Second))

	if pb.showProgressBar { //nolint:nestif // nested conditions for bar width calculation
		availableWidth := width - len(infoText) - 3 //nolint:mnd // 3 for "[", "]", and space
		if availableWidth < 10 {                    //nolint:mnd // minimum bar width
			availableWidth = 10 // Минимальная ширина прогресс-бара
		}

		barLength := int(ratio * float64(availableWidth))
		bar := strings.Repeat("█", barLength) + strings.Repeat("░", availableWidth-barLength)

		fmt.Printf("\r[%s]%s", bar, infoText)
	} else {
		fmt.Printf("\r%s", infoText)
	}
}

// Finish draws the final state of the bar and moves to a new line.
func (pb *Progress) Finish() {
	pb.Update(pb.total)
	fmt.Println() //nolint:forbidigo // terminal newline after progress bar completion
}

// getTerminalWidth returns the width of the terminal or a default width
// when it cannot be detected.
func getTerminalWidth() int {
	const defaultWidth = 50

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		return width
	}

	return defaultWidth
}
