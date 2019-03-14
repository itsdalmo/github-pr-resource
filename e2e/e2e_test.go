// +build e2e

package e2e_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	resource "github.com/telia-oss/github-pr-resource"
)

var (
	targetCommitID       = "a5114f6ab89f4b736655642a11e8d15ce363d882"
	targetPullRequestID  = "4"
	targetDateTime       = time.Date(2018, time.May, 11, 8, 43, 48, 0, time.UTC)
	latestCommitID       = "890a7e4f0d5b05bda8ea21b91f4604e3e0313581"
	latestPullRequestID  = "5"
	latestDateTime       = time.Date(2018, time.May, 14, 10, 51, 58, 0, time.UTC)
	developCommitID      = "ac771f3b69cbd63b22bbda553f827ab36150c640"
	developPullRequestID = "6"
	developDateTime      = time.Date(2018, time.May, 14, 10, 51, 58, 0, time.UTC)
)

func TestCheckE2E(t *testing.T) {

	tests := []struct {
		description string
		source      resource.Source
		version     resource.Version
		expected    resource.CheckResponse
	}{
		{
			description: "check returns the latest version if there is no previous",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN"),
			},
			version: resource.Version{},
			expected: resource.CheckResponse{
				resource.Version{PR: latestPullRequestID, Commit: latestCommitID, CommittedDate: latestDateTime},
			},
		},

		{
			description: "check returns the previous version when its still latest",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN"),
			},
			version: resource.Version{PR: latestPullRequestID, Commit: latestCommitID, CommittedDate: latestDateTime},
			expected: resource.CheckResponse{
				resource.Version{PR: latestPullRequestID, Commit: latestCommitID, CommittedDate: latestDateTime},
			},
		},

		{
			description: "check returns all new versions since the last",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN"),
			},
			version: resource.Version{PR: targetPullRequestID, Commit: targetCommitID, CommittedDate: targetDateTime},
			expected: resource.CheckResponse{
				resource.Version{PR: latestPullRequestID, Commit: latestCommitID, CommittedDate: latestDateTime},
			},
		},

		{
			description: "check will only return versions that match the specified paths",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN"),
				Paths:       []string{"*.md"},
			},
			version: resource.Version{},
			expected: resource.CheckResponse{
				resource.Version{PR: targetPullRequestID, Commit: targetCommitID, CommittedDate: targetDateTime},
			},
		},

		{
			description: "check will skip versions which only match the ignore paths",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN"),
				IgnorePaths: []string{"*.txt"},
			},
			version: resource.Version{},
			expected: resource.CheckResponse{
				resource.Version{PR: targetPullRequestID, Commit: targetCommitID, CommittedDate: targetDateTime},
			},
		},

		{
			description: "check works with custom endpoints",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN"),
				V3Endpoint:  "https://api.github.com/",
				V4Endpoint:  "https://api.github.com/graphql",
			},
			version: resource.Version{},
			expected: resource.CheckResponse{
				resource.Version{PR: latestPullRequestID, Commit: latestCommitID, CommittedDate: latestDateTime},
			},
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

func TestGetAndPutE2E(t *testing.T) {
	tests := []struct {
		description    string
		source         resource.Source
		version        resource.Version
		getParameters  resource.GetParameters
		putParameters  resource.PutParameters
		versionString  string
		metadataString string
	}{
		{
			description: "get and put works",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				V3Endpoint:  "https://api.github.com/",
				V4Endpoint:  "https://api.github.com/graphql",
				AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN"),
			},
			version: resource.Version{
				PR:            targetPullRequestID,
				Commit:        targetCommitID,
				CommittedDate: time.Time{},
			},
			getParameters:  resource.GetParameters{},
			putParameters:  resource.PutParameters{},
			versionString:  `{"pr":"4","commit":"a5114f6ab89f4b736655642a11e8d15ce363d882","committed":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"4"},{"name":"url","value":"https://github.com/itsdalmo/test-repository/pull/4"},{"name":"head_name","value":"my_second_pull"},{"name":"head_sha","value":"a5114f6ab89f4b736655642a11e8d15ce363d882"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"93eeeedb8a16e6662062d1eca5655108977cc59a"},{"name":"message","value":"Push 2."},{"name":"author","value":"itsdalmo"}]`,
		},
		{
			description: "get works with non-master bases",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				V3Endpoint:  "https://api.github.com/",
				V4Endpoint:  "https://api.github.com/graphql",
				AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN"),
			},
			version: resource.Version{
				PR:            developPullRequestID,
				Commit:        developCommitID,
				CommittedDate: time.Time{},
			},
			getParameters:  resource.GetParameters{},
			putParameters:  resource.PutParameters{},
			versionString:  `{"pr":"6","commit":"ac771f3b69cbd63b22bbda553f827ab36150c640","committed":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"6"},{"name":"url","value":"https://github.com/itsdalmo/test-repository/pull/6"},{"name":"head_name","value":"test-develop-pr"},{"name":"head_sha","value":"ac771f3b69cbd63b22bbda553f827ab36150c640"},{"name":"base_name","value":"develop"},{"name":"base_sha","value":"93eeeedb8a16e6662062d1eca5655108977cc59a"},{"name":"message","value":"[skip ci] Add a PR with a non-master base"},{"name":"author","value":"itsdalmo"}]`,
		},
		{
			description: "get works when ssl verification is disabled",
			source: resource.Source{
				Repository:          "itsdalmo/test-repository",
				V3Endpoint:          "https://api.github.com/",
				V4Endpoint:          "https://api.github.com/graphql",
				AccessToken:         os.Getenv("GITHUB_ACCESS_TOKEN"),
				SkipSSLVerification: true,
			},
			version: resource.Version{
				PR:            targetPullRequestID,
				Commit:        targetCommitID,
				CommittedDate: time.Time{},
			},
			getParameters:  resource.GetParameters{},
			putParameters:  resource.PutParameters{},
			versionString:  `{"pr":"4","commit":"a5114f6ab89f4b736655642a11e8d15ce363d882","committed":"0001-01-01T00:00:00Z"}`,
			metadataString: `[{"name":"pr","value":"4"},{"name":"url","value":"https://github.com/itsdalmo/test-repository/pull/4"},{"name":"head_name","value":"my_second_pull"},{"name":"head_sha","value":"a5114f6ab89f4b736655642a11e8d15ce363d882"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"93eeeedb8a16e6662062d1eca5655108977cc59a"},{"name":"message","value":"Push 2."},{"name":"author","value":"itsdalmo"}]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			// Create temporary directory
			dir, err := ioutil.TempDir("", "github-pr-resource")
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			github, err := resource.NewGithubClient(&tc.source)
			require.NoError(t, err)

			git, err := resource.NewGitClient(&tc.source, dir, ioutil.Discard)
			require.NoError(t, err)

			// Get (output and files)
			getRequest := resource.GetRequest{Source: tc.source, Version: tc.version, Params: tc.getParameters}
			getOutput, err := resource.Get(getRequest, github, git, dir)

			require.NoError(t, err)
			assert.Equal(t, tc.version, getOutput.Version)

			version := readTestFile(t, filepath.Join(dir, ".git", "resource", "version.json"))
			assert.Equal(t, tc.versionString, version)

			metadata := readTestFile(t, filepath.Join(dir, ".git", "resource", "metadata.json"))
			assert.Equal(t, tc.metadataString, metadata)

			// Put
			putRequest := resource.PutRequest{Source: tc.source, Params: tc.putParameters}
			putOutput, err := resource.Put(putRequest, github, dir)

			require.NoError(t, err)
			assert.Equal(t, tc.version, putOutput.Version)
		})
	}
}

func readTestFile(t *testing.T, path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read: %s: %s", path, err)
	}
	return string(b)
}
