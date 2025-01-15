package cali

import (
	"testing"

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
