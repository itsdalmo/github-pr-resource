package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendDepth(t *testing.T) {
	tests := []struct {
		description string
		args        []string
		depth       int
		expected    []string
	}{
		{
			description: "git clone depth 1",
			args:        []string{"git", "clone"},
			depth:       1,
			expected:    []string{"git", "clone", "--depth", "1"},
		},
		{
			description: "git clone no depth",
			args:        []string{"git", "clone"},
			depth:       0,
			expected:    []string{"git", "clone"},
		},
		{
			description: "no args no depth",
			args:        []string{},
			depth:       0,
			expected:    []string{},
		},
		{
			description: "no args with depth",
			args:        []string{},
			depth:       5,
			expected:    []string{"--depth", "5"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			output := appendDepth(tc.args, tc.depth)
			assert.Equal(t, tc.expected, output)
		})
	}
}
