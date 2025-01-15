package cali

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		desc string
		in   Event
		err  error
	}{
		{
			desc: "invalid start day",
			in: Event{
				StartDay: "not-a-day",
			},
			err: ErrorInvalidStartDay,
		}, {
			desc: "invalid end day",
			in: Event{
				StartDay: "2008-01-01",
				EndDay:   "not-a-day",
			},
			err: ErrorInvalidEndDay,
		}, {
			desc: "invalid start time",
			in: Event{
				StartDay:  "2008-01-01",
				EndDay:    "2008-01-01",
				StartTime: "HH:mm",
			},
			err: ErrorInvalidStartTime,
		}, {
			desc: "invalid end time",
			in: Event{
				StartDay:  "2008-01-01",
				EndDay:    "2008-01-01",
				StartTime: "13:00",
				EndTime:   "HH:mm",
			},
			err: ErrorInvalidEndTime,
		}, {
			desc: "start day is after end day",
			in: Event{
				StartDay:  "2008-01-02",
				EndDay:    "2008-01-01",
				StartTime: "13:00",
				EndTime:   "14:00",
			},
			err: ErrorStartDayIsAfterEndDay,
		}, {
			desc: "start time is after end time",
			in: Event{
				StartDay:  "2008-01-01",
				EndDay:    "2008-01-01",
				StartTime: "15:00",
				EndTime:   "14:00",
			},
			err: ErrorStartTimeIsAfterEndTime,
		}, {
			desc: "invalid zone",
			in: Event{
				StartDay:  "2008-01-01",
				EndDay:    "2008-01-01",
				StartTime: "13:00",
				EndTime:   "14:00",
				Zone:      "not-a-zone",
			},
			err: ErrorInvalidZone,
		}, {
			desc: "missing repeating pattern",
			in: Event{
				StartDay:    "2008-01-01",
				EndDay:      "2008-01-01",
				StartTime:   "13:00",
				EndTime:     "14:00",
				Zone:        "America/Denver",
				IsRepeating: true,
			},
			err: ErrorMissingRepeatPattern,
		}, {
			desc: "repeat occurrence too large",
			in: Event{
				StartDay:    "2008-01-01",
				EndDay:      "2008-01-01",
				StartTime:   "13:00",
				EndTime:     "14:00",
				Zone:        "America/Denver",
				IsRepeating: true,
				Repeat:      &Repeat{RepeatOccurrences: 10000},
			},
			err: ErrorRepeatOccurrenceTooLarge,
		}, {
			desc: "repeat occurrence too small",
			in: Event{
				StartDay:    "2008-01-01",
				EndDay:      "2008-01-01",
				StartTime:   "13:00",
				EndTime:     "14:00",
				Zone:        "America/Denver",
				IsRepeating: true,
				Repeat:      &Repeat{RepeatOccurrences: 1},
			},
			err: ErrorRepeatOccurrenceTooSmall,
		}, {
			desc: "missing end of repeat",
			in: Event{
				StartDay:    "2008-01-01",
				EndDay:      "2008-01-01",
				StartTime:   "13:00",
				EndTime:     "14:00",
				Zone:        "America/Denver",
				IsRepeating: true,
				Repeat:      &Repeat{RepeatOccurrences: 0, RepeatStopDate: nil},
			},
			err: ErrorMissingEndOfRepeat,
		}, {
			desc: "repeat stop date is before start",
			in: Event{
				StartDay:    "2008-01-01",
				EndDay:      "2008-01-01",
				StartTime:   "13:00",
				EndTime:     "14:00",
				Zone:        "America/Denver",
				IsRepeating: true,
				Repeat:      &Repeat{RepeatStopDate: _t(time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC))},
			},
			err: ErrorRepeatStopDateIsBeforeStart,
		}, {
			desc: "repeat stop date too large",
			in: Event{
				StartDay:    "2008-01-01",
				EndDay:      "2008-01-01",
				StartTime:   "13:00",
				EndTime:     "14:00",
				Zone:        "America/Denver",
				IsRepeating: true,
				Repeat:      &Repeat{RepeatStopDate: _t(time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC))},
			},
			err: ErrorRepeatStopDateTooLarge,
		}, {
			desc: "invalid day of the week",
			in: Event{
				StartDay:    "2008-01-01",
				EndDay:      "2008-01-01",
				StartTime:   "13:00",
				EndTime:     "14:00",
				Zone:        "America/Denver",
				IsRepeating: true,
				Repeat:      &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: 0, RepeatStopDate: _t(time.Date(2008, time.January, 20, 0, 0, 0, 0, time.UTC))},
			},
			err: ErrorInvalidDayOfWeek,
		}, {
			desc: "invalid repeat type",
			in: Event{
				StartDay:    "2008-01-01",
				EndDay:      "2008-01-01",
				StartTime:   "13:00",
				EndTime:     "14:00",
				Zone:        "America/Denver",
				IsRepeating: true,
				Repeat:      &Repeat{RepeatType: -1, DayOfWeek: 0, RepeatStopDate: _t(time.Date(2008, time.January, 20, 0, 0, 0, 0, time.UTC))},
			},
			err: ErrorInvalidRepeatType,
		}, {
			desc: "success",
			in: Event{
				StartDay:    "2008-01-01",
				EndDay:      "2008-01-01",
				StartTime:   "13:00",
				EndTime:     "14:00",
				Zone:        "America/Denver",
				IsRepeating: true,
				Repeat:      &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekTuesday, RepeatStopDate: _t(time.Date(2008, time.January, 20, 0, 0, 0, 0, time.UTC))},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(tc.desc)
			err := Validate(tc.in)
			if tc.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateInvite(t *testing.T) {
	testCases := []struct {
		desc string
		in   Invite
		err  error
	}{
		{
			desc: "invalid invite status",
			in: Invite{
				Status: InviteStatus(-1),
			},
			err: ErrorInvalidInviteStatus,
		}, {
			desc: "missing invite permission",
			in: Invite{
				Permission: 0,
			},
			err: ErrorMissingInvitePermission,
		}, {
			desc: "missing read permission",
			in: Invite{
				Permission: PermissionModify | PermissionCancel,
			},
			err: ErrorIncompatibleInvitePermission,
		}, {
			desc: "missing invite permission",
			in: Invite{
				Permission: PermissionRead | PermissionModify,
			},
			err: ErrorIncompatibleInvitePermission,
		}, {
			desc: "missing modify permission",
			in: Invite{
				Permission: PermissionRead | PermissionInvite | PermissionCancel | PermissionDelete,
			},
			err: ErrorIncompatibleInvitePermission,
		}, {
			desc: "missing cancel permission",
			in: Invite{
				Permission: PermissionRead | PermissionInvite | PermissionModify | PermissionDelete,
			},
			err: ErrorIncompatibleInvitePermission,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(tc.desc)
			err := ValidateInvite(tc.in)
			if tc.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
