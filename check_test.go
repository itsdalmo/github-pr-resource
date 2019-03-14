package resource_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	resource "github.com/telia-oss/github-pr-resource"
	"github.com/telia-oss/github-pr-resource/fakes"
)

var (
	testPullRequests = []*resource.PullRequest{
		createTestPR(1, true, false),
		createTestPR(2, false, false),
		createTestPR(3, false, false),
		createTestPR(4, false, false),
		createTestPR(5, false, true),
		createTestPR(6, false, false),
	}
)

func TestCheck(t *testing.T) {
	tests := []struct {
		description  string
		source       resource.Source
		version      resource.Version
		files        [][]string
		pullRequests []*resource.PullRequest
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
				resource.NewVersion(testPullRequests[1]),
			},
		},

		{
			description: "check returns the previous version when its still latest",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version:      resource.NewVersion(testPullRequests[1]),
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[1]),
			},
		},

		{
			description: "check returns all new versions since the last",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version:      resource.NewVersion(testPullRequests[3]),
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[2]),
				resource.NewVersion(testPullRequests[1]),
			},
		},

		{
			description: "check will only return versions that match the specified paths",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
				Paths:       []string{"terraform/*/*.tf", "terraform/*/*/*.tf"},
			},
			version:      resource.NewVersion(testPullRequests[3]),
			pullRequests: testPullRequests,
			files: [][]string{
				{"README.md", "travis.yml"},
				{"terraform/modules/ecs/main.tf", "README.md"},
				{"terraform/modules/variables.tf", "travis.yml"},
			},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[2]),
			},
		},

		{
			description: "check will skip versions which only match the ignore paths",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
				IgnorePaths: []string{"*.md", "*.yml"},
			},
			version:      resource.NewVersion(testPullRequests[3]),
			pullRequests: testPullRequests,
			files: [][]string{
				{"README.md", "travis.yml"},
				{"terraform/modules/ecs/main.tf", "README.md"},
				{"terraform/modules/variables.tf", "travis.yml"},
			},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[2]),
			},
		},
		{
			description: "check correctly ignores [skip ci] when specified",
			source: resource.Source{
				Repository:    "itsdalmo/test-repository",
				AccessToken:   "oauthtoken",
				DisableCISkip: true,
			},
			version:      resource.NewVersion(testPullRequests[1]),
			pullRequests: testPullRequests,
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[0]),
			},
		},
		{
			description: "check correctly ignores cross repo pull requests",
			source: resource.Source{
				Repository:   "itsdalmo/test-repository",
				AccessToken:  "oauthtoken",
				DisableForks: true,
			},
			version:      resource.NewVersion(testPullRequests[5]),
			pullRequests: testPullRequests,
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[3]),
				resource.NewVersion(testPullRequests[2]),
				resource.NewVersion(testPullRequests[1]),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			github := new(fakes.FakeGithub)
			github.ListOpenPullRequestsReturns(tc.pullRequests, nil)

			for i, file := range tc.files {
				github.ListModifiedFilesReturnsOnCall(i, file, nil)
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

func TestContainsSkipCI(t *testing.T) {
	tests := []struct {
		description string
		message     string
		want        bool
	}{
		{
			description: "does not just match any symbol in the regexp",
			message:     "(",
			want:        false,
		},
		{
			description: "does not match when it should not",
			message:     "test",
			want:        false,
		},
		{
			description: "matches [ci skip]",
			message:     "[ci skip]",
			want:        true,
		},
		{
			description: "matches [skip ci]",
			message:     "[skip ci]",
			want:        true,
		},
		{
			description: "matches trailing skip ci",
			message:     "trailing [skip ci]",
			want:        true,
		},
		{
			description: "matches leading skip ci",
			message:     "[skip ci] leading",
			want:        true,
		},
		{
			description: "is case insensitive",
			message:     "case[Skip CI]insensitive",
			want:        true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			got := resource.ContainsSkipCI(tc.message)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFilterPath(t *testing.T) {
	cases := []struct {
		description string
		pattern     string
		files       []string
		want        []string
	}{
		{
			description: "returns all matching files",
			pattern:     "*.txt",
			files: []string{
				"file1.txt",
				"test/file2.txt",
			},
			want: []string{
				"file1.txt",
			},
		},
		{
			description: "works with wildcard",
			pattern:     "test/*",
			files: []string{
				"file1.txt",
				"test/file2.txt",
			},
			want: []string{
				"test/file2.txt",
			},
		},
		{
			description: "excludes unmatched files",
			pattern:     "*/*.txt",
			files: []string{
				"test/file1.go",
				"test/file2.txt",
			},
			want: []string{
				"test/file2.txt",
			},
		},
		{
			description: "handles prefix matches",
			pattern:     "foo/",
			files: []string{
				"foo/a",
				"foo/a.txt",
				"foo/a/b/c/d.txt",
				"foo",
				"bar",
				"bar/a.txt",
			},
			want: []string{
				"foo/a",
				"foo/a.txt",
				"foo/a/b/c/d.txt",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := resource.FilterPath(tc.files, tc.pattern)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestFilterIgnorePath(t *testing.T) {
	cases := []struct {
		description string
		pattern     string
		files       []string
		want        []string
	}{
		{
			description: "excludes all matching files",
			pattern:     "*.txt",
			files: []string{
				"file1.txt",
				"test/file2.txt",
			},
			want: []string{
				"test/file2.txt",
			},
		},
		{
			description: "works with wildcard",
			pattern:     "test/*",
			files: []string{
				"file1.txt",
				"test/file2.txt",
			},
			want: []string{
				"file1.txt",
			},
		},
		{
			description: "includes unmatched files",
			pattern:     "*/*.txt",
			files: []string{
				"test/file1.go",
				"test/file2.txt",
			},
			want: []string{
				"test/file1.go",
			},
		},
		{
			description: "handles prefix matches",
			pattern:     "foo/",
			files: []string{
				"foo/a",
				"foo/a.txt",
				"foo/a/b/c/d.txt",
				"foo",
				"bar",
				"bar/a.txt",
			},
			want: []string{
				"foo",
				"bar",
				"bar/a.txt",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := resource.FilterIgnorePath(tc.files, tc.pattern)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestIsInsidePath(t *testing.T) {
	cases := []struct {
		description string
		parent      string

		expectChildren    []string
		expectNotChildren []string

		want bool
	}{
		{
			description: "basic test",
			parent:      "foo/bar",
			expectChildren: []string{
				"foo/bar",
				"foo/bar/baz",
			},
			expectNotChildren: []string{
				"foo/barbar",
				"foo/baz/bar",
			},
		},
		{
			description: "does not match parent directories against child files",
			parent:      "foo/",
			expectChildren: []string{
				"foo/bar",
			},
			expectNotChildren: []string{
				"foo",
			},
		},
		{
			description: "matches parents without trailing slash",
			parent:      "foo/bar",
			expectChildren: []string{
				"foo/bar",
				"foo/bar/baz",
			},
		},
		{
			description: "handles children that are shorter than the parent",
			parent:      "foo/bar/baz",
			expectNotChildren: []string{
				"foo",
				"foo/bar",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			for _, expectedChild := range tc.expectChildren {
				if !resource.IsInsidePath(tc.parent, expectedChild) {
					t.Errorf("Expected \"%s\" to be inside \"%s\"", expectedChild, tc.parent)
				}
			}

			for _, expectedNotChild := range tc.expectNotChildren {
				if resource.IsInsidePath(tc.parent, expectedNotChild) {
					t.Errorf("Expected \"%s\" to not be inside \"%s\"", expectedNotChild, tc.parent)
				}
			}
		})
	}
}
