package resource_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	resource "github.com/telia-oss/github-pr-resource"
	"github.com/telia-oss/github-pr-resource/fakes"
)

func TestPut(t *testing.T) {

	tests := []struct {
		description string
		source      resource.Source
		version     resource.Version
		parameters  resource.PutParameters
		pullRequest *resource.PullRequest
	}{
		{
			description: "put with no parameters does nothing",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:            "pr1",
				Commit:        "commit1",
				CommittedDate: time.Time{},
			},
			parameters:  resource.PutParameters{},
			pullRequest: createTestPR(1, false, false),
		},

		{
			description: "we can set status on a commit",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:            "pr1",
				Commit:        "commit1",
				CommittedDate: time.Time{},
			},
			parameters: resource.PutParameters{
				Status: "success",
			},
			pullRequest: createTestPR(1, false, false),
		},

		{
			description: "we can provide a custom context for the status",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:            "pr1",
				Commit:        "commit1",
				CommittedDate: time.Time{},
			},
			parameters: resource.PutParameters{
				Status:  "failure",
				Context: "build",
			},
			pullRequest: createTestPR(1, false, false),
		},

		{
			description: "we can comment on the pull request",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:            "pr1",
				Commit:        "commit1",
				CommittedDate: time.Time{},
			},
			parameters: resource.PutParameters{
				Comment: "comment",
			},
			pullRequest: createTestPR(1, false, false),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			github := new(fakes.FakeGithub)
			github.GetPullRequestReturns(tc.pullRequest, nil)

			git := new(fakes.FakeGit)
			git.RevParseReturns("sha", nil)

			dir := createTestDirectory(t)
			defer os.RemoveAll(dir)

			// Run get so we have version and metadata for the put request
			// (This is tested in in_test.go)
			getInput := resource.GetRequest{Source: tc.source, Version: tc.version, Params: resource.GetParameters{}}
			_, err := resource.Get(getInput, github, git, dir)
			require.NoError(t, err)

			putInput := resource.PutRequest{Source: tc.source, Params: tc.parameters}
			output, err := resource.Put(putInput, github, dir)

			// Validate output
			if assert.NoError(t, err) {
				assert.Equal(t, tc.version, output.Version)
			}

			// Validate method calls put on Github.
			if tc.parameters.Status != "" {
				if assert.Equal(t, 1, github.UpdateCommitStatusCallCount()) {
					commit, context, status := github.UpdateCommitStatusArgsForCall(0)
					assert.Equal(t, tc.version.Commit, commit)
					assert.Equal(t, tc.parameters.Context, context)
					assert.Equal(t, tc.parameters.Status, status)
				}
			}
			if tc.parameters.Comment != "" {
				if assert.Equal(t, 1, github.PostCommentCallCount()) {
					pr, comment := github.PostCommentArgsForCall(0)
					assert.Equal(t, tc.version.PR, pr)
					assert.Equal(t, tc.parameters.Comment, comment)
				}
			}
		})
	}
}
