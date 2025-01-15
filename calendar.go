package cali

import (
	"fmt"
	"time"
)

// Calendar is an object that can interact with a data store
// (usually a database or cache of some kind) to create, update,
// and retrieve calendar events. It handles the logic of dealing
// with repating events as well as validating events and
// invitations to make sure they don't have conflicting values.
// It tries to be as stateless as possible.
type Calendar struct {
	// dataStore is the implementation of the data store that the
	// event and invitation data will be stored in
	dataStore DataStore
}

// NewCalendar creates a new calendar with the given data store
func NewCalendar(dataStore DataStore) *Calendar {
	c := &Calendar{
		dataStore: dataStore,
	}
	return c
}

// Get grabs a single event by id
func (c *Calendar) Get(eventId int64) (*Event, error) {
	return c.dataStore.Get(eventId)
}

// Query collects a list of events using the provided query parameters
func (c *Calendar) Query(q Query) ([]*Event, error) {
	return c.dataStore.Query(q)
}

// Create an event with the given values. Created and Updated fields will be set automatically. Repeating events will also be created automatically.
func (c *Calendar) Create(e Event) (*Event, int64, error) {
	if err := Validate(e); err != nil {
		return nil, 0, err
	}

	if !e.IsRepeating {
		newEvent, err := c.dataStore.Create(e)
		var count int64 = 0
		if newEvent != nil {
			count++
		}
		return newEvent, count, err
	}

	events, err := GenerateRepeatEvents(e)
	if err != nil {
		return nil, 0, err
	}

	if events == nil || len(events) == 0 {
		return nil, 0, ErrorEmptyRepeatingEvents
	}

	var results []*Event
	var count int64 = 0
	for _, event := range events {
		newEvent, err := c.dataStore.Create(*event)
		if err != nil {
			return nil, 0, err
		}
		if newEvent != nil {
			count++
		}
		results = append(results, newEvent)
	}

	return events[0], count, nil
}

// UpdateTime changes the day and time values of the event
func (c *Calendar) UpdateTime(eventId int64, startDay string, startTime string, endDay string, endTime string, zone string, isAllDay bool, editType RepeatEditType) error {
	if err := ValidTimes(startDay, startTime, endDay, endTime, zone, isAllDay); err != nil {
		return err
	}
	return c.dataStore.SetTime(eventId, startDay, startTime, endDay, endTime, zone, isAllDay)
}

// Cancel sets the status of the event to StatusCanceled
func (c *Calendar) Cancel(eventId int64, editType RepeatEditType) error {
	return c.dataStore.SetStatus(eventId, StatusCanceled)
}

// Remove sets the status of the event to StatusRemoved (we never delete things here)
func (c *Calendar) Remove(eventId int64, editType RepeatEditType) error {
	return c.dataStore.SetStatus(eventId, StatusRemoved)
}

// UpdateTitle sets the title of the event
func (c *Calendar) UpdateTitle(eventId int64, title string, editType RepeatEditType) error {
	return c.dataStore.SetTitle(eventId, title)
}

// UpdateDescription sets the description of the event
func (c *Calendar) UpdateDescription(eventId int64, description *string, editType RepeatEditType) error {
	return c.dataStore.SetDescription(eventId, description)
}

// UpdateUrl sets the url link of the event
func (c *Calendar) UpdateUrl(eventId int64, url *string, editType RepeatEditType) error {
	return c.dataStore.SetUrl(eventId, url)
}

// UpdateUserData sets the user data for the event
func (c *Calendar) UpdateUserData(eventId int64, userData map[string]interface{}, editType RepeatEditType) error {
	return c.dataStore.SetUserData(eventId, userData)
}

// ///////////////////////
// Invites
// ///////////////////////

// GetInvitation grabs a single matching invite from the data store or nil if it does not exist
func (c *Calendar) GetInvitation(eventId int64, userId int64) (*Invite, error) {
	return c.dataStore.GetInvite(eventId, userId)
}

// AcceptInvitation changes the status of an invitation to InviteStatusConfirmed
func (c *Calendar) AcceptInvitation(eventId int64, userId int64, editType RepeatEditType) error {
	return c.dataStore.SetInviteStatus(eventId, userId, InviteStatusConfirmed)
}

// DeclineInvitation changes the status of an invitation to InviteStatusDeclined
func (c *Calendar) DeclineInvitation(eventId int64, userId int64, editType RepeatEditType) error {
	return c.dataStore.SetInviteStatus(eventId, userId, InviteStatusDeclined)
}

// RevokeInvitation changes the status of an invitation to InviteStatusRevoked (we never delete things)
func (c *Calendar) RevokeInvitation(eventId int64, userId int64, editType RepeatEditType) error {
	return c.dataStore.SetInviteStatus(eventId, userId, InviteStatusRevoked)
}

// InviteUser creates a pending invitation for a user on an event
func (c *Calendar) InviteUser(eventId int64, userId int64, permission Permission, editType RepeatEditType) error {
	i := Invite{
		EventId:    eventId,
		UserId:     userId,
		Status:     InviteStatusPending,
		Permission: permission,
		Created:    time.Now(),
	}
	i.Updated = i.Created
	if err := ValidateInvite(i); err != nil {
		return err
	}
	_, err := c.dataStore.AddInvite(i)
	return err
}

// UpdateInvitationPermission sets the permission of a user on an event
func (c *Calendar) UpdateInvitationPermission(eventId int64, userId int64, permission Permission, editType RepeatEditType) error {
	return c.dataStore.SetInvitePermissions(eventId, userId, permission)
}

// ///////////////////////
// Helpers
// ///////////////////////
var heatDeathOfTheUniverse = time.Unix(999999999, 0)

// getAllRepeatingEvents collects all the events that match the parent id of this event (including this event).
// Or if the parent id is nil, then it just returns this event.
func (c *Calendar) getAllRepeatingEvents(e Event) ([]*Event, error) {
	var result []*Event
	if e.ParentId == nil {
		result = append(result, &e)
		return result, nil
	}
	return c.dataStore.Query(Query{
		Start:     time.Unix(0, 0),
		End:       heatDeathOfTheUniverse,
		ParentIds: []int64{*e.ParentId},
	})
}

// getAllRepeatingEventsThisAndAfter collects all the events that match the parent id of this event (including this event) and are at or after the start day and time of this event.
// Or if the parent id is nil, then it just returns this event.
func (c *Calendar) getAllRepeatingEventsThisAndAfter(e Event) ([]*Event, error) {
	var result []*Event
	if e.ParentId == nil {
		result = append(result, &e)
		return result, nil
	}
	start, err := time.Parse(fmt.Sprintf("%sT%s", time.DateOnly, TimeFormat), fmt.Sprintf("%sT%s", e.StartDay, e.StartTime))
	if err != nil {
		return nil, ErrorInvalidStartDay
	}
	return c.dataStore.Query(Query{
		Start:     start,
		End:       heatDeathOfTheUniverse,
		ParentIds: []int64{*e.ParentId},
	})
}
