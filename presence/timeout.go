package presence

import "time"

// Timeout describes the the end of presence
type Timeout struct {
	Sid              string
	Code             string
	Start            time.Time
	LastNotification time.Time
	End              time.Time
	Duration         time.Duration
	Wait             time.Duration
}
