package cali

import "time"

func GenerateRepeatEvents(e Event) ([]*Event, error) {
	if !e.IsRepeating {
		return nil, ErrorNotRepeatingEvent
	}

	startDay, err := time.Parse(time.DateOnly, e.StartDay)
	if err != nil {
		return nil, ErrorInvalidStartDay
	}
	endDay, err := time.Parse(time.DateOnly, e.EndDay)
	if err != nil {
		return nil, ErrorInvalidEndDay
	}
	nextStart := startDay
	nextEnd := endDay

	if err := Validate(e); err != nil {
		return nil, err
	}
	r := e.Repeat

	var events []*Event

	switch e.Repeat.RepeatType {
	case RepeatTypeDaily, RepeatTypeMonthly, RepeatTypeYearly:
		events = append(events, &e)
		// daily, monthly, and yearly repeats are all the same
		// kind of repeating
		year, month, day := 0, 0, 0
		switch e.Repeat.RepeatType {
		case RepeatTypeDaily:
			day++
		case RepeatTypeMonthly:
			month++
		case RepeatTypeYearly:
			year++
		}
		if r.RepeatOccurrences >= 2 {
			// loop until there are a specific number of events
			for len(events) < int(r.RepeatOccurrences) {
				nextEvent := e

				nextStart = nextStart.AddDate(year, month, day)
				nextEnd = nextEnd.AddDate(year, month, day)

				nextEvent.StartDay = nextStart.Format(time.DateOnly)
				nextEvent.EndDay = nextEnd.Format(time.DateOnly)

				events = append(events, &nextEvent)
			}
		} else if r.RepeatStopDate != nil {
			// loop until the next start date is after the stop date
			for !nextStart.After(*r.RepeatStopDate) {
				// if there are more event repeats than allowed, throw error
				if len(events) > int(MaxRepeatOccurrence) {
					return nil, ErrorTooManyRepeatOccurrences
				}
				nextEvent := e

				nextStart = nextStart.AddDate(year, month, day)
				nextEnd = nextEnd.AddDate(year, month, day)

				nextEvent.StartDay = nextStart.Format(time.DateOnly)
				nextEvent.EndDay = nextEnd.Format(time.DateOnly)

				events = append(events, &nextEvent)
			}
		}
	case RepeatTypeWeekly:
		// weekly repeating happens based on the day of the week which
		// means the initial event could actually be not in the repeating
		// events. Ex: initial event is on a Wednesday, but the DayOfWeek
		// says Tuesday and Thursday.
		if r.RepeatOccurrences >= 2 {
			// loop until there are a specific number of events
			for len(events) < int(r.RepeatOccurrences) {
				day := dayOfWeekFromWeekday(nextStart.Weekday())
				if !r.DayOfWeek.HasFlag(day) {
					continue
				}

				nextEvent := e
				nextEvent.StartDay = nextStart.Format(time.DateOnly)
				nextEvent.EndDay = nextEnd.Format(time.DateOnly)
				events = append(events, &nextEvent)

				// go to the next day (do this at the end of the for loop
				// since we need to check the original event)
				nextStart = nextStart.AddDate(0, 0, 1)
				nextEnd = nextEnd.AddDate(0, 0, 1)
			}
		} else if r.RepeatStopDate != nil {
			// loop until the next start date is after the stop date
			for !nextStart.After(*r.RepeatStopDate) {
				// if there are more event repeats than allowed, throw error
				if len(events) > int(MaxRepeatOccurrence) {
					return nil, ErrorTooManyRepeatOccurrences
				}

				day := dayOfWeekFromWeekday(nextStart.Weekday())
				if !r.DayOfWeek.HasFlag(day) {
					continue
				}

				nextEvent := e
				nextEvent.StartDay = nextStart.Format(time.DateOnly)
				nextEvent.EndDay = nextEnd.Format(time.DateOnly)
				events = append(events, &nextEvent)

				// go to the next day (do this at the end of the for loop
				// since we need to check the original event)
				nextStart = nextStart.AddDate(0, 0, 1)
				nextEnd = nextEnd.AddDate(0, 0, 1)
			}
		}
	}

	if events == nil || len(events) == 0 {
		return nil, ErrorEmptyRepeatingEvents
	}

	return nil, nil
}
