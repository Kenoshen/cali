package cali

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalendar(t *testing.T) {
	d := &InMemoryDataStore{}
	c := NewCalendar(d)

	a, count, err := c.Create(Event{
		StartDay:  "2008-01-01",
		StartTime: "09:00",
		EndDay:    "2008-01-01",
		EndTime:   "10:00",
		Zone:      "America/Denver",
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
	require.NotNil(t, a)

	err = c.UpdateTime(a.Id, "2008-02-01", "10:00", "2008-02-01", "11:00", "America/Denver", false, RepeatEditTypeThis)
	require.NoError(t, err)

	originalStatus := a.Status
	assert.NotEqual(t, StatusCanceled, a.Status)
	err = c.Cancel(a.Id, RepeatEditTypeThis)
	require.NoError(t, err)
	assert.NotEqual(t, originalStatus, a.Status)
	assert.Equal(t, StatusCanceled, a.Status)

	originalStatus = a.Status
	assert.NotEqual(t, StatusRemoved, a.Status)
	err = c.Remove(a.Id, RepeatEditTypeThis)
	require.NoError(t, err)
	assert.NotEqual(t, originalStatus, a.Status)
	assert.Equal(t, StatusRemoved, a.Status)

	originalTitle := a.Title
	assert.NotEqual(t, "New Title", a.Title)
	err = c.UpdateTitle(a.Id, "New Title", RepeatEditTypeThis)
	require.NoError(t, err)
	assert.NotEqual(t, originalTitle, a.Title)
	assert.Equal(t, "New Title", a.Title)

	originalUserData := a.UserData
	assert.NotEqual(t, map[string]interface{}{"key": "value"}, a.UserData)
	err = c.UpdateUserData(a.Id, map[string]interface{}{"key": "value"}, RepeatEditTypeThis)
	require.NoError(t, err)
	assert.NotEqual(t, originalUserData, a.UserData)
	assert.Equal(t, map[string]interface{}{"key": "value"}, a.UserData)

	err = c.InviteUser(a.Id, 7, PermissionInvitee, RepeatEditTypeThis)
	require.NoError(t, err)
	invite, err := c.GetInvitation(a.Id, 7)
	require.NoError(t, err)
	require.NotNil(t, invite)

	originalInvitationStatus := invite.Status
	assert.NotEqual(t, InviteStatusConfirmed, invite.Status)
	err = c.AcceptInvitation(a.Id, 7, RepeatEditTypeThis)
	require.NoError(t, err)
	assert.NotEqual(t, originalInvitationStatus, invite.Status)
	assert.Equal(t, InviteStatusConfirmed, invite.Status)

	originalInvitationStatus = invite.Status
	assert.NotEqual(t, InviteStatusDeclined, invite.Status)
	err = c.DeclineInvitation(a.Id, 7, RepeatEditTypeThis)
	require.NoError(t, err)
	assert.NotEqual(t, originalInvitationStatus, invite.Status)
	assert.Equal(t, InviteStatusDeclined, invite.Status)
}

func TestRepeatEventsOnCalendar(t *testing.T) {
	d := &InMemoryDataStore{}
	c := NewCalendar(d)

	a, count, err := c.Create(Event{
		Id:          -10,
		StartDay:    "2008-01-01",
		EndDay:      "2008-01-01",
		Zone:        "America/Denver",
		IsAllDay:    true,
		IsRepeating: true,
		Repeat: &Repeat{
			RepeatType:     RepeatTypeWeekly,
			DayOfWeek:      DayOfWeekTuesday | DayOfWeekThursday,
			RepeatStopDate: _t(time.Date(2008, time.January, 17, 0, 0, 0, 0, time.UTC)),
		},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(6), count)
	assert.Len(t, d.events, 6)
	require.NotNil(t, a)

	events, err := c.Query(Query{})
	require.NoError(t, err)
	assert.Len(t, events, 6)

	foreach(events, func(e Event) {
		assert.Equalf(t, StatusActive, e.Status, "failed on event with id: %v", e.Id)
	})
	err = c.Cancel(a.Id, RepeatEditTypeAll)
	require.NoError(t, err)
	foreach(events, func(e Event) {
		assert.Equalf(t, StatusCanceled, e.Status, "failed on event with id: %v", e.Id)
	})

	foreach(events, func(e Event) {
		assert.Equalf(t, "", e.Title, "failed on event with id: %v", e.Id)
	})
	err = c.UpdateTitle(events[3].Id, "New Title", RepeatEditTypeThisAndAfter)
	require.NoError(t, err)
	foreach(events[:3], func(e Event) {
		assert.Equalf(t, "", e.Title, "failed on event with id: %v", e.Id)
	})
	foreach(events[3:], func(e Event) {
		assert.Equalf(t, "New Title", e.Title, "failed on event with id: %v", e.Id)
	})

	foreach(events, func(e Event) {
		assert.Nilf(t, e.Description, "failed on event with id: %v", e.Id)
	})
	desc := "Some description"
	err = c.UpdateDescription(events[1].Id, &desc, RepeatEditTypeThis)
	require.NoError(t, err)
	foreach(events[:1], func(e Event) {
		assert.Nilf(t, e.Description, "failed on event with id: %v", e.Id)
	})
	foreach(events[2:], func(e Event) {
		assert.Nilf(t, e.Description, "failed on event with id: %v", e.Id)
	})
	foreach(events[1:1], func(e Event) {
		assert.NotNilf(t, e.Description, "failed on event with id: %v", e.Id)
		if e.Description != nil {
			assert.Equalf(t, "Some description", *e.Description, "failed on event with id: %v", e.Id)
		}
	})
}

const den = "America/Denver"

func TestUpdateTimeOnRepeatEvent(t *testing.T) {
	// Events:
	// #1 Jan 01 08:00-09:00
	// #2 Jan 03 08:00-09:00
	// #3 Jan 08 08:00-09:00
	// #4 Jan 10 08:00-09:00
	// #5 Jan 15 08:00-09:00
	// #6 Jan 17 08:00-09:00
	testCases := []struct {
		name      string
		eventId   int64
		startDay  string
		startTime string
		endDay    string
		endTime   string
		zone      string
		isAllDay  bool
		editType  RepeatEditType
		times     []string
		err       string
	}{
		{
			name:      "event not found",
			eventId:   -1,
			startDay:  "2008-01-01",
			startTime: "08:00",
			endDay:    "2008-01-01",
			endTime:   "09:00",
			zone:      den,
			isAllDay:  false,
			editType:  RepeatEditTypeThis,
			times:     nil,
			err:       ErrorEventNotFound.Error(),
		},
		{
			name:      "no change",
			eventId:   4,
			startDay:  "2008-01-10",
			startTime: "08:00",
			endDay:    "2008-01-10",
			endTime:   "09:00",
			zone:      den,
			isAllDay:  false,
			editType:  RepeatEditTypeThis,
			times: []string{
				"2008-01-01 08:00 - 2008-01-01 09:00",
				"2008-01-03 08:00 - 2008-01-03 09:00",
				"2008-01-08 08:00 - 2008-01-08 09:00",
				"2008-01-10 08:00 - 2008-01-10 09:00",
				"2008-01-15 08:00 - 2008-01-15 09:00",
				"2008-01-17 08:00 - 2008-01-17 09:00",
			},
		},
		{
			name:      "single event time change",
			eventId:   4,
			startDay:  "2008-01-10",
			startTime: "13:00",
			endDay:    "2008-01-10",
			endTime:   "13:45",
			zone:      den,
			isAllDay:  false,
			editType:  RepeatEditTypeThis,
			times: []string{
				"2008-01-01 08:00 - 2008-01-01 09:00",
				"2008-01-03 08:00 - 2008-01-03 09:00",
				"2008-01-08 08:00 - 2008-01-08 09:00",
				"2008-01-10 13:00 - 2008-01-10 13:45",
				"2008-01-15 08:00 - 2008-01-15 09:00",
				"2008-01-17 08:00 - 2008-01-17 09:00",
			},
		},
		{
			name:      "singe event time and day change",
			eventId:   4,
			startDay:  "2008-01-11",
			startTime: "13:00",
			endDay:    "2008-01-11",
			endTime:   "13:45",
			zone:      den,
			isAllDay:  false,
			editType:  RepeatEditTypeThis,
			times: []string{
				"2008-01-01 08:00 - 2008-01-01 09:00",
				"2008-01-03 08:00 - 2008-01-03 09:00",
				"2008-01-08 08:00 - 2008-01-08 09:00",
				"2008-01-11 13:00 - 2008-01-11 13:45",
				"2008-01-15 08:00 - 2008-01-15 09:00",
				"2008-01-17 08:00 - 2008-01-17 09:00",
			},
		},
		{
			name:      "all event time changes",
			eventId:   4,
			startDay:  "2008-01-10",
			startTime: "13:00",
			endDay:    "2008-01-10",
			endTime:   "13:45",
			zone:      den,
			isAllDay:  false,
			editType:  RepeatEditTypeAll,
			times: []string{
				"2008-01-01 13:00 - 2008-01-01 13:45",
				"2008-01-03 13:00 - 2008-01-03 13:45",
				"2008-01-08 13:00 - 2008-01-08 13:45",
				"2008-01-10 13:00 - 2008-01-10 13:45",
				"2008-01-15 13:00 - 2008-01-15 13:45",
				"2008-01-17 13:00 - 2008-01-17 13:45",
			},
		},
		{
			name:      "all events after or on event time change",
			eventId:   4,
			startDay:  "2008-01-10",
			startTime: "13:00",
			endDay:    "2008-01-10",
			endTime:   "13:45",
			zone:      den,
			isAllDay:  false,
			editType:  RepeatEditTypeThisAndAfter,
			times: []string{
				"2008-01-01 08:00 - 2008-01-01 09:00",
				"2008-01-03 08:00 - 2008-01-03 09:00",
				"2008-01-08 08:00 - 2008-01-08 09:00",
				"2008-01-10 13:00 - 2008-01-10 13:45",
				"2008-01-15 13:00 - 2008-01-15 13:45",
				"2008-01-17 13:00 - 2008-01-17 13:45",
			},
		},
		{
			name:      "invalid date change to all",
			eventId:   4,
			startDay:  "2008-01-11",
			startTime: "08:00",
			endDay:    "2008-01-11",
			endTime:   "09:00",
			zone:      den,
			isAllDay:  false,
			editType:  RepeatEditTypeAll,
			err:       "something",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			t.Parallel()

			d := &InMemoryDataStore{}
			c := NewCalendar(d)

			a, count, err := c.Create(Event{
				StartDay:    "2008-01-01",
				StartTime:   "08:00",
				EndDay:      "2008-01-01",
				EndTime:     "09:00",
				Zone:        "America/Denver",
				IsAllDay:    false,
				IsRepeating: true,
				Repeat: &Repeat{
					RepeatType:     RepeatTypeWeekly,
					DayOfWeek:      DayOfWeekTuesday | DayOfWeekThursday,
					RepeatStopDate: _t(time.Date(2008, time.January, 17, 0, 0, 0, 0, time.UTC)),
				},
			})
			require.NoError(t, err)
			assert.Equal(t, int64(6), count)
			assert.Len(t, d.events, 6)
			require.NotNil(t, a)

			// get all events in the database
			events, err := c.Query(Query{})
			require.NoError(t, err)
			assert.Len(t, events, 6)

			err = c.UpdateTime(tc.eventId, tc.startDay, tc.startTime, tc.endDay, tc.endTime, tc.zone, tc.isAllDay, tc.editType)
			if tc.err == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err)
				// stop processing if there is an error here
				return
			}

			var times []string
			for _, e := range events {
				if e.IsAllDay {
					times = append(times, fmt.Sprintf("%s - %s", e.StartDay, e.EndDay))
				} else {
					times = append(times, fmt.Sprintf("%s %s - %s %s", e.StartDay, e.StartTime, e.EndDay, e.EndTime))
				}
			}
			assert.Equal(t, tc.times, times)
		})
	}

}

func foreach(events []*Event, f func(e Event)) {
	for _, e := range events {
		if e != nil {
			f(*e)
		}
	}
}
