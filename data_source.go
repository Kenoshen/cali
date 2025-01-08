package cali

import "time"

type DataSource interface {
	Create(event Event) (*Event, error)
	Update(event Event) (*Event, error)
	Get(eventId int64) (*Event, error)
	Query(q Query) ([]*Event, error)
}

type Query struct {
	Start time.Time
	End   time.Time
	// EventIds is a list of specific events that you want to query
	EventIds []int64
	// UserIds is a check if the user has an attendance record for the event that is not declined
	UserIds []int64
	// EventTypes is a check if the event has a specific event type
	EventTypes []EventType
	// SourceIds is an OR check on the source ids
	SourceIds []int64
	// Text is an OR search for specific words
	Text []string
	// Statuses is an OR search for specific statuses
	Statuses []Status
}
