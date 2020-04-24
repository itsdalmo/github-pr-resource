package resource_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	resource "github.com/telia-oss/github-pr-resource"
	"github.com/telia-oss/github-pr-resource/fakes"
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
		files          []resource.ChangedFileObject
		filesString    string
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
			pullRequest:    createTestPR(1, "master", false, false, 0, nil),
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"author_email","value":"user@example.com"},{"name":"title","value":"pr1 title"}]`,
		},
		{
			description: "get supports unlocking with git crypt",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
				GitCryptKey: "gitcryptkey",
			},
			version: resource.Version{
				PR:            "pr1",
				Commit:        "commit1",
				CommittedDate: time.Time{},
			},
			parameters:     resource.GetParameters{},
			pullRequest:    createTestPR(1, "master", false, false, 0, nil),
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"author_email","value":"user@example.com"},{"name":"title","value":"pr1 title"}]`,
		},
		{
			description: "get supports rebasing",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:            "pr1",
				Commit:        "commit1",
				CommittedDate: time.Time{},
			},
			parameters: resource.GetParameters{
				IntegrationTool: "rebase",
			},
			pullRequest:    createTestPR(1, "master", false, false, 0, nil),
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"author_email","value":"user@example.com"},{"name":"title","value":"pr1 title"}]`,
		},
		{
			description: "get supports checkout",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:            "pr1",
				Commit:        "commit1",
				CommittedDate: time.Time{},
			},
			parameters: resource.GetParameters{
				IntegrationTool: "checkout",
			},
			pullRequest:    createTestPR(1, "master", false, false, 0, nil),
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"author_email","value":"user@example.com"},{"name":"title","value":"pr1 title"}]`,
		},
		{
			description: "get supports git_depth",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:            "pr1",
				Commit:        "commit1",
				CommittedDate: time.Time{},
			},
			parameters: resource.GetParameters{
				GitDepth: 2,
			},
			pullRequest:    createTestPR(1, "master", false, false, 0, nil),
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"author_email","value":"user@example.com"},{"name":"title","value":"pr1 title"}]`,
		},
		{
			description: "get supports list_changed_files",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:            "pr1",
				Commit:        "commit1",
				CommittedDate: time.Time{},
			},
			parameters: resource.GetParameters{
				ListChangedFiles: true,
			},
			pullRequest: createTestPR(1, "master", false, false, 0, nil),
			files: []resource.ChangedFileObject{
				{
					Path: "README.md",
				},
				{
					Path: "Other.md",
				},
			},
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"author_email","value":"user@example.com"},{"name":"title","value":"pr1 title"}]`,
			filesString:    "README.md\nOther.md\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			github := new(fakes.FakeGithub)
			github.GetPullRequestReturns(tc.pullRequest, nil)

			if tc.files != nil {
				github.GetChangedFilesReturns(tc.files, nil)
			}

			git := new(fakes.FakeGit)
			git.RevParseReturns("sha", nil)

			dir := createTestDirectory(t)
			defer os.RemoveAll(dir)

			input := resource.GetRequest{Source: tc.source, Version: tc.version, Params: tc.parameters}
			output, err := resource.Get(input, github, git, dir)

			// Validate output
			if assert.NoError(t, err) {
				assert.Equal(t, tc.version, output.Version)

				// Verify written files
				version := readTestFile(t, filepath.Join(dir, ".git", "resource", "version.json"))
				assert.Equal(t, tc.versionString, version)

				metadata := readTestFile(t, filepath.Join(dir, ".git", "resource", "metadata.json"))
				assert.Equal(t, tc.metadataString, metadata)

				// Verify individual files
				files := map[string]string{
					"pr":           "1",
					"url":          "pr1 url",
					"head_name":    "pr1",
					"head_sha":     "oid1",
					"base_name":    "master",
					"base_sha":     "sha",
					"message":      "commit message1",
					"author":       "login1",
					"author_email": "user@example.com",
					"title":        "pr1 title",
				}

				for filename, expected := range files {
					actual := readTestFile(t, filepath.Join(dir, ".git", "resource", filename))
					assert.Equal(t, expected, actual)
				}

				if tc.files != nil {
					changedFiles := readTestFile(t, filepath.Join(dir, ".git", "resource", "changed_files"))
					assert.Equal(t, tc.filesString, changedFiles)
				}
			}

			// Validate Github calls
			if assert.Equal(t, 1, github.GetPullRequestCallCount()) {
				pr, commit := github.GetPullRequestArgsForCall(0)
				assert.Equal(t, tc.version.PR, pr)
				assert.Equal(t, tc.version.Commit, commit)
			}

			// Validate Git calls
			if assert.Equal(t, 1, git.InitCallCount()) {
				base := git.InitArgsForCall(0)
				assert.Equal(t, tc.pullRequest.BaseRefName, base)
			}

			if assert.Equal(t, 1, git.PullCallCount()) {
				url, base, depth := git.PullArgsForCall(0)
				assert.Equal(t, tc.pullRequest.Repository.URL, url)
				assert.Equal(t, tc.pullRequest.BaseRefName, base)
				assert.Equal(t, tc.parameters.GitDepth, depth)
			}

			if assert.Equal(t, 1, git.RevParseCallCount()) {
				base := git.RevParseArgsForCall(0)
				assert.Equal(t, tc.pullRequest.BaseRefName, base)
			}

			if assert.Equal(t, 1, git.FetchCallCount()) {
				url, pr, depth := git.FetchArgsForCall(0)
				assert.Equal(t, tc.pullRequest.Repository.URL, url)
				assert.Equal(t, tc.pullRequest.Number, pr)
				assert.Equal(t, tc.parameters.GitDepth, depth)
			}

			switch tc.parameters.IntegrationTool {
			case "rebase":
				if assert.Equal(t, 1, git.RebaseCallCount()) {
					branch, tip := git.RebaseArgsForCall(0)
					assert.Equal(t, tc.pullRequest.BaseRefName, branch)
					assert.Equal(t, tc.pullRequest.Tip.OID, tip)
				}
			case "checkout":
				if assert.Equal(t, 1, git.CheckoutCallCount()) {
					branch, sha := git.CheckoutArgsForCall(0)
					assert.Equal(t, tc.pullRequest.HeadRefName, branch)
					assert.Equal(t, tc.pullRequest.Tip.OID, sha)
				}
			default:
				if assert.Equal(t, 1, git.MergeCallCount()) {
					tip := git.MergeArgsForCall(0)
					assert.Equal(t, tc.pullRequest.Tip.OID, tip)
				}
			}
			if tc.source.GitCryptKey != "" {
				if assert.Equal(t, 1, git.GitCryptUnlockCallCount()) {
					key := git.GitCryptUnlockArgsForCall(0)
					assert.Equal(t, tc.source.GitCryptKey, key)
				}
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
			github := new(fakes.FakeGithub)
			git := new(fakes.FakeGit)
			dir := createTestDirectory(t)
			defer os.RemoveAll(dir)

			// Run the get and check output
			input := resource.GetRequest{Source: tc.source, Version: tc.version, Params: tc.parameters}
			output, err := resource.Get(input, github, git, dir)

			if assert.NoError(t, err) {
				assert.Equal(t, tc.version, output.Version)
			}
		})
	}
}

func createTestPR(
	count int,
	baseName string,
	skipCI bool,
	isCrossRepo bool,
	approvedReviews int,
	labels []string,
) *resource.PullRequest {
	n := strconv.Itoa(count)
	d := time.Now().AddDate(0, 0, -count)
	m := fmt.Sprintf("commit message%s", n)
	if skipCI {
		m = "[skip ci]" + m
	}
	approvedCount := approvedReviews

	var labelObjects []resource.LabelObject
	for _, l := range labels {
		lObject := resource.LabelObject{
			Name: l,
		}

		labelObjects = append(labelObjects, lObject)
	}

	return &resource.PullRequest{
		PullRequestObject: resource.PullRequestObject{
			ID:          fmt.Sprintf("pr%s", n),
			Number:      count,
			Title:       fmt.Sprintf("pr%s title", n),
			URL:         fmt.Sprintf("pr%s url", n),
			BaseRefName: baseName,
			HeadRefName: fmt.Sprintf("pr%s", n),
			Repository: struct{ URL string }{
				URL: fmt.Sprintf("repo%s url", n),
			},
			IsCrossRepository: isCrossRepo,
		},
		Tip: resource.CommitObject{
			ID:            fmt.Sprintf("commit%s", n),
			OID:           fmt.Sprintf("oid%s", n),
			CommittedDate: githubv4.DateTime{Time: d},
			Message:       m,
			Author: struct {
				User  struct{ Login string }
				Email string
			}{
				User: struct{ Login string }{
					Login: fmt.Sprintf("login%s", n),
				},
				Email: "user@example.com",
			},
		},
		ApprovedReviewCount: approvedCount,
		Labels:              labelObjects,
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
