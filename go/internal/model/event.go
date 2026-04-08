package model

import "time"

// Event represents an event from the client (impression, click, conversion)
type Event struct {
	ID        int64     `json:"id"`
	TestID    string    `json:"test_id"`
	UserID    string    `json:"user_id"`
	Variant   string    `json:"variant"`    // 'A' or 'B'
	Type      string    `json:"event_type"` // 'impression', 'click', 'conversion'
	Value     float64   `json:"value"`      // weight of the event (default 1.0)
	CreatedAt time.Time `json:"created_at"`
}
