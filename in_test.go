package resource_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	resource "github.com/telia-oss/github-pr-resource"
	"github.com/telia-oss/github-pr-resource/fakes"
	"github.com/telia-oss/github-pr-resource/pullrequest"
)

func TestGet(t *testing.T) {
	tests := []struct {
		description    string
		source         resource.Source
		version        resource.Version
		parameters     resource.GetParameters
		pullRequest    pullrequest.PullRequest
		versionString  string
		metadataString string
		files          []string
		filesString    string
	}{
		{
			description: "get works",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:          1,
				Commit:      "commit1",
				UpdatedDate: time.Time{},
			},
			parameters:     resource.GetParameters{},
			pullRequest:    createTestPR(1, "master", false, false, false, false, 0, nil),
			versionString:  `{"pr":"1","commit":"commit1","updated":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"head_short_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"events","value":"[]"}]`,
		},
		{
			description: "get supports unlocking with git crypt",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
				GitCryptKey: "gitcryptkey",
			},
			version: resource.Version{
				PR:          1,
				Commit:      "commit1",
				UpdatedDate: time.Time{},
			},
			parameters:     resource.GetParameters{},
			pullRequest:    createTestPR(1, "master", false, false, false, false, 0, nil),
			versionString:  `{"pr":"1","commit":"commit1","updated":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"head_short_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"events","value":"[]"}]`,
		},
		{
			description: "get supports rebasing",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:          1,
				Commit:      "commit1",
				UpdatedDate: time.Time{},
			},
			parameters: resource.GetParameters{
				IntegrationTool: "rebase",
			},
			pullRequest:    createTestPR(1, "master", false, false, false, false, 0, nil),
			versionString:  `{"pr":"1","commit":"commit1","updated":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"head_short_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"events","value":"[]"}]`,
		},
		{
			description: "get supports merge",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:          1,
				Commit:      "commit1",
				UpdatedDate: time.Time{},
			},
			parameters: resource.GetParameters{
				IntegrationTool: "merge",
			},
			pullRequest:    createTestPR(1, "master", false, false, false, false, 0, nil),
			versionString:  `{"pr":"1","commit":"commit1","updated":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"head_short_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"events","value":"[]"}]`,
		},
		{
			description: "get supports git_depth",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:          1,
				Commit:      "commit1",
				UpdatedDate: time.Time{},
			},
			parameters: resource.GetParameters{
				GitDepth: 2,
			},
			pullRequest:    createTestPR(1, "master", false, false, false, false, 0, nil),
			versionString:  `{"pr":"1","commit":"commit1","updated":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"head_short_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"events","value":"[]"}]`,
		},
		{
			description: "get supports list_changed_files",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version: resource.Version{
				PR:          1,
				Commit:      "commit1",
				UpdatedDate: time.Time{},
			},
			parameters: resource.GetParameters{
				ListChangedFiles: true,
			},
			pullRequest:    createTestPR(1, "master", false, false, false, false, 0, nil),
			files:          []string{"README.md", "Other.md"},
			versionString:  `{"pr":"1","commit":"commit1","updated":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"head_short_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"events","value":"[]"}]`,
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
					"pr":             "1",
					"url":            "pr1 url",
					"head_name":      "pr1",
					"head_sha":       "oid1",
					"head_short_sha": "oid1",
					"base_name":      "master",
					"base_sha":       "sha",
					"message":        "commit message1",
					"author":         "login1",
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
			if assert.Equal(t, 1, git.CloneCallCount()) {
				url, base, depth := git.CloneArgsForCall(0)
				assert.Equal(t, tc.pullRequest.RepositoryURL, url)
				assert.Equal(t, tc.pullRequest.BaseRefName, base)
				assert.Equal(t, tc.parameters.GitDepth, depth)
			}

			if assert.Equal(t, 1, git.FetchCallCount()) {
				pr, depth := git.FetchArgsForCall(0)
				assert.Equal(t, tc.pullRequest.Number, pr)
				assert.Equal(t, tc.parameters.GitDepth, depth)
			}

			switch tc.parameters.IntegrationTool {
			case "rebase":
				if assert.Equal(t, 1, git.RevParseCallCount()) {
					base := git.RevParseArgsForCall(0)
					assert.Equal(t, tc.pullRequest.BaseRefName, base)
				}

				if assert.Equal(t, 1, git.RebaseCallCount()) {
					branch, tip := git.RebaseArgsForCall(0)
					assert.Equal(t, tc.pullRequest.BaseRefName, branch)
					assert.Equal(t, tc.version.Commit, tip)
				}
			case "checkout", "":
				if assert.Equal(t, 1, git.CheckoutCallCount()) {
					branch, sha := git.CheckoutArgsForCall(0)
					assert.Equal(t, tc.pullRequest.HeadRefName, branch)
					assert.Equal(t, tc.version.Commit, sha)
				}
			case "merge":
				if assert.Equal(t, 1, git.RevParseCallCount()) {
					base := git.RevParseArgsForCall(0)
					assert.Equal(t, tc.pullRequest.BaseRefName, base)
				}

				if assert.Equal(t, 1, git.MergeCallCount()) {
					tip := git.MergeArgsForCall(0)
					assert.Equal(t, tc.version.Commit, tip)
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
				PR:          1,
				Commit:      "commit1",
				UpdatedDate: time.Time{},
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
