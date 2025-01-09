package cali

import (
	"errors"
	"time"
)

var (
	ErrorNotRepeatingEvent          = errors.New("must be a repeating event")
	ErrorRepeatOccurrenceTooLarge   = errors.New("repeat occurrences is over the maximum count")
	ErrorRepeatOccurrenceTooSmall   = errors.New("repeat occurrences must be at least 2")
	ErrorRepeatStopDateTooLarge     = errors.New("repeat stop date is over the maximum duration")
	ErrorRepeatStopDateIsBeforeStart  = errors.New("repeat stop date must be after start")
	ErrorStartDayIsAfterEndDay      = errors.New("start day must be equal or less than end day")
	ErrorStartTimeIsAfterEndTime    = errors.New("start time must be equal or less than end time")
	ErrorMissingEndOfRepeat         = errors.New("repeating events must have some end")
	ErrorEmptyRepeatingEvents       = errors.New("repeating event list is empty")
	ErrorMissingRepeatPattern       = errors.New("missing repeat pattern")
	ErrorInvalidRepeatType          = errors.New("invalid repeat type")
	ErrorSeparationCountLessThanOne = errors.New("separation count must be 1 or greater")
	ErrorMissingDayOfWeek           = errors.New("missing day of the week (SMTWTFS)")
	ErrorInvalidStartDay            = errors.New("invalid start day")
	ErrorInvalidStartTime           = errors.New("invalid start time")
	ErrorInvalidEndDay              = errors.New("invalid end day")
	ErrorInvalidEndTime             = errors.New("invalid end time")
	ErrorTooManyRepeatOccurrences   = errors.New("too many event occurrences in repeat calculation")
	ErrorInvalidDayOfWeek           = errors.New("invalid day of week")
	ErrorInvalidZone           = errors.New("invalid zone")
)

func Validate(e Event) error {
	startDay, err := time.Parse(time.DateOnly, e.StartDay)
	if err != nil {
		return ErrorInvalidStartDay
	}
	_, err = time.Parse(time.DateOnly, e.EndDay)
	if err != nil {
		return ErrorInvalidEndDay
	}
	if !e.IsAllDay {
		_, err = time.Parse(TimeFormat, e.StartTime)
		if err != nil {
			return ErrorInvalidStartTime
		}
		_, err = time.Parse(TimeFormat, e.EndTime)
		if err != nil {
			return ErrorInvalidEndTime
		}
	}
	if e.StartDay > e.EndDay {
		return ErrorStartDayIsAfterEndDay
	} else if e.StartDay == e.EndDay && e.StartTime > e.EndTime {
		return ErrorStartTimeIsAfterEndTime
	}

	_, err = time.LoadLocation(e.Zone)
	if err != nil {
    return ErrorInvalidZone
	}

	if e.IsRepeating {
		if e.Repeat == nil {
			return ErrorMissingRepeatPattern
		}
		if e.Repeat.RepeatOccurrences > MaxRepeatOccurrence {
			return ErrorRepeatOccurrenceTooLarge
		}
		if e.Repeat.RepeatOccurrences == 1 || e.Repeat.RepeatOccurrences < 0 {
			return ErrorRepeatOccurrenceTooSmall
		}
		if e.Repeat.RepeatStopDate == nil && e.Repeat.RepeatOccurrences == 0 {
			return ErrorMissingEndOfRepeat
		}
		if e.Repeat.RepeatStopDate != nil {
			if !e.Repeat.RepeatStopDate.After(startDay) {
				return ErrorRepeatStopDateIsBeforeStart
      }
			if e.Repeat.RepeatStopDate.After(startDay.Add(24 * time.Hour).Add(MaxRepeatDuration)) {
				return ErrorRepeatStopDateTooLarge
			}
		}

		switch e.Repeat.RepeatType {
		case RepeatTypeDaily:
		case RepeatTypeWeekly:
			if e.Repeat.DayOfWeek <= 0 {
				return ErrorInvalidDayOfWeek
			}
		case RepeatTypeMonthly:
		case RepeatTypeYearly:
		default:
			return ErrorInvalidRepeatType
		}
	}

	return nil
}
