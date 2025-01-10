package cali

type Calendar struct {
	dataStore DataStore
}

func NewCalendar(dataStore DataStore) *Calendar {
	c := &Calendar{
		dataStore: dataStore,
	}
	return c
}

func (c *Calendar) Get(eventId int64) (*Event, error) {
	return c.dataStore.Get(eventId)
}

func (c *Calendar) Query(q Query) ([]*Event, error) {
	return c.dataStore.Query(q)
}

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

func (c *Calendar) UpdateTime(eventId int64) error {
	// TODO: implement me
	return nil
}

func (c *Calendar) Cancel(eventId int64) error {
	// TODO: implement me
	return nil
}

func (c *Calendar) Remove(eventId int64) error {
	// TODO: implement me
	return nil
}

func (c *Calendar) UpdateDetails(details Event) error {
	// TODO: implement me
	return nil
}

func (c *Calendar) SetUserData(eventId int64, userData map[string]interface{}) error {
	// TODO: implement me
	return nil
}

// ///////////////////////
// Invites
// ///////////////////////

func (c *Calendar) AcceptInvitation(eventId int64, userId int64) error {
	// TODO: implement me
	return nil
}

func (c *Calendar) DeclineInvitation(eventId int64, userId int64) error {
	// TODO: implement me
	return nil
}

func (c *Calendar) InviteUser(eventId int64, userId int64, permission Permission) error {
	// TODO: implement me
	return nil
}

func (c *Calendar) UpdateInvitationPermission(eventId int64, userId int64, permission Permission) error {
	// TODO: implement me
	return nil
}

func (c *Calendar) RevokeInvitation(eventId int64, userId int64) error {
	// TODO: implement me
	return nil
}
