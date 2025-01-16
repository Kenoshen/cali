package cali

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tt is a helper test function to convert a time string to a
// time.Time value. Expects YYYY-MM-DD HH:MM
func tt(s string) *time.Time {
	t, err := time.Parse(time.DateOnly+" "+TimeFormat, s)
	if err != nil {
		panic(err.Error())
	}
	return &t
}

func TestTT(t *testing.T) {
	assert.Equal(t, "2008-01-01 00:00", tt("2008-01-01 00:00").Format(DayTimeFormat))
	assert.Equal(t, "2009-02-03 13:45", tt("2009-02-03 13:45").Format(DayTimeFormat))
}

func TestQueryMatch(t *testing.T) {
	testCases := []struct {
		name string
		q    Query
		out  []int
	}{
		{
			name: "outside time range",
			q: Query{
				Start: tt("2006-01-01 00:00"),
				End:   tt("2007-01-01 00:00"),
			},
			out: nil,
		},
		{
			name: "all repeating",
			q: Query{
				Start:     tt("2008-01-01 00:00"),
				End:       tt("2009-01-01 00:00"),
				ParentIds: []int64{10},
			},
			out: []int{10, 11, 12},
		},
		{
			name: "all",
			q:    Query{},
			out:  []int{1, 10, 2, 11, 12},
		},
		{
			name: "on or after",
			q:    Query{Start: tt("2008-01-02 00:00")},
			out:  []int{2, 11, 12},
		},
		{
			name: "on or before",
			q:    Query{End: tt("2008-01-02 08:00")},
			out:  []int{1, 10, 2, 11},
		},
	}

	i := func(i int64) *int64 {
		return &i
	}

	events := []*Event{
		{
			Id:        1,
			StartDay:  "2008-01-01",
			StartTime: "08:00",
			EndDay:    "2008-01-01",
			EndTime:   "09:00",
		},
		{
			Id:        10,
			ParentId:  i(10),
			StartDay:  "2008-01-01",
			StartTime: "08:00",
			EndDay:    "2008-01-01",
			EndTime:   "09:00",
		},
		{
			Id:       2,
			IsAllDay: true,
			StartDay: "2008-01-02",
			EndDay:   "2008-01-02",
		},
		{
			Id:        11,
			ParentId:  i(10),
			StartDay:  "2008-01-02",
			StartTime: "08:00",
			EndDay:    "2008-01-02",
			EndTime:   "09:00",
		},
		{
			Id:        12,
			ParentId:  i(10),
			StartDay:  "2008-01-03",
			StartTime: "08:00",
			EndDay:    "2008-01-03",
			EndTime:   "09:00",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			t.Log(tc.name)

			var result []int
			for _, e := range events {
				if tc.q.Matches(e) {
					t.Log("Match:", e.StartDay, tc.q.Start)
					result = append(result, int(e.Id))
				}
			}
			assert.Equal(t, tc.out, result)
		})
	}
}

func TestParseDayTime(t *testing.T) {
	testCases := []struct {
		name    string
		day     string
		hourMin string
		out     string
		err     string
	}{
		{
			name: "no day",
			err:  "invalid day value",
		},
		{
			name: "invalid day",
			day:  "not-a-day",
			err:  "cannot parse",
		},
		{
			name: "day only",
			day:  "2008-01-01",
			out:  "2008-01-01 00:00",
		},
		{
			name:    "day and invalid time",
			day:     "2008-01-01",
			hourMin: "blah",
			err:     "cannot parse",
		},
		{
			name:    "day and time",
			day:     "2008-01-01",
			hourMin: "13:00",
			out:     "2008-01-01 13:00",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			t.Log(tc.name)
			result, err := parseDayTime(tc.day, tc.hourMin)
			if tc.err != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.out, result.Format(DayTimeFormat))
		})
	}
}
