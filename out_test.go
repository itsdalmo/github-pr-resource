package resource_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/telia-oss/github-pr-resource"
	"github.com/telia-oss/github-pr-resource/mocks"
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			github := mocks.NewMockGithub(ctrl)
			github.EXPECT().GetPullRequest(tc.version.PR, tc.version.Commit).Times(1).Return(tc.pullRequest, nil)

			git := mocks.NewMockGit(ctrl)
			gomock.InOrder(
				git.EXPECT().Init(tc.pullRequest.BaseRefName).Times(1).Return(nil),
				git.EXPECT().Pull(tc.pullRequest.Repository.URL, tc.pullRequest.BaseRefName).Times(1).Return(nil),
				git.EXPECT().RevParse(tc.pullRequest.BaseRefName).Times(1).Return("sha", nil),
				git.EXPECT().Fetch(tc.pullRequest.Repository.URL, tc.pullRequest.Number).Times(1).Return(nil),
				git.EXPECT().Merge(tc.pullRequest.Tip.OID).Times(1).Return(nil),
			)

			dir := createTestDirectory(t)
			defer os.RemoveAll(dir)

			// Run get so we have version and metadata for the put request
			getInput := resource.GetRequest{Source: tc.source, Version: tc.version, Params: resource.GetParameters{}}
			_, err := resource.Get(getInput, github, git, dir)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			// Set expectations
			if tc.parameters.Status != "" {
				github.EXPECT().UpdateCommitStatus(tc.version.Commit, tc.parameters.Context, tc.parameters.Status).Times(1).Return(nil)
			}
			if tc.parameters.Comment != "" {
				github.EXPECT().PostComment(tc.version.PR, tc.parameters.Comment).Times(1).Return(nil)
			}

			// Run put and verify output
			putInput := resource.PutRequest{Source: tc.source, Params: tc.parameters}
			output, err := resource.Put(putInput, github, dir)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if got, want := output.Version, tc.version; !reflect.DeepEqual(got, want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
			}
		})
	}
}
