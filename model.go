package cali

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Event is a single record of an event and also contains links to other events through
// the parent id as well as a ref to the repeat logic used to create the event if it is
// a repeating event.
type Event struct {
	// Id is the unique id for this event
	Id int64 `json:"id"`
	// SourceId represents an id for an external source object that this event is directly tied to
	SourceId *int64 `json:"sourceId"`
	// ParentId is the id of another event that this event is related to via repeating events
	// and can be used to update other related repeating events when this one changes
	ParentId *int64 `json:"parentId"`
	// OwnerId is the id of the user that created this event
	OwnerId int64 `json:"ownerId"`
	// EventType represents the overall type of the event. This is just an int, so you can set this
	// to what ever you would like
	EventType EventType `json:"eventType"`

	// Title is the value that will be shown for this event when displayed on a calendar interface
	Title string `json:"title"`
	// Description is a longer field description of what the event is
	Description *string `json:"description"`
	// Url is a quick way to set the destination on an event that is clicked on in an interface
	Url *string `json:"url"`
	// Status represents the current status of the event, defaults to active, but events can also
	// be canceled or removed
	Status Status `json:"status"`

	// IsAllDay is true if the event is an all day event which will set the time values to 00:00
	IsAllDay bool `json:"isAllDay"`

	// IsRepeating is true if this event is a part of a repeating series
	IsRepeating bool `json:"isRepeating"`
	// Repeat is the pattern to repeat the event
	Repeat *Repeat `json:"repeat"`

	// Zone must be a valid time.Location name like "UTC" or "America/New_York"
	Zone string `json:"zone"`

	// StartDay is the YYYY-MM-DD value representing the start day of this event
	StartDay string `json:"startDay"`
	// StartTime is the HH:MM value representing the start time of this event
	StartTime string `json:"startTime"`

	// EndDay is the YYYY-MM-DD value representing the end day of this event
	EndDay string `json:"endDay"`
	// EndTime is the HH:MM value representing the end time of this event
	EndTime string `json:"endTime"`

	// Created is a UTC timestamp for when the event was created
	Created time.Time `json:"created"`
	// Updated is a UTC timestamp for when the event was modified last
	Updated time.Time `json:"updated"`

	// UserData is a custom and optional blob of JSON saved to the event
	UserData map[string]interface{} `json:"userData"`
}

// Start gets the time.Time value using the StartDay and StartTime fields
func (e Event) Start() (time.Time, error) {
	return parseDayTime(e.StartDay, e.StartTime)
}

// End gets the time.Time value using the EndDay and EndTime fields
func (e Event) End() (time.Time, error) {
	return parseDayTime(e.EndDay, e.EndTime)
}

const iCalDateTimeFormat = "20060102T150400Z"

// MarshallToICal marshalls this event to an ical format
func (e Event) MarshallToICal() string {
	start, _ := e.Start()
	end, _ := e.Start()
	s := []string{
		"BEGIN:VEVENT",
		fmt.Sprintf("UID:%v", e.Id),
		fmt.Sprintf("DTSTAMP:%v", start.Format(iCalDateTimeFormat)),
		fmt.Sprintf("DTSTART:%v", start.Format(iCalDateTimeFormat)),
		fmt.Sprintf("DTEND:%v", end.Format(iCalDateTimeFormat)),
		fmt.Sprintf("SUMMARY:%v", strings.ReplaceAll(e.Title, "\n", " ")),
		"CLASS:PRIVATE",
	}
	if e.Description != nil && len(*e.Description) > 0 {
		s = append(s, fmt.Sprintf("DESCRIPTION:", *e.Description))
	}

	s = append(s, "END:VEVENT")
	return strings.Join(s, "\n")
}

// parseDayTime takes a day of YYYY-MM-DD and an hourMin as HH-mm (or "")
// and converts it into a time.Time object
func parseDayTime(day, hourMin string) (time.Time, error) {
	if day == "" {
		return time.Time{}, fmt.Errorf("invalid day value")
	}
	if hourMin == "" {
		return time.Parse(time.DateOnly, day)
	}

	return time.Parse(DayTimeFormat, fmt.Sprintf("%s %s", day, hourMin))
}

// DayTimeFormat is the time package format style for YYYY-MM-DD HH:mm
const DayTimeFormat = time.DateOnly + " 15:04"

// TimeFormat is the time package format style for HH:mm
const TimeFormat = "15:04"

type Details struct {
	// Id is the unique id for this event
	Id int64

	// Title is the value that will be shown for this event when displayed on a calendar interface
	Title string
	// Description is a longer field description of what the event is
	Description *string
	// Url is a quick way to set the destination on an event that is clicked on in an interface
	Url *string
	// Status represents the current status of the event, defaults to active, but events can also
	// be canceled or removed
	Status Status

	// IsAllDay is true if the event is an all day event which will set the time values to 00:00
	IsAllDay bool

	// Zone must be a valid time.Location name like "UTC" or "America/New_York"
	Zone string

	// StartDay is the YYYY-MM-DD value representing the start day of this event
	StartDay string
	// StartTime is the HH:MM value representing the start time of this event
	StartTime string

	// EndDay is the YYYY-MM-DD value representing the end day of this event
	EndDay string
	// EndTime is the HH:MM value representing the end time of this event
	EndTime string
}

type Status int64

const (
	// StatusActive is the default for events and it means it will show up on the calendar as a standard event
	StatusActive Status = 0
	// StatusCanceled is for when the owner of the event cancels the event but it still shows up on the calendar
	// as a faded out event with a line through the text
	StatusCanceled Status = 1
	// StatusAbandoned is when all of the invitees of the event (including the owner of the event) declines the event
	// which then causes the event to disappear from the calendar
	StatusAbandoned Status = -2
	// StatusRemoved is when the event was deleted by the owner of the event and it disappears from the calendar
	StatusRemoved Status = -1
)

// EventType must be defined by the user of this library
type EventType = int64

// Invite is a record of an invitation to a specific user for a specific event
type Invite struct {
	// EventId is a reference to the unique identifier for a specific event
	EventId int64
	// UserId is the reference for the user who's invite is in question
	UserId int64
	// Status refers to the response of the user to the invite of an event
	// and defaults to pending which is kind of like a soft confirm
	Status InviteStatus
	// Permission is a bitmask for the allowed permissions for this user on this event
	Permission Permission
	// Created is a timestamp for when the invite invitation was created
	Created time.Time
	// Updated is a timestamp for when the invite invitation was modified last
	Updated time.Time
}

func (i Invite) String() string {
	return fmt.Sprintf("{Event:%v, User:%v, Status:%v, Perm:%v}", i.EventId, i.UserId, i.Status, i.Permission)
}

// InviteStatus conveys the invitation status of this invitation. Statuses that are equal or
// greater to zero will be considered positive statuses for the purpose of showing the event
// on that user's calendar. Anything less than 0 will be hidden on the user's calendar.
type InviteStatus int64

const (
	// InviteStatusPending is the default and refers to a non-answer, pending invite will still be treated
	// as a soft confirm and the event will remain on the user's calendar but be outlined
	InviteStatusPending InviteStatus = 0
	// InviteStatusConfirmed is an acknowledgment that the user is going to attend the event
	InviteStatusConfirmed InviteStatus = 1
	// InviteStatusDeclined is when the user decides tho not attend the event, if all users decline an event
	// it becomes abandoned
	InviteStatusDeclined InviteStatus = -1
	// InviteStatusRevoked is when a user with the correct permission forcibly removes a user's invitation
	InviteStatusRevoked InviteStatus = -2
)

type Bitmask uint32

func (f Bitmask) HasFlag(flag Bitmask) bool {
	return f&flag != 0
}

func (f *Bitmask) AddFlag(flag Bitmask) {
	*f |= flag
}

func (f *Bitmask) ClearFlag(flag Bitmask) {
	*f &= ^flag
}

func (f *Bitmask) ToggleFlag(flag Bitmask) {
	*f ^= flag
}

type Permission = Bitmask

const (
	PermissionRead = 1 << iota
	PermissionModify
	PermissionInvite
	PermissionCancel
	PermissionDelete
)

const (
	PermissionOwner   = PermissionDelete | PermissionCancel | PermissionModify | PermissionInvite | PermissionRead
	PermissionInvitee = PermissionRead
)

// MaxRepeatOccurrence is set to 30 events
const MaxRepeatOccurrence int64 = 30

// MaxRepeatDuration is set to 2 years
const MaxRepeatDuration = time.Duration(24*365*2) * time.Hour

// Repeat contains all of the values required to be able to repeat an event
// over a period of time or for a number of occurrences
type Repeat struct {
	// RepeatType is a enumeration of the valid types of repeat events (daily,
	// weekly, monthly, or yearly)
	RepeatType RepeatType `json:"repeatType"`
	// DayOfWeek is a bitmask of the days of the week (SMTWTFS)
	DayOfWeek DayOfWeek `json:"dayOfWeek"`
	// RepeatOccurrences is a number of times the event should repeat.
	// It should be 0 if RepeatStopDate is not nil.
	// It can't be more than MaxRepeatOccurrence.
	RepeatOccurrences int64 `json:"repeatOccrrences"`
	// RepeatStopDate is a timestamp for when the repeating event should stop.
	// It should be nil if RepeatOccurrences > 1.
	// It can't be more than MaxRepeatDuration.
	RepeatStopDate *time.Time `json:"repeatStopDate"`
}

type RepeatType int64

const (
	RepeatTypeDaily   RepeatType = 0
	RepeatTypeWeekly  RepeatType = 1
	RepeatTypeMonthly RepeatType = 2
	RepeatTypeYearly  RepeatType = 3
)

type DayOfWeek = Bitmask

const (
	DayOfWeekSunday = 1 << iota
	DayOfWeekMonday
	DayOfWeekTuesday
	DayOfWeekWednesday
	DayOfWeekThursday
	DayOfWeekFriday
	DayOfWeekSaturday
)

func dayOfWeekFromWeekday(w time.Weekday) DayOfWeek {
	switch w {
	case time.Sunday:
		return DayOfWeekSunday
	case time.Monday:
		return DayOfWeekMonday
	case time.Tuesday:
		return DayOfWeekTuesday
	case time.Wednesday:
		return DayOfWeekWednesday
	case time.Thursday:
		return DayOfWeekThursday
	case time.Friday:
		return DayOfWeekFriday
	case time.Saturday:
		return DayOfWeekSaturday
	}
	return DayOfWeekSunday
}

func _t(t time.Time) *time.Time {
	return &t
}

// Query is the object that the data store uses to try and find the list of matching events
type Query struct {
	// Start is an inclusive timestamp and should be compared against the end timestamp of other events (overlap)
	Start *time.Time
	// End is an inclusive timestamp and should be compared against the start timestamp of other events (overlap)
	End *time.Time
	// EventIds is a list of specific events that you want to query
	EventIds []int64
	// ParentIds is a list of parent ids that should be searched for and will find all events that have a match to the parent id
	ParentIds []int64
	// UserIds is a check if the user has an invite record for the event that is not
	// declined or revoked
	UserIds []int64
	// EventTypes is a check if the event has a specific event type
	EventTypes []EventType
	// SourceIds is an OR check on the source ids
	SourceIds []int64
	// Statuses is an OR search for specific statuses
	Statuses []Status
	// Text is an OR search for specific words
	Text []string
}

// Matches does a local check if the given event matches the query
func (q Query) Matches(event *Event) bool {
	if event == nil {
		return false
	}

	if q.Start != nil {
		startDay := q.Start.Format(time.DateOnly)
		startTime := q.Start.Format(TimeFormat)
		if startDay > event.EndDay {
			return false
		}
		if event.EndTime != "" && startDay+startTime > event.EndDay+event.EndTime {
			return false
		}
	}

	if q.End != nil {
		endDay := q.End.Format(time.DateOnly)
		endTime := q.End.Format(TimeFormat)
		if endDay < event.StartDay {
			return false
		}
		if event.StartTime != "" && endDay+endTime < event.StartDay+event.StartTime {
			return false
		}
	}

	found := false
	if len(q.EventIds) > 0 {
		for _, id := range q.EventIds {
			if event.Id == id {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(q.ParentIds) > 0 {
		found = false
		for _, id := range q.ParentIds {
			if event.ParentId != nil && *event.ParentId == id {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(q.EventTypes) > 0 {
		found = false
		for _, eventType := range q.EventTypes {
			if event.EventType == eventType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(q.SourceIds) > 0 {
		found = false
		for _, id := range q.SourceIds {
			if event.SourceId != nil && *event.SourceId == id {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(q.Statuses) > 0 {
		found = false
		for _, status := range q.Statuses {
			if event.Status == status {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(q.Text) > 0 {
		found = false
		for _, text := range q.Text {
			if strings.Contains(event.Title, text) {
				found = true
				break
			}
			if event.Description != nil && strings.Contains(*event.Description, text) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

type RepeatEditType int64

const (
	RepeatEditTypeThis         RepeatEditType = 0
	RepeatEditTypeAll          RepeatEditType = 1
	RepeatEditTypeThisAndAfter RepeatEditType = 2
)

// Sort events by their start day and time where earlier events
// are first and later events are last
func Sort(e []*Event) []*Event {
	sort.SliceStable(e, func(a int, b int) bool {
		A := e[a]
		B := e[b]
		if A == nil {
			return true
		}
		if B == nil {
			return false
		}
		if A.StartDay < B.StartDay {
			return true
		} else if A.StartDay > B.StartDay {
			return false
		}
		if A.StartTime <= B.StartTime {
			return true
		}
		return false
	})
	return e
}
