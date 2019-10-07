package resource_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	resource "github.com/telia-oss/github-pr-resource"
	"github.com/telia-oss/github-pr-resource/fakes"
	"github.com/telia-oss/github-pr-resource/pullrequest"
)

var (
	testPullRequests = []pullrequest.PullRequest{
		// earliest
		createTestPR(1, "master", true, false, false, false, 0, nil),
		createTestPR(2, "master", false, false, false, false, 0, nil),
		// new pr
		createTestPR(3, "master", false, false, true, false, 0, nil),
		// new pr w/old commitdate
		createTestPR(4, "master", false, false, true, true, 0, nil),
		// old pr w/old commitdate
		createTestPR(5, "master", false, true, false, true, 0, nil),
		createTestPR(6, "master", false, false, false, false, 0, nil),
		createTestPR(7, "develop", false, false, false, false, 0, []string{"enhancement"}),
		createTestPR(8, "master", true, false, false, false, 1, []string{"wontfix"}),
		createTestPR(9, "master", false, false, false, false, 0, nil),
		// latest
	}
)

func TestCheck(t *testing.T) {
	tests := []struct {
		description  string
		source       resource.Source
		version      resource.Version
		files        [][]string
		pullRequests []pullrequest.PullRequest
		expected     resource.CheckResponse
	}{
		{
			description: "check returns the latest version if there is no previous",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version:      resource.Version{},
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[8]),
			},
		},
		{
			description: "check returns the latest version if there is no previous w/basebranch",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
				BaseBranch:  "master",
			},
			version:      resource.Version{},
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[8]),
			},
		},
		{
			description: "check supports specifying base branch",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
				BaseBranch:  "develop",
			},
			version:      resource.Version{},
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[6]),
			},
		},

		{
			description: "check correctly ignores PRs with no approved reviews when specified",
			source: resource.Source{
				Repository:              "itsdalmo/test-repository",
				AccessToken:             "oauthtoken",
				RequiredReviewApprovals: 1,
			},
			version:      resource.NewVersion(testPullRequests[8]),
			pullRequests: testPullRequests,
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[8]),
			},
		},
		{
			description: "check returns latest version from a PR with at least one of the desired labels on it",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
				Labels:      []string{"enhancement"},
			},
			version:      resource.Version{},
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[6]),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			github := new(fakes.FakeGithub)

			github.ListOpenPullRequestsReturns(tc.pullRequests, nil)

			for i, file := range tc.files {
				github.GetChangedFilesReturnsOnCall(i, file, nil)
			}

			input := resource.CheckRequest{Source: tc.source, Version: tc.version}
			output, err := resource.Check(input, github)

			if assert.NoError(t, err) {
				assert.Equal(t, tc.expected, output)
			}
			assert.Equal(t, 1, github.ListOpenPullRequestsCallCount())
		})
	}
}
