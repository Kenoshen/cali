package cali

type Calendar struct {
	dataSource DataSource
}

func NewCalendar(dataSource DataSource) *Calendar {
	c := &Calendar{
		dataSource: dataSource,
	}
	return c
}

func (c *Calendar) Get(eventId int64) (*Event, error) {
	return c.dataSource.Get(eventId)
}

func (c *Calendar) Query(q Query) ([]*Event, error) {
	return c.dataSource.Query(q)
}

func (c *Calendar) Create(e Event) (*Event, int64, error) {
	if err := Validate(e); err != nil {
		return nil, 0, err
	}

	if !e.IsRepeating {
		newEvent, err := c.dataSource.Create(e)
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

	return events[0], int64(len(events)), nil
}

func (c *Calendar) UpdateTime(eventId int64) error {
	return nil
}

func (c *Calendar) Cancel(eventId int64) error {
	return nil
}

func (c *Calendar) Remove(eventId int64) error {
	return nil
}

func (c *Calendar) UpdateDetails(details Event) error {
	return nil
}

func (c *Calendar) SetUserData(eventId int64, userData map[string]interface{}) error {
	return nil
}

// ///////////////////////
// Invites
// ///////////////////////

func (c *Calendar) AcceptInvitation(eventId int64, userId int64) error {
	return nil
}

func (c *Calendar) DeclineInvitation(eventId int64, userId int64) error {
	return nil
}

func (c *Calendar) InviteUser(eventId int64, userId int64, permission Permission) error {
	return nil
}

func (c *Calendar) UpdateInvitationPermission(eventId int64, userId int64, permission Permission) error {
	return nil
}

func (c *Calendar) RevokeInvitation(eventId int64, userId int64) error {
	return nil
}
