package resource_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/shurcooL/githubv4"
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
			pullRequest: createTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
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
			pullRequest: createTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
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
			pullRequest: createTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can provide a custom base context for the status",
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
				Status:      "failure",
				BaseContext: "concourse-ci-custom",
				Context:     "build",
			},
			pullRequest: createTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can provide a custom target url for the status",
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
				Status:    "failure",
				TargetURL: "https://targeturl.com/concourse",
			},
			pullRequest: createTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can provide a custom description for the status",
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
				Status:      "failure",
				Description: "Concourse CI build",
			},
			pullRequest: createTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
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
			pullRequest: createTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can delete previous comments made on the pull request",
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
				DeletePreviousComments: true,
			},
			pullRequest: createTestPR(1, "master", false, false, 0, []string{}, false, githubv4.PullRequestStateOpen),
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
					commit, baseContext, context, status, targetURL, description := github.UpdateCommitStatusArgsForCall(0)
					assert.Equal(t, tc.version.Commit, commit)
					assert.Equal(t, tc.parameters.BaseContext, baseContext)
					assert.Equal(t, tc.parameters.Context, context)
					assert.Equal(t, tc.parameters.TargetURL, targetURL)
					assert.Equal(t, tc.parameters.Description, description)
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

			if tc.parameters.DeletePreviousComments {
				if assert.Equal(t, 1, github.DeletePreviousCommentsCallCount()) {
					pr := github.DeletePreviousCommentsArgsForCall(0)
					assert.Equal(t, tc.version.PR, pr)
				}
			}
		})
	}
}

func TestVariableSubstitution(t *testing.T) {

	var (
		variableName  = "BUILD_JOB_NAME"
		variableValue = "my-job"
		variableURL   = "https://concourse-ci.org/"
	)

	tests := []struct {
		description       string
		source            resource.Source
		version           resource.Version
		parameters        resource.PutParameters
		expectedComment   string
		expectedTargetURL string
		pullRequest       *resource.PullRequest
	}{

		{
			description: "we can substitute environment variables for Comment",
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
				Comment: fmt.Sprintf("$%s", variableName),
			},
			expectedComment: variableValue,
			pullRequest:     createTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can substitute environment variables for TargetURL",
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
				Status:    "failure",
				TargetURL: fmt.Sprintf("%s$%s", variableURL, variableName),
			},
			expectedTargetURL: fmt.Sprintf("%s%s", variableURL, variableValue),
			pullRequest:       createTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we do not substitute variables other then concourse build metadata",
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
				Comment: "$THIS_IS_NOT_SUBSTITUTED",
			},
			expectedComment: "$THIS_IS_NOT_SUBSTITUTED",
			pullRequest:     createTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
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
			getInput := resource.GetRequest{Source: tc.source, Version: tc.version, Params: resource.GetParameters{}}
			_, err := resource.Get(getInput, github, git, dir)
			require.NoError(t, err)

			oldValue := os.Getenv(variableName)
			defer os.Setenv(variableName, oldValue)

			os.Setenv(variableName, variableValue)

			putInput := resource.PutRequest{Source: tc.source, Params: tc.parameters}
			_, err = resource.Put(putInput, github, dir)

			if tc.parameters.TargetURL != "" {
				if assert.Equal(t, 1, github.UpdateCommitStatusCallCount()) {
					_, _, _, _, targetURL, _ := github.UpdateCommitStatusArgsForCall(0)
					assert.Equal(t, tc.expectedTargetURL, targetURL)
				}
			}

			if tc.parameters.Comment != "" {
				if assert.Equal(t, 1, github.PostCommentCallCount()) {
					_, comment := github.PostCommentArgsForCall(0)
					assert.Equal(t, tc.expectedComment, comment)
				}
			}

		})
	}
}
