// +build integration

package resource_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	resource "github.com/telia-oss/github-pr-resource"
)

const (
	defaultRepository = "itsdalmo/test-repository"
)

var (
	source resource.Source
)

func TestComponentCheck(t *testing.T) {
	tests := []struct {
		description string
		source      resource.Source
		version     resource.Version
		expected    resource.CheckResponse
		err         error
	}{
		{
			description: "check returns the latest version if there is no previous",
			source:      buildSource(t, false, false),
			version:     resource.Version{},
			expected: resource.CheckResponse{
				prVersion(4),
			},
			err: nil,
		},
		{
			description: "check returns latest version if latest matches input version",
			source:      buildSource(t, false, false),
			version:     prVersion(5),
			expected: resource.CheckResponse{
				prVersion(5),
			},
			err: nil,
		},
		{
			description: "check returns multiple versions sorted by UpdatedDate",
			source:      buildSource(t, false, false),
			version:     prVersion(1),
			expected: resource.CheckResponse{
				prVersion(3),
				prVersion(5),
				prVersion(4),
			},
			err: nil,
		},
		{
			description: "check enables previews",
			source:      buildSource(t, true, false),
			version:     prVersion(3),
			expected: resource.CheckResponse{
				prVersion(5),
				prVersion(4),
			},
			err: nil,
		},
		{
			description: "check disables skipci",
			source:      buildSource(t, false, true),
			version:     prVersion(1),
			expected: resource.CheckResponse{
				prVersion(3),
				prVersion(5),
				prVersion(6),
				prVersion(4),
			},
			err: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			github, err := resource.NewGithubClient(&tc.source)
			require.NoError(t, err)

			input := resource.CheckRequest{Source: tc.source, Version: tc.version}
			output, err := resource.Check(input, github)

			if assert.NoError(t, err) {
				assert.Equal(t, tc.expected, output)
			}
		})
	}
}

func TestComponentIn(t *testing.T) {

}

func TestComponentOut(t *testing.T) {

}

func buildSource(t *testing.T, preview bool, skipci bool) resource.Source {
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	if token == "" {
		t.Fatal("environment variable GITHUB_ACCESS_TOKEN is required")
	}

	repo := os.Getenv("GITHUB_REPOSITORY")
	if repo == "" {
		repo = defaultRepository
	}

	s := resource.Source{
		Repository:  repo,
		AccessToken: token,
	}

	if preview {
		s.PreviewSchema = true
	}

	if skipci {
		s.DisableCISkip = true
	}

	host := os.Getenv("GITHUB_HOST")
	if host != "" {
		s.V3Endpoint = fmt.Sprintf("%s/api/v3/", host)
		s.V4Endpoint = fmt.Sprintf("%s/api/graphql", host)
	}

	return s
}

func prVersion(number int) resource.Version {
	switch number {
	case 1:
		return resource.Version{PR: "1", Commit: "444503178704846a540b17707fc8fa238314664b", UpdatedDate: time.Date(2018, time.May, 10, 10, 52, 20, 0, time.UTC)}
	case 3:
		return resource.Version{PR: "3", Commit: "23dc9f552bf989d1a4aeb65ce23351dee0ec9019", UpdatedDate: time.Date(2018, time.May, 11, 7, 30, 57, 0, time.UTC)}
	case 4:
		return resource.Version{PR: "4", Commit: "a5114f6ab89f4b736655642a11e8d15ce363d882", UpdatedDate: time.Date(2019, time.August, 9, 9, 49, 45, 0, time.UTC)}
	case 5:
		return resource.Version{PR: "5", Commit: "890a7e4f0d5b05bda8ea21b91f4604e3e0313581", UpdatedDate: time.Date(2018, time.May, 14, 10, 52, 20, 0, time.UTC)}
	case 6:
		return resource.Version{PR: "6", Commit: "ac771f3b69cbd63b22bbda553f827ab36150c640", UpdatedDate: time.Date(2018, time.September, 25, 21, 0, 36, 0, time.UTC)}
	}

	return resource.Version{}
}
