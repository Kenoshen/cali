package cali

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRepeatEvent(t *testing.T) {
	testCases := []struct {
		desc string
		in   Event
		out  []*Event
		err  error
	}{
		{
			desc: "not repeating event",
			in: Event{
				IsRepeating: false,
			},
			err: ErrorNotRepeatingEvent,
		}, {
			desc: "failed to parse start day",
			in: Event{
				IsRepeating: true,
				StartDay:    "not-a-day",
				Repeat:      &Repeat{},
			},
			err: ErrorInvalidStartDay,
		}, {
			desc: "failed to parse end day",
			in: Event{
				IsRepeating: true,
				StartDay:    "2008-01-01",
				EndDay:      "",
				Repeat:      &Repeat{},
			},
			err: ErrorInvalidEndDay,
		}, {
			desc: "empty events",
			in: Event{
				IsRepeating: true,
				StartDay:    "2008-01-01",
				EndDay:      "2008-01-01",
				Repeat: &Repeat{
					RepeatType:     RepeatTypeWeekly,
					RepeatStopDate: _t(time.Date(2008, time.January, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
			err: ErrorEmptyRepeatingEvents,
		}, {
			desc: "daily 3 times",
			in: Event{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeDaily, RepeatOccurrences: 3},
			},
			out: []*Event{{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeDaily, RepeatOccurrences: 3},
			}, {
				IsRepeating: true,
				StartDay:    "2008-01-02", EndDay: "2008-01-02",
				Repeat: &Repeat{RepeatType: RepeatTypeDaily, RepeatOccurrences: 3},
			}, {
				IsRepeating: true,
				StartDay:    "2008-01-03", EndDay: "2008-01-03",
				Repeat: &Repeat{RepeatType: RepeatTypeDaily, RepeatOccurrences: 3},
			}},
		}, {
			desc: "daily too many times",
			in: Event{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeDaily, RepeatStopDate: _t(time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC))},
			},
			err: ErrorTooManyRepeatOccurrences,
		}, {
			desc: "weekly too many times",
			in: Event{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekSunday, RepeatStopDate: _t(time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC))},
			},
			err: ErrorTooManyRepeatOccurrences,
		}, {
			desc: "monthly 3 times",
			in: Event{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeMonthly, RepeatOccurrences: 3},
			},
			out: []*Event{{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeMonthly, RepeatOccurrences: 3},
			}, {
				IsRepeating: true,
				StartDay:    "2008-02-01", EndDay: "2008-02-01",
				Repeat: &Repeat{RepeatType: RepeatTypeMonthly, RepeatOccurrences: 3},
			}, {
				IsRepeating: true,
				StartDay:    "2008-03-01", EndDay: "2008-03-01",
				Repeat: &Repeat{RepeatType: RepeatTypeMonthly, RepeatOccurrences: 3},
			}},
		}, {
			desc: "yearly 3 times",
			in: Event{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeYearly, RepeatOccurrences: 3},
			},
			out: []*Event{{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeYearly, RepeatOccurrences: 3},
			}, {
				IsRepeating: true,
				StartDay:    "2009-01-01", EndDay: "2009-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeYearly, RepeatOccurrences: 3},
			}, {
				IsRepeating: true,
				StartDay:    "2010-01-01", EndDay: "2010-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeYearly, RepeatOccurrences: 3},
			}},
		}, {
			desc: "weekly 3 times on Tuesday",
			in: Event{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekTuesday, RepeatOccurrences: 3},
			},
			out: []*Event{{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekTuesday, RepeatOccurrences: 3},
			}, {
				IsRepeating: true,
				StartDay:    "2008-01-08", EndDay: "2008-01-08",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekTuesday, RepeatOccurrences: 3},
			}, {
				IsRepeating: true,
				StartDay:    "2008-01-15", EndDay: "2008-01-15",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekTuesday, RepeatOccurrences: 3},
			}},
		}, {
			desc: "weekly 5 times on Wednesday and Thursday",
			in: Event{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekWednesday | DayOfWeekThursday, RepeatOccurrences: 5},
			},
			out: []*Event{{
				IsRepeating: true,
				StartDay:    "2008-01-02", EndDay: "2008-01-02",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekWednesday | DayOfWeekThursday, RepeatOccurrences: 5},
			}, {
				IsRepeating: true,
				StartDay:    "2008-01-03", EndDay: "2008-01-03",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekWednesday | DayOfWeekThursday, RepeatOccurrences: 5},
			}, {
				IsRepeating: true,
				StartDay:    "2008-01-09", EndDay: "2008-01-09",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekWednesday | DayOfWeekThursday, RepeatOccurrences: 5},
			}, {
				IsRepeating: true,
				StartDay:    "2008-01-10", EndDay: "2008-01-10",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekWednesday | DayOfWeekThursday, RepeatOccurrences: 5},
			}, {
				IsRepeating: true,
				StartDay:    "2008-01-16", EndDay: "2008-01-16",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekWednesday | DayOfWeekThursday, RepeatOccurrences: 5},
			}},
		}, {
			desc: "repeat on Thursday but stop on Wednesday",
			in: Event{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekThursday, RepeatStopDate: _t(time.Date(2008, time.January, 2, 0, 0, 0, 0, time.UTC))},
			},
			err: ErrorEmptyRepeatingEvents,
		}, {
			desc: "weekly 5 times on Wednesday and Thursday stop on date",
			in: Event{
				IsRepeating: true,
				StartDay:    "2008-01-01", EndDay: "2008-01-01",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekWednesday | DayOfWeekThursday, RepeatStopDate: _t(time.Date(2008, time.January, 16, 0, 0, 0, 0, time.UTC))},
			},
			out: []*Event{{
				IsRepeating: true,
				StartDay:    "2008-01-02", EndDay: "2008-01-02",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekWednesday | DayOfWeekThursday, RepeatStopDate: _t(time.Date(2008, time.January, 16, 0, 0, 0, 0, time.UTC))},
			}, {
				IsRepeating: true,
				StartDay:    "2008-01-03", EndDay: "2008-01-03",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekWednesday | DayOfWeekThursday, RepeatStopDate: _t(time.Date(2008, time.January, 16, 0, 0, 0, 0, time.UTC))},
			}, {
				IsRepeating: true,
				StartDay:    "2008-01-09", EndDay: "2008-01-09",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekWednesday | DayOfWeekThursday, RepeatStopDate: _t(time.Date(2008, time.January, 16, 0, 0, 0, 0, time.UTC))},
			}, {
				IsRepeating: true,
				StartDay:    "2008-01-10", EndDay: "2008-01-10",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekWednesday | DayOfWeekThursday, RepeatStopDate: _t(time.Date(2008, time.January, 16, 0, 0, 0, 0, time.UTC))},
			}, {
				IsRepeating: true,
				StartDay:    "2008-01-16", EndDay: "2008-01-16",
				Repeat: &Repeat{RepeatType: RepeatTypeWeekly, DayOfWeek: DayOfWeekWednesday | DayOfWeekThursday, RepeatStopDate: _t(time.Date(2008, time.January, 16, 0, 0, 0, 0, time.UTC))},
			}},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(tc.desc)
			out, err := GenerateRepeatEvents(tc.in)
			if tc.err != nil {
				require.Errorf(t, err, "instead, got: %v", out)
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, len(tc.out), len(out), "event length is different")
			for i, a := range tc.out {
				b := out[i]
				require.NotNil(t, a, "expected event is nil")
				require.NotNil(t, b, "actual event is nil")
				assert.Equal(t, *a, *b)
			}
		})
	}
}
