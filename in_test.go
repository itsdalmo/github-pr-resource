package resource_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shurcooL/githubv4"
	"github.com/telia-oss/github-pr-resource"
	"github.com/telia-oss/github-pr-resource/mocks"
)

func TestGet(t *testing.T) {

	tests := []struct {
		description    string
		source         resource.Source
		version        resource.Version
		parameters     resource.GetParameters
		pullRequest    *resource.PullRequest
		versionString  string
		metadataString string
	}{
		{
			description: "get works",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:            "pr1",
				Commit:        "commit1",
				CommittedDate: time.Time{},
			},
			parameters:     resource.GetParameters{},
			pullRequest:    createTestPR(1, false),
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_sha","value":"oid1"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"}]`,
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

			// Run the get and check output
			input := resource.GetRequest{Source: tc.source, Version: tc.version, Params: tc.parameters}
			output, err := resource.Get(input, github, git, dir)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if got, want := output.Version, tc.version; !reflect.DeepEqual(got, want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
			}

			// Verify written files
			version := readTestFile(t, filepath.Join(dir, ".git", "resource", "version.json"))
			if got, want := version, tc.versionString; got != want {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
			}

			metadata := readTestFile(t, filepath.Join(dir, ".git", "resource", "metadata.json"))
			if got, want := metadata, tc.metadataString; got != want {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
			}
		})
	}
}

func TestGetSkipDownload(t *testing.T) {

	tests := []struct {
		description string
		source      resource.Source
		version     resource.Version
		parameters  resource.GetParameters
	}{
		{
			description: "skip download works",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:            "pr1",
				Commit:        "commit1",
				CommittedDate: time.Time{},
			},
			parameters: resource.GetParameters{SkipDownload: true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			github := mocks.NewMockGithub(ctrl)
			git := mocks.NewMockGit(ctrl)
			dir := createTestDirectory(t)
			defer os.RemoveAll(dir)

			// Run the get and check output
			input := resource.GetRequest{Source: tc.source, Version: tc.version, Params: tc.parameters}
			output, err := resource.Get(input, github, git, dir)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if got, want := output.Version, tc.version; !reflect.DeepEqual(got, want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
			}
		})
	}
}

func createTestPR(count int, skipCI bool) *resource.PullRequest {
	n := strconv.Itoa(count)
	d := time.Now().AddDate(0, 0, -count)
	m := fmt.Sprintf("commit message%s", n)
	if skipCI {
		m = "[skip ci]" + m
	}

	return &resource.PullRequest{
		PullRequestObject: resource.PullRequestObject{
			ID:          fmt.Sprintf("pr%s", n),
			Number:      count,
			Title:       fmt.Sprintf("pr%s title", n),
			URL:         fmt.Sprintf("pr%s url", n),
			BaseRefName: "master",
			HeadRefName: fmt.Sprintf("pr%s", n),
			Repository: struct{ URL string }{
				URL: fmt.Sprintf("repo%s url", n),
			},
		},
		Tip: resource.CommitObject{
			ID:            fmt.Sprintf("commit%s", n),
			OID:           fmt.Sprintf("oid%s", n),
			CommittedDate: githubv4.DateTime{Time: d},
			Message:       m,
			Author: struct{ User struct{ Login string } }{
				User: struct{ Login string }{
					Login: fmt.Sprintf("login%s", n),
				},
			},
		},
	}
}

func createTestDirectory(t *testing.T) string {
	dir, err := ioutil.TempDir("", "github-pr-resource")
	if err != nil {
		t.Fatalf("failed to create temporary directory")
	}
	return dir
}

func readTestFile(t *testing.T, path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read: %s: %s", path, err)
	}
	return string(b)
}
