package cali

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// tt is a helper test function to convert a time string to a
// time.Time value. Expects YYYY-MM-DD HH:MM
func tt(s string) time.Time {
	t, _ := time.Parse(time.DateOnly+" "+TimeFormat, s)
	return t
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
				Start: tt("2007-01-01 00:00"),
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
		// TODO: add more tests cases
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
		// TODO: add more events to test with
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			t.Log(tc.name)

			var result []int
			for _, e := range events {
				if tc.q.Matches(e) {
					result = append(result, int(e.Id))
				}
			}
			assert.Equal(t, tc.out, result)
		})
	}
}
