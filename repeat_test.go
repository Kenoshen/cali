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
	}{{
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
			StartDay:    "2008-01-01",
			EndDay:      "2008-01-01",
			Repeat: &Repeat{
				RepeatType:        RepeatTypeDaily,
				RepeatOccurrences: 3,
			},
		},
		out: []*Event{{
			IsRepeating: true,
			StartDay:    "2008-01-01",
			EndDay:      "2008-01-01",
			Repeat: &Repeat{
				RepeatType:        RepeatTypeDaily,
				RepeatOccurrences: 3,
			},
		}, {
			IsRepeating: true,
			StartDay:    "2008-01-02",
			EndDay:      "2008-01-02",
			Repeat: &Repeat{
				RepeatType:        RepeatTypeDaily,
				RepeatOccurrences: 3,
			},
		}, {
			IsRepeating: true,
			StartDay:    "2008-01-03",
			EndDay:      "2008-01-03",
			Repeat: &Repeat{
				RepeatType:        RepeatTypeDaily,
				RepeatOccurrences: 3,
			},
		}},
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			out, err := GenerateRepeatEvents(tc.in)
			if tc.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, len(tc.out), len(out), "event length is different")
			for i, a := range tc.out {
				b := out[i]
				assert.Equal(t, a, *b)
			}
		})
	}
}
