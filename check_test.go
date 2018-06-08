package resource_test

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/itsdalmo/github-pr-resource"
	"github.com/itsdalmo/github-pr-resource/mocks"
)

var (
	testPullRequests = []*resource.PullRequest{
		{
			PullRequestObject: createTestPR(1),
			Tip:               createTestCommit(1, true),
		},
		{
			PullRequestObject: createTestPR(2),
			Tip:               createTestCommit(2, false),
		},
		{
			PullRequestObject: createTestPR(3),
			Tip:               createTestCommit(3, false),
		},
		{
			PullRequestObject: createTestPR(4),
			Tip:               createTestCommit(4, false),
		},
	}
	testPreviousPullRequests = []*resource.PullRequest{
		{
			PullRequestObject: createTestPR(1),
			Tip:               createTestCommit(3, true),
		},
		{
			PullRequestObject: createTestPR(2),
			Tip:               createTestCommit(4, false),
		},
		{
			PullRequestObject: createTestPR(3),
			Tip:               createTestCommit(2, false),
		},
		{
			PullRequestObject: createTestPR(4),
			Tip:               createTestCommit(6, false),
		},
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
				resource.NewVersion(testPullRequests[1], resource.GenerateVersion(testPullRequests[1:])),
			},
		},

		{
			description: "check returns the previous version when its still latest",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version:      resource.NewVersion(testPreviousPullRequests[1], resource.GenerateVersion(testPreviousPullRequests)),
			pullRequests: testPreviousPullRequests,
			files:        [][]string{},
			expected: resource.CheckResponse{
				resource.NewVersion(testPreviousPullRequests[1], resource.GenerateVersion(testPreviousPullRequests)),
			},
		},

		{
			description: "check returns all new versions since the last",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version:      resource.NewVersion(testPreviousPullRequests[3], resource.GenerateVersion(testPreviousPullRequests)),
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[3], resource.GenerateVersion(testPullRequests[1:4])),
				resource.NewVersion(testPullRequests[1], resource.GenerateVersion(testPullRequests[1:4])),
			},
		},

		{
			description: "check will only return versions that match the specified paths and are newer",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
				Paths:       []string{"terraform/*/*.tf", "terraform/*/*/*.tf"},
			},
			version:      resource.NewVersion(testPreviousPullRequests[3], resource.GenerateVersion(testPreviousPullRequests)),
			pullRequests: testPullRequests,
			files: [][]string{
				{"README.md", "travis.yml"},
				{"terraform/modules/ecs/main.tf", "README.md"},
				{"terraform/modules/variables.tf", "travis.yml"},
			},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[3], resource.GenerateVersion(testPullRequests[2:])),
			},
		},

		{
			description: "check will skip versions which only match the ignore paths",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
				IgnorePaths: []string{"*.md", "*.yml"},
			},
			version:      resource.NewVersion(testPullRequests[3], resource.GenerateVersion(testPullRequests[3:])),
			pullRequests: testPullRequests,
			files: [][]string{
				{"README.md", "travis.yml"},
				{"terraform/modules/ecs/main.tf", "README.md"},
				{"terraform/modules/variables.tf", "travis.yml"},
			},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[2], resource.GenerateVersion(testPullRequests[2:])),
			},
		},
		{
			description: "check correctly ignores [skip ci] when specified",
			source: resource.Source{
				Repository:    "itsdalmo/test-repository",
				AccessToken:   "oauthtoken",
				DisableCISkip: "true",
			},
			version:      resource.NewVersion(testPullRequests[1], resource.GenerateVersion(testPullRequests[1:])),
			pullRequests: testPullRequests,
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[0], resource.GenerateVersion(testPullRequests)),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			github := mocks.NewMockGithub(ctrl)
			github.EXPECT().ListOpenPullRequests().Times(1).Return(tc.pullRequests, nil)

			if len(tc.files) > 0 {
				// TODO: Figure out how to do this in a loop with variables. As is, it will break when adding new tests.
				gomock.InOrder(
					github.EXPECT().ListModifiedFiles(gomock.Any()).Times(1).Return(tc.files[0], nil),
					github.EXPECT().ListModifiedFiles(gomock.Any()).Times(1).Return(tc.files[1], nil),
					github.EXPECT().ListModifiedFiles(gomock.Any()).Times(1).Return(tc.files[2], nil),
				)
			}

			input := resource.CheckRequest{Source: tc.source, Version: tc.version}
			output, err := resource.Check(input, github)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if got, want := output, tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
			}
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
			if got != tc.want {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, tc.want)
			}
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
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := resource.FilterPath(tc.files, tc.pattern)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%s\n", got, tc.want)
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
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := resource.FilterIgnorePath(tc.files, tc.pattern)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%s\n", got, tc.want)
			}
		})
	}
}
