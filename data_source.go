package cali

import "time"

type DataSource interface {
	// Create should save an event in the data store and handle setting the Created and Updated and Id fields
	Create(event Event) (*Event, error)
	// Update uses the given event Id field to update the event in the data store. It is also responsible
	// for updating the Updated field with the current timestamp
	Update(event Event) (*Event, error)
	// Get retrieves a single event from the data store by its Id field. If none is found, it returns nil, nil
	Get(eventId int64) (*Event, error)
	// Query finds a list of events from the data store using the query object to conduct the search
	Query(q Query) ([]*Event, error)
}

type Query struct {
	// Start is an inclusive timestamp and should be compared against the end timestamp of other events (overlap)
	Start time.Time
	// End is an inclusive timestamp and should be compared against the start timestamp of other events (overlap)
	End time.Time
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

// InMemoryDataSource implements the DataSource interface and is useful for a mock data source
type InMemoryDataSource struct {
	events  []*Event
	invites []*Attendance
	curId   int64
}

func (d *InMemoryDataSource) Create(event Event) (*Event, error) {
	return nil, nil
}

func (d *InMemoryDataSource) Update(event Event) (*Event, error) {
	return nil, nil
}

func (d *InMemoryDataSource) Get(eventId int64) (*Event, error) {
	return nil, nil
}

func (d *InMemoryDataSource) Query(q Query) ([]*Event, error) {
	return nil, nil
}

// id generates the next id value
func (d *InMemoryDataSource) id() int64 {
	d.curId++
	return d.curId
}
