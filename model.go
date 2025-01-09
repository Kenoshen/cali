package cali

import (
	"strings"
	"time"
)

type Event struct {
	// Id is the unique id for this event
	Id int64
	// SourceId represents an id for an external source object that this event is directly tied to
	SourceId *int64
	// ParentId is the id of another event that this event is related to via repeating events
	// and can be used to update other related repeating events when this one changes
	ParentId *int64
	// OwnerId is the id of the user that created this event
	OwnerId int64
	// EventType represents the overall type of the event. This is just an int, so you can set this
	// to what ever you would like
	EventType EventType

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

	// IsRepeating is true if this event is a part of a repeating series
	IsRepeating bool
	// Repeat is the pattern to repeat the event
	Repeat *Repeat

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

	// Created is a UTC timestamp for when the event was created
	Created time.Time
	// Updated is a UTC timestamp for when the event was modified last
	Updated time.Time

	// UserData is a custom and optional blob of JSON saved to the event
	UserData map[string]interface{}
}

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

// Attendance is a record of an invitation to a specific user for a specific event
type Attendance struct {
	// EventId is a reference to the unique identifier for a specific event
	EventId int64
	// UserId is the reference for the user who's attendance is in question
	UserId int64
	// Status refers to the response of the user to the attendance of an event
	// and defaults to pending which is kind of like a soft confirm
	Status AttendanceStatus
	// Permission is a bitmask for the allowed permissions for this user on this event
	Permission Permission
	// Created is a timestamp for when the attendance invitation was created
	Created time.Time
	// Updated is a timestamp for when the attendance invitation was modified last
	Updated time.Time
}

type AttendanceStatus int64

const (
	// AttendanceStatusPending is the default and refers to a non-answer, pending attendance will still be treated
	// as a soft confirm and the event will remain on the user's calendar but be outlined
	AttendanceStatusPending AttendanceStatus = 0
	// AttendanceStatusConfirmed is an acknowledgment that the user is going to attend the event
	AttendanceStatusConfirmed AttendanceStatus = 1
	// AttendanceStatusDeclined is when the user decides tho not attend the event, if all users decline an event
	// it becomes abandoned
	AttendanceStatusDeclined AttendanceStatus = 2
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

type Repeat struct {
	// RepeatType is a enumeration of the valid types of repeat events (daily,
	// weekly, monthly, or yearly)
	RepeatType RepeatType
	// DayOfWeek is a bitmask of the days of the week (SMTWTFS)
	DayOfWeek DayOfWeek
	// RepeatOccurrences is a number of times the event should repeat. It can't
	// be more than MaxRepeatOccurrence.
	RepeatOccurrences int64
	// RepeatStopDate is a timestamp for when the repeating event should stop.
	// It can't be more than MaxRepeatDuration.
	RepeatStopDate *time.Time
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
	// Statuses is an OR search for specific statuses
	Statuses []Status
	// Text is an OR search for specific words
	Text []string
}

func (q Query) Matches(event *Event) bool {
	startDay := q.Start.Format(time.DateOnly)
	startTime := q.Start.Format(TimeFormat)
	endDay := q.End.Format(time.DateOnly)
	endTime := q.End.Format(TimeFormat)
	if event.StartDay > endDay {
		return false
	}
	if event.EndDay < startDay {
		return false
	}
	if event.StartDay+event.StartTime > endDay+endTime {
		return false
	}
	if event.EndDay+event.EndTime < startDay+startTime {
		return false
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
		found = false
	}

	if len(q.UserIds) > 0 {
		for _, id := range q.UserIds {
			if event.OwnerId == id {
				found = true
				break
			}
		}
		if !found {
			return false
		}
		found = false
	}
	if len(q.EventTypes) > 0 {
		for _, eventType := range q.EventTypes {
			if event.EventType == eventType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
		found = false
	}
	if len(q.SourceIds) > 0 {
		for _, id := range q.SourceIds {
			if event.SourceId != nil && *event.SourceId == id {
				found = true
				break
			}
		}
		if !found {
			return false
		}
		found = false
	}
	if len(q.Statuses) > 0 {
		for _, status := range q.Statuses {
			if event.Status == status {
				found = true
				break
			}
		}
		if !found {
			return false
		}
		found = false
	}
	if len(q.Text) > 0 {
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
		found = false
	}
	return true
}
