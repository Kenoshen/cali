package cali

import (
	"errors"
	"time"
)

var (
	ErrorNotRepeatingEvent            = errors.New("must be a repeating event")
	ErrorRepeatOccurrenceTooLarge     = errors.New("repeat occurrences is over the maximum count")
	ErrorRepeatOccurrenceTooSmall     = errors.New("repeat occurrences must be at least 2")
	ErrorRepeatStopDateTooLarge       = errors.New("repeat stop date is over the maximum duration")
	ErrorRepeatStopDateIsBeforeStart  = errors.New("repeat stop date must be after start")
	ErrorStartDayIsAfterEndDay        = errors.New("start day must be equal or less than end day")
	ErrorStartTimeIsAfterEndTime      = errors.New("start time must be equal or less than end time")
	ErrorMissingEndOfRepeat           = errors.New("repeating events must have some end")
	ErrorEmptyRepeatingEvents         = errors.New("repeating event list is empty")
	ErrorMissingRepeatPattern         = errors.New("missing repeat pattern")
	ErrorInvalidRepeatType            = errors.New("invalid repeat type")
	ErrorSeparationCountLessThanOne   = errors.New("separation count must be 1 or greater")
	ErrorMissingDayOfWeek             = errors.New("missing day of the week (SMTWTFS)")
	ErrorInvalidStartDay              = errors.New("invalid start day")
	ErrorInvalidStartTime             = errors.New("invalid start time")
	ErrorInvalidEndDay                = errors.New("invalid end day")
	ErrorInvalidEndTime               = errors.New("invalid end time")
	ErrorTooManyRepeatOccurrences     = errors.New("too many event occurrences in repeat calculation")
	ErrorInvalidDayOfWeek             = errors.New("invalid day of week")
	ErrorInvalidZone                  = errors.New("invalid zone")
	ErrorInvalidInviteStatus          = errors.New("invalid invite status")
	ErrorMissingInvitePermission      = errors.New("missing invite permission")
	ErrorIncompatibleInvitePermission = errors.New("incompatible invite permission")
	ErrorEventNotFound                = errors.New("there is no event with that id")
	ErrorInvalidStatus                = errors.New("invalid status")
	ErrorInviteNotFound               = errors.New("invitation not found")
	ErrorInvalidRepeatEditType        = errors.New("invalid repeat edit type")
)

// VAlidate makes sure the event object doesn't have conflicting values
func Validate(e Event) error {
	if err := ValidTimes(e.StartDay, e.StartTime, e.EndDay, e.EndTime, e.Zone, e.IsAllDay); err != nil {
		return err
	}

	if err := ValidRepeat(e); err != nil {
		return err
	}

	if !ValidStatus(e.Status) {
		return ErrorInvalidStatus
	}

	return nil
}

// ValidateInvite makes sure the invite object doesn't have conflicting values
func ValidateInvite(a Invite) error {
	switch a.Status {
	case InviteStatusPending, InviteStatusConfirmed, InviteStatusDeclined:
	default:
		return ErrorInvalidInviteStatus
	}

	if a.Permission <= 0 {
		return ErrorMissingInvitePermission
	}

	if !a.Permission.HasFlag(PermissionRead) && (a.Permission.HasFlag(PermissionDelete) || a.Permission.HasFlag(PermissionCancel) || a.Permission.HasFlag(PermissionInvite) || a.Permission.HasFlag(PermissionModify)) {
		return ErrorIncompatibleInvitePermission
	}

	if !a.Permission.HasFlag(PermissionInvite) && a.Permission.HasFlag(PermissionModify) {
		return ErrorIncompatibleInvitePermission
	}

	if !a.Permission.HasFlag(PermissionModify) && (a.Permission.HasFlag(PermissionDelete) || a.Permission.HasFlag(PermissionCancel)) {
		return ErrorIncompatibleInvitePermission
	}

	if !a.Permission.HasFlag(PermissionCancel) && a.Permission.HasFlag(PermissionDelete) {
		return ErrorIncompatibleInvitePermission
	}

	return nil
}

// ValidStatus returns true if the status is one of the pre-defined statuses from this library
func ValidStatus(s Status) bool {
	switch s {
	case StatusActive, StatusCanceled, StatusAbandoned, StatusRemoved:
		return true
	default:
		return false
	}
}

// ValidRepeat checks the event.Repeat if event.IsRepeating is true to see if there are invalid values within the repeat
func ValidRepeat(e Event) error {
	if e.IsRepeating {
		startDay, err := time.Parse(time.DateOnly, e.StartDay)
		if err != nil {
			return ErrorInvalidStartDay
		}
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
			// allows stop date to be equal to start day since stop date
			// is inclusive
			if e.Repeat.RepeatStopDate.Before(startDay) {
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

// ValidateTimeValues compares two HH:mm values to make sure they are
// correctly formatted and start time is equal or before the end time
func ValidateTimeValues(startTime, endTime string) error {
	_, err := time.Parse(TimeFormat, startTime)
	if err != nil {
		return ErrorInvalidStartTime
	}
	_, err = time.Parse(TimeFormat, endTime)
	if err != nil {
		return ErrorInvalidEndTime
	}
	if startTime > endTime {
		return ErrorStartTimeIsAfterEndTime
	}
	return nil
}

// ValidateDayValues compares two YYYY-MM-DD values to make sure they are
// correctly formatted and start day is equal or before the end day
func ValidateDayValues(startDay, endDay string) error {
	_, err := time.Parse(time.DateOnly, startDay)
	if err != nil {
		return ErrorInvalidStartDay
	}
	_, err = time.Parse(time.DateOnly, endDay)
	if err != nil {
		return ErrorInvalidEndDay
	}
	if startDay > endDay {
		return ErrorStartDayIsAfterEndDay
	}
	return nil
}

// ValidateDayTimeValues makes sure that the start and end dates and times are valid values
func ValidateDayTimeValues(startDay, startTime, endDay, endTime string) error {
	_, err := time.Parse(time.DateOnly, startDay)
	if err != nil {
		return ErrorInvalidStartDay
	}
	_, err = time.Parse(time.DateOnly, endDay)
	if err != nil {
		return ErrorInvalidEndDay
	}
	_, err = time.Parse(TimeFormat, startTime)
	if err != nil {
		return ErrorInvalidStartTime
	}
	_, err = time.Parse(TimeFormat, endTime)
	if err != nil {
		return ErrorInvalidEndTime
	}
	if startDay > endDay {
		return ErrorStartDayIsAfterEndDay
	} else if startDay == endDay && startTime > endTime {
		return ErrorStartTimeIsAfterEndTime
	}


	return nil
}
