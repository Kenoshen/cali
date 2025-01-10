package cali

import "time"

type DataStore interface {
	// Create should save an event in the data store and handle setting the Created and Updated and Id fields
	Create(event Event) (*Event, error)
	// Update uses the given event Id field to update the event in the data store. It is also responsible
	// for updating the Updated field with the current timestamp
	Update(details Details) (*Event, error)
	// Get retrieves a single event from the data store by its Id field. If none is found, it returns nil, nil
	Get(eventId int64) (*Event, error)
	// Query finds a list of events from the data store using the query object to conduct the search
	Query(q Query) ([]*Event, error)

	// AddAttendance adds a new attendance record to the data store and handles
	// setting the Created and Updated fields
	AddAttendance(attendance Attendance) (*Attendance, error)
	// UpdateAttendance uses the EventId and UserId to update the attendance in
	// the data store. Also responsible for updating the Updated field with the
	// current timestamp.
	UpdateAttendance(attendance Attendance) (*Attendance, error)
	// GetAttendance retrieves a single Attendance by the EventId and UserId fields.
	// If none is found, it returns nil, nil
	GetAttendance(eventId int64, userId int64) (*Attendance, error)
}

// InMemoryDataStore implements the DataStore interface and is useful for a mock data source
type InMemoryDataStore struct {
	events  []*Event
	invites []*Attendance
	curId   int64
}

func (d *InMemoryDataStore) Create(event Event) (*Event, error) {
	err := Validate(event)
	if err != nil {
		return nil, err
	}
	event.Id = d.id()
	event.Created = time.Now()
	event.Updated = event.Created

	_, err = d.AddAttendance(Attendance{
		EventId:    event.Id,
		UserId:     event.OwnerId,
		Status:     AttendanceStatusConfirmed,
		Permission: PermissionOwner,
	})
	if err != nil {
		return nil, err
	}

	d.events = append(d.events, &event)

	return &event, nil
}

func (d *InMemoryDataStore) Update(details Details) (*Event, error) {
	for _, other := range d.events {
		if other.Id == details.Id {
			other.Title = details.Title
			other.Description = details.Description
			other.Url = details.Url
			other.Status = details.Status
			other.IsAllDay = details.IsAllDay
			other.Zone = details.Zone
			other.StartDay = details.StartDay
			other.StartTime = details.StartTime
			other.EndDay = details.EndDay
			other.EndTime = details.EndTime
			other.Updated = time.Now()
			return other, nil
		}
	}
	return nil, nil
}

func (d *InMemoryDataStore) Get(eventId int64) (*Event, error) {
	for _, event := range d.events {
		if event.Id == eventId {
			return event, nil
		}
	}
	return nil, nil
}

func (d *InMemoryDataStore) Query(q Query) ([]*Event, error) {
	var result []*Event

	for _, event := range d.events {
		if q.Matches(event) {
			result = append(result, event)
		}
	}

	return result, nil
}

func (d *InMemoryDataStore) AddAttendance(a Attendance) (*Attendance, error) {
	a.Created = time.Now()
	a.Updated = a.Created
	err := ValidateAttendance(a)
	if err != nil {
		return nil, err
	}
	d.invites = append(d.invites, &a)
	return &a, nil
}

func (d *InMemoryDataStore) UpdateAttendance(a Attendance) (*Attendance, error) {
	for _, attendance := range d.invites {
		if attendance.EventId == a.EventId && attendance.UserId == a.UserId {
			attendance.Status = a.Status
			attendance.Permission = a.Permission
			attendance.Updated = time.Now()
			return attendance, nil
		}
	}
	return nil, nil
}

func (d *InMemoryDataStore) GetAttendance(eventId int64, userId int64) (*Attendance, error) {
	for _, attendance := range d.invites {
		if attendance.EventId == eventId && attendance.UserId == userId {
			return attendance, nil
		}
	}
	return nil, nil
}

// id generates the next id value
func (d *InMemoryDataStore) id() int64 {
	d.curId++
	return d.curId
}
