package cali

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryDataStore(t *testing.T) {
	// check that the TestInMemoryDataStore is an actual implementation of
	// the DataStore interface
	var dataStoreInterface DataStore = &InMemoryDataStore{}
	var d *InMemoryDataStore = dataStoreInterface.(*InMemoryDataStore)

	a, err := d.Create(Event{Status: StatusActive, StartDay: "2008-01-01", EndDay: "2008-01-01", IsAllDay: true})
	require.NoError(t, err)
	assert.Len(t, d.events, 1)
	assert.Len(t, d.invites, 1)
	assert.Equal(t, d.events[0], a)

	a1, err := d.Get(a.Id)
	require.NoError(t, err)
	assert.Len(t, d.events, 1)
	assert.Len(t, d.invites, 1)
	assert.Equal(t, a, a1)

	// save a copy of the original before it gets updated
	original := *a
	a2, err := d.Update(Details{Id: a.Id, Status: StatusCanceled, StartDay: "2008-01-01", EndDay: "2008-01-01"})
	require.NoError(t, err)
	assert.Len(t, d.events, 1)
	assert.Len(t, d.invites, 1)
	assert.NotEqual(t, a2, &original)
	assert.Equal(t, a2.Status, StatusCanceled)

	d.Create(Event{Status: StatusActive, StartDay: "2008-01-01", EndDay: "2008-01-01", IsAllDay: true})
	d.Create(Event{Status: StatusRemoved, StartDay: "2008-01-01", EndDay: "2008-01-01", IsAllDay: true})
	d.Create(Event{Status: StatusActive, StartDay: "2008-01-01", EndDay: "2008-01-01", IsAllDay: true})
	d.Create(Event{Status: StatusRemoved, StartDay: "2008-01-01", EndDay: "2008-01-01", IsAllDay: true})
	assert.Len(t, d.events, 5)
	assert.Len(t, d.invites, 5)

	res, err := d.Query(Query{Start: time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC), Statuses: []Status{StatusActive}})
	assert.Len(t, res, 2)
}
