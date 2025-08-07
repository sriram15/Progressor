package events

import "time"

const (
	// CardStoppedTopic is the topic for when a card is stopped.
	CardStoppedTopic = "card:stopped"
	// CardStartedTopic is the topic for when a card is started.
	CardStartedTopic = "card:started"
)

// CardStoppedEvent is the data for the event when a card is stopped.
type CardStoppedEvent struct {
	CardID    int64
	ProjectID int64
	UserID    int64
	TimeSpent time.Duration
	StoppedAt time.Time
}

// CardStartedEvent is the data for the event when a card is started.
type CardStartedEvent struct {
	CardID    int64
	ProjectID int64
	UserID    int64
	StartedAt time.Time
}
