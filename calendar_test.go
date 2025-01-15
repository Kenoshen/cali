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

	events, err := c.Query(Query{
		Start: time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	assert.Len(t, events, 6)
	foreach(events, func(e Event) {
		t.Log("Event", e.Id, "P:", *e.ParentId)
	})

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

			events, err := c.Query(Query{
				Start: time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
			})
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
					times = append(times, fmt.Sprintf("%s - %s %s", e.StartDay, e.EndDay, e.Zone))
				} else {
					times = append(times, fmt.Sprintf("%sT%s - %sT%s %s", e.StartDay, e.StartTime, e.EndDay, e.EndTime, e.Zone))
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
