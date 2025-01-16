package cali

import (
	"time"
)

type DataStore interface {
	// Create should save an event in the data store and handle setting the Created and Updated and Id fields
	Create(event Event) (*Event, error)
	// SetTime updates the time values for a specific event
	SetTime(eventId int64, startDay, startTime, endDay, endTime, zone string, isAllDay bool) error
	// SetStatus applies the given status to the event. If the event already has the status it returns nil
	SetStatus(eventId int64, status Status) error
	// SetTitle updates the event with the given title
	SetTitle(eventId int64, title string) error
	// SetDescription updates the event with the given description
	SetDescription(eventId int64, description *string) error
	// SetUrl updates the event with the url value
	SetUrl(eventId int64, url *string) error
	// SetUserData updates the event with the user data
	SetUserData(eventId int64, userData map[string]interface{}) error
	// Get retrieves a single event from the data store by its Id field. If none is found, it returns nil, nil
	Get(eventId int64) (*Event, error)
	// Query finds a list of events from the data store using the query object to conduct the search
	Query(q Query) ([]*Event, error)

	// AddInvite adds a new invite record to the data store and handles
	// setting the Created and Updated fields
	AddInvite(invite Invite) (*Invite, error)
	// SetInviteStatus uses the EventId and UserId to update the status of the invite and updates the Updated date too
	SetInviteStatus(eventId, userId int64, status InviteStatus) error
	// SetInvitePermissions uses the EventId and UserId to update the permissions of the invite and updates the Updated date too
	SetInvitePermissions(eventId, userId int64, permissions Permission) error
	// GetInvite retrieves a single Invite by the EventId and UserId fields.
	// If none is found, it returns nil, nil
	GetInvite(eventId, userId int64) (*Invite, error)
}

// InMemoryDataStore implements the DataStore interface and is useful for a mock data source
type InMemoryDataStore struct {
	events  []*Event
	invites []*Invite
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

	// if the event is a repeating event, but doesn't have the ParentId
	// field set, then this must be the first event of the repeat and
	// should also have its own Id as the ParentId
	if event.IsRepeating && event.ParentId == nil {
		event.ParentId = &event.Id
	}

	_, err = d.AddInvite(Invite{
		EventId:    event.Id,
		UserId:     event.OwnerId,
		Status:     InviteStatusConfirmed,
		Permission: PermissionOwner,
	})
	if err != nil {
		return nil, err
	}

	d.events = append(d.events, &event)
	return &event, nil
}

func (d *InMemoryDataStore) SetTime(eventId int64, startDay, startTime, endDay, endTime, zone string, isAllDay bool) error {
	if err := ValidateDayTimeValues(startDay, startTime, endDay, endTime, zone, isAllDay); err != nil {
		return err
	}

	for _, other := range d.events {
		if other.Id == eventId {
			other.StartDay = startDay
			other.StartTime = startTime
			other.EndDay = endDay
			other.EndTime = endTime
			other.IsAllDay = isAllDay
			other.Zone = zone
			return nil
		}
	}
	return ErrorEventNotFound
}

func (d *InMemoryDataStore) SetStatus(eventId int64, status Status) error {
	if !ValidStatus(status) {
		return ErrorInvalidStatus
	}

	for _, other := range d.events {
		if other.Id == eventId {
			other.Status = status
			return nil
		}
	}
	return ErrorEventNotFound
}

func (d *InMemoryDataStore) SetTitle(eventId int64, title string) error {
	for _, other := range d.events {
		if other.Id == eventId {
			other.Title = title
			return nil
		}
	}
	return ErrorEventNotFound
}

func (d *InMemoryDataStore) SetDescription(eventId int64, description *string) error {
	for _, other := range d.events {
		if other.Id == eventId {
			other.Description = description
			return nil
		}
	}
	return ErrorEventNotFound
}

func (d *InMemoryDataStore) SetUrl(eventId int64, url *string) error {
	for _, other := range d.events {
		if other.Id == eventId {
			other.Url = url
			return nil
		}
	}
	return ErrorEventNotFound
}

func (d *InMemoryDataStore) SetUserData(eventId int64, userData map[string]interface{}) error {
	for _, other := range d.events {
		if other.Id == eventId {
			other.UserData = userData
			return nil
		}
	}
	return ErrorEventNotFound
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

func (d *InMemoryDataStore) AddInvite(a Invite) (*Invite, error) {
	a.Created = time.Now()
	a.Updated = a.Created
	err := ValidateInvite(a)
	if err != nil {
		return nil, err
	}
	d.invites = append(d.invites, &a)
	return &a, nil
}

func (d *InMemoryDataStore) SetInviteStatus(eventId, userId int64, status InviteStatus) error {
	for _, invite := range d.invites {
		if invite.EventId == eventId && invite.UserId == userId {
			invite.Status = status
			invite.Updated = time.Now()
			return nil
		}
	}
	return ErrorInviteNotFound
}

func (d *InMemoryDataStore) SetInvitePermissions(eventId, userId int64, permissions Permission) error {
	for _, invite := range d.invites {
		if invite.EventId == eventId && invite.UserId == userId {
			invite.Permission = permissions
			invite.Updated = time.Now()
			return nil
		}
	}
	return ErrorInviteNotFound
}

func (d *InMemoryDataStore) GetInvite(eventId int64, userId int64) (*Invite, error) {
	for _, invite := range d.invites {
		if invite.EventId == eventId && invite.UserId == userId {
			return invite, nil
		}
	}
	return nil, nil
}

// id generates the next id value
func (d *InMemoryDataStore) id() int64 {
	d.curId++
	return d.curId
}
