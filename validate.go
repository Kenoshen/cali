package cali

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrorNotRepeatingEvent          = errors.New("must be a repeating event")
	ErrorRepeatOccurrenceTooLarge   = errors.New("repeat occurrences is over the maximum count")
	ErrorRepeatOccurrenceTooSmall   = errors.New("repeat occurrences must be at least 2")
	ErrorRepeatStopDateTooLarge     = errors.New("repeat stop date is over the maximum duration")
	ErrorRepeatStopDateIsBeforeNow  = errors.New("repeat stop date must be after now")
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
)

func Validate(e Event) error {
	if e.StartDay > e.EndDay {
		return ErrorStartDayIsAfterEndDay
	} else if e.StartDay == e.EndDay && e.StartTime > e.EndTime {
		return ErrorStartTimeIsAfterEndTime
	}

	_, err := time.LoadLocation(e.Zone)
	if err != nil {
		// TODO: might need to wrap this error so it can be tested easier
		return err
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
			startDay, err := time.Parse(time.DateOnly, e.StartDay)
			if err != nil {
				return ErrorInvalidStartDay
			}
			if !e.Repeat.RepeatStopDate.After(startDay) {
				return fmt.Errorf("%v: %v | %v", ErrorRepeatStopDateIsBeforeNow, e.Repeat.RepeatStopDate, startDay)
			}
			if e.Repeat.RepeatStopDate.After(startDay.Add(24 * time.Hour).Add(MaxRepeatDuration)) {
				return ErrorRepeatStopDateTooLarge
			}
		}

		switch e.Repeat.RepeatType {
		case RepeatTypeDaily:
		case RepeatTypeWeekly:
			if e.Repeat.DayOfWeek <= 0 {

			}
		case RepeatTypeMonthly:
		case RepeatTypeYearly:
		default:
			return ErrorInvalidRepeatType
		}
	}

	return nil
}
