package healthz

import "time"

// Health defines the structure of a healthz check state
type Health struct {
	Healthy   bool      `json:"healthy"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
