package pullrequest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLatest(t *testing.T) {
	tests := []struct {
		description string
		input       []time.Time
		expect      time.Time
	}{
		{
			description: "simple test w/3 input arguments",
			input: []time.Time{
				time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).AddDate(0, 1, 0),
				time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).AddDate(0, -1, 0),
			},
			expect: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).AddDate(0, 1, 0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := latest(tc.input...)
			assert.Equal(t, out, tc.expect)
		})
	}
}
