package out

import "fmt"

var (
	// Level is the current LogLevel.
	Level = 0
	// LvlMute is the minimum LogLevel to Mute all messages.
	LvlMute = 3
	// LvlError is the maximum LogLevel to receive error messages.
	LvlError = 2
	// LvlWarning is the maximum LogLevel to receive warnings.
	LvlWarning = 1
	// LvlStatus is the maximum LogLevel to receive status messages.
	LvlStatus = 0
)

// Error prints an error message to the
// standard output.
func Error(msg string) {
	if Level <= LvlError {
		fmt.Println("[ERR]", msg)
	}
}

// Warning prints a warning to the standard
// output.
func Warning(msg string) {
	if Level <= LvlWarning {
		fmt.Println("[WRN]", msg)
	}
}

// Status prints a status message to the
// standard output.
func Status(msg string) {
	if Level <= LvlWarning {
		fmt.Println("[STA]", msg)
	}
}
