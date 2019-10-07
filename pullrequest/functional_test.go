package pullrequest_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/telia-oss/github-pr-resource/pullrequest"
)

// Test Filters

func TestSkipCI(t *testing.T) {
	tests := []struct {
		description string
		disabled    bool
		pull        pullrequest.PullRequest
		expect      bool
	}{
		{
			description: "match title v1",
			disabled:    false,
			pull: pullrequest.PullRequest{
				Title: "WIP: Weeee... [skip ci]",
				HeadRef: pullrequest.Commit{
					Message: "this is a test",
				},
			},
			expect: true,
		},
		{
			description: "match title v2",
			disabled:    false,
			pull: pullrequest.PullRequest{
				Title: "WIP: Weeee... [ci skip]",
				HeadRef: pullrequest.Commit{
					Message: "this is a test",
				},
			},
			expect: true,
		},
		{
			description: "match message v1",
			disabled:    false,
			pull: pullrequest.PullRequest{
				Title: "WIP: Weeee...",
				HeadRef: pullrequest.Commit{
					Message: "this is a test [skip ci]",
				},
			},
			expect: true,
		},
		{
			description: "match message v2",
			disabled:    false,
			pull: pullrequest.PullRequest{
				Title: "WIP: Weeee...",
				HeadRef: pullrequest.Commit{
					Message: "this is a test [ci skip]",
				},
			},
			expect: true,
		},
		{
			description: "match nothing",
			disabled:    false,
			pull: pullrequest.PullRequest{
				Title: "WIP: Weeee... I should skip to the CI",
				HeadRef: pullrequest.Commit{
					Message: "this is a test ci [ skip]",
				},
			},
			expect: false,
		},
		{
			description: "match nothing disabled",
			disabled:    true,
			pull: pullrequest.PullRequest{
				Title: "WIP: Weeee... [ci skip]",
				HeadRef: pullrequest.Commit{
					Message: "this is a test [ci skip]",
				},
			},
			expect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.SkipCI(tc.disabled)(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestFork(t *testing.T) {
	tests := []struct {
		description string
		disabled    bool
		pull        pullrequest.PullRequest
		expect      bool
	}{
		{
			description: "match",
			disabled:    true,
			pull: pullrequest.PullRequest{
				IsCrossRepository: true,
			},
			expect: true,
		},
		{
			description: "no match",
			disabled:    false,
			pull: pullrequest.PullRequest{
				IsCrossRepository: false,
			},
			expect: false,
		},
		{
			description: "no match disabled",
			disabled:    true,
			pull: pullrequest.PullRequest{
				IsCrossRepository: false,
			},
			expect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.Fork(tc.disabled)(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestBaseBranch(t *testing.T) {
	tests := []struct {
		description string
		branch      string
		pull        pullrequest.PullRequest
		expect      bool
	}{
		{
			description: "match set",
			branch:      "develop",
			pull: pullrequest.PullRequest{
				BaseRefName: "master",
			},
			expect: true,
		},
		{
			description: "no match not set",
			branch:      "",
			pull: pullrequest.PullRequest{
				BaseRefName: "master",
			},
			expect: false,
		},
		{
			description: "no match set",
			branch:      "master",
			pull: pullrequest.PullRequest{
				BaseRefName: "master",
			},
			expect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.BaseBranch(tc.branch)(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestLabels(t *testing.T) {
	tests := []struct {
		description string
		labels      []string
		pull        pullrequest.PullRequest
		expect      bool
	}{
		{
			description: "match label",
			labels:      []string{"develop"},
			pull: pullrequest.PullRequest{
				Labels: []string{"develop", "wip", "stupid-ball"},
			},
			expect: false,
		},
		{
			description: "no match not set",
			labels:      []string{},
			pull: pullrequest.PullRequest{
				Labels: []string{"develop", "wip", "stupid-ball"},
			},
			expect: false,
		},
		{
			description: "no match set",
			labels:      []string{"stupid-ball"},
			pull: pullrequest.PullRequest{
				Labels: []string{"develop", "wip"},
			},
			expect: true,
		},
		{
			description: "no match set",
			labels:      []string{"stupid-ball"},
			pull: pullrequest.PullRequest{
				Labels: []string{},
			},
			expect: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.Labels(tc.labels)(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestRequiredApprovals(t *testing.T) {
	tests := []struct {
		description string
		approvals   int
		pull        pullrequest.PullRequest
		expect      bool
	}{
		{
			description: "0 is not less than 0 requirement",
			approvals:   0,
			pull: pullrequest.PullRequest{
				ApprovedReviewCount: 0,
			},
			expect: false,
		},
		{
			description: "1 is not less than 0 requirement",
			approvals:   0,
			pull: pullrequest.PullRequest{
				ApprovedReviewCount: 1,
			},
			expect: false,
		},
		{
			description: "0 is less than 1 requirement",
			approvals:   1,
			pull: pullrequest.PullRequest{
				ApprovedReviewCount: 0,
			},
			expect: true,
		},
		{
			description: "1 is not less than 1 requirement",
			approvals:   1,
			pull: pullrequest.PullRequest{
				ApprovedReviewCount: 1,
			},
			expect: false,
		},
		{
			description: "2 is not less than 1 requirement",
			approvals:   1,
			pull: pullrequest.PullRequest{
				ApprovedReviewCount: 2,
			},
			expect: false,
		},
		{
			description: "1 is less than 2 requirement",
			approvals:   2,
			pull: pullrequest.PullRequest{
				ApprovedReviewCount: 1,
			},
			expect: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.ApprovedReviewCount(tc.approvals)(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestCreated(t *testing.T) {
	tests := []struct {
		description string
		pull        pullrequest.PullRequest
		since       time.Time
		expect      bool
	}{
		{
			description: "match Created==Updated",
			pull: pullrequest.PullRequest{
				CreatedAt: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			since:  time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			expect: true,
		},
		{
			description: "match Created after head commit and since last check",
			pull: pullrequest.PullRequest{
				CreatedAt: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
				HeadRef: pullrequest.Commit{
					CommittedDate: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
					PushedDate:    time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC),
					AuthoredDate:  time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			since:  time.Date(2018, 12, 1, 0, 0, 0, 0, time.UTC),
			expect: true,
		},
		{
			description: "no match Created after head commit but not since last check",
			pull: pullrequest.PullRequest{
				CreatedAt: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
				HeadRef: pullrequest.Commit{
					CommittedDate: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
					PushedDate:    time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC),
					AuthoredDate:  time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			since:  time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC),
			expect: false,
		},
		{
			description: "no match",
			pull: pullrequest.PullRequest{
				CreatedAt: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
				HeadRef: pullrequest.Commit{
					CommittedDate: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
					PushedDate:    time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC),
					AuthoredDate:  time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			since:  time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			expect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.Created(tc.since)(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestNewCommits(t *testing.T) {
	tests := []struct {
		description string
		versionDate time.Time
		pull        pullrequest.PullRequest
		expect      bool
	}{
		{
			description: "match empty version",
			versionDate: time.Time{},
			pull: pullrequest.PullRequest{
				HeadRef: pullrequest.Commit{
					CommittedDate: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
					PushedDate:    time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC),
					AuthoredDate:  time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expect: true,
		},
		{
			description: "no match version date now",
			versionDate: time.Now(),
			pull: pullrequest.PullRequest{
				HeadRef: pullrequest.Commit{
					CommittedDate: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
					PushedDate:    time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC),
					AuthoredDate:  time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.NewCommits(tc.versionDate)(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestBuildCI(t *testing.T) {
	tests := []struct {
		description string
		pull        pullrequest.PullRequest
		expect      bool
	}{
		{
			description: "match v1",
			pull: pullrequest.PullRequest{
				Comments: []pullrequest.Comment{
					{Body: "weeee I don't want to build ci"},
					{Body: "weeee I do want to [build ci]"},
				},
			},
			expect: true,
		},
		{
			description: "match v2",
			pull: pullrequest.PullRequest{
				Comments: []pullrequest.Comment{
					{Body: "weeee I don't want to build ci"},
					{Body: "weeee I do want to [ci build]"},
				},
			},
			expect: true,
		},
		{
			description: "no match",
			pull: pullrequest.PullRequest{
				Comments: []pullrequest.Comment{
					{Body: "weeee I don't want to build ci"},
				},
			},
			expect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.BuildCI()(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestBaseRefChanged(t *testing.T) {
	tests := []struct {
		description string
		pull        pullrequest.PullRequest
		expect      bool
	}{
		{
			description: "match single",
			pull: pullrequest.PullRequest{
				Events: []pullrequest.Event{
					{
						Type:      pullrequest.BaseRefChangedEvent,
						CreatedAt: time.Now(),
					},
				},
			},
			expect: true,
		},
		{
			description: "match many",
			pull: pullrequest.PullRequest{
				Events: []pullrequest.Event{
					{
						Type:      pullrequest.BaseRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
					{
						Type:      pullrequest.BaseRefChangedEvent,
						CreatedAt: time.Now(),
					},
				},
			},
			expect: true,
		},
		{
			description: "no match",
			pull: pullrequest.PullRequest{
				Events: []pullrequest.Event{
					{
						Type:      pullrequest.BaseRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
				},
			},
			expect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.BaseRefChanged()(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestBaseRefForcePushed(t *testing.T) {
	tests := []struct {
		description string
		pull        pullrequest.PullRequest
		expect      bool
	}{
		{
			description: "match single",
			pull: pullrequest.PullRequest{
				Events: []pullrequest.Event{
					{
						Type:      pullrequest.BaseRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
				},
			},
			expect: true,
		},
		{
			description: "match many",
			pull: pullrequest.PullRequest{
				Events: []pullrequest.Event{
					{
						Type:      pullrequest.HeadRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
					{
						Type:      pullrequest.HeadRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
					{
						Type:      pullrequest.BaseRefChangedEvent,
						CreatedAt: time.Now(),
					},
					{
						Type:      pullrequest.BaseRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
				},
			},
			expect: true,
		},
		{
			description: "no match",
			pull: pullrequest.PullRequest{
				Events: []pullrequest.Event{
					{
						Type:      pullrequest.BaseRefChangedEvent,
						CreatedAt: time.Now(),
					},
				},
			},
			expect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.BaseRefForcePushed()(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestHeadRefForcePushed(t *testing.T) {
	tests := []struct {
		description string
		pull        pullrequest.PullRequest
		expect      bool
	}{
		{
			description: "match single",
			pull: pullrequest.PullRequest{
				Events: []pullrequest.Event{
					{
						Type:      pullrequest.HeadRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
				},
			},
			expect: true,
		},
		{
			description: "match many",
			pull: pullrequest.PullRequest{
				Events: []pullrequest.Event{
					{
						Type:      pullrequest.HeadRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
					{
						Type:      pullrequest.HeadRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
					{
						Type:      pullrequest.BaseRefChangedEvent,
						CreatedAt: time.Now(),
					},
					{
						Type:      pullrequest.BaseRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
				},
			},
			expect: true,
		},
		{
			description: "no match",
			pull: pullrequest.PullRequest{
				Events: []pullrequest.Event{
					{
						Type:      pullrequest.BaseRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
				},
			},
			expect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.HeadRefForcePushed()(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestReopened(t *testing.T) {
	tests := []struct {
		description string
		pull        pullrequest.PullRequest
		expect      bool
	}{
		{
			description: "match single",
			pull: pullrequest.PullRequest{
				Events: []pullrequest.Event{
					{
						Type:      pullrequest.ReopenedEvent,
						CreatedAt: time.Now(),
					},
				},
			},
			expect: true,
		},
		{
			description: "match many",
			pull: pullrequest.PullRequest{
				Events: []pullrequest.Event{
					{
						Type:      pullrequest.HeadRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
					{
						Type:      pullrequest.HeadRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
					{
						Type:      pullrequest.ReopenedEvent,
						CreatedAt: time.Now(),
					},
					{
						Type:      pullrequest.BaseRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
				},
			},
			expect: true,
		},
		{
			description: "no match",
			pull: pullrequest.PullRequest{
				Events: []pullrequest.Event{
					{
						Type:      pullrequest.HeadRefForcePushedEvent,
						CreatedAt: time.Now(),
					},
				},
			},
			expect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.Reopened()(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestFiles(t *testing.T) {
	tests := []struct {
		description string
		patterns    []string
		invert      bool
		pull        pullrequest.PullRequest
		expect      bool
	}{
		{
			description: "match txt files @ root level",
			patterns:    []string{"*.txt"},
			invert:      false,
			pull: pullrequest.PullRequest{
				Files: []string{
					"file1.txt",
				},
			},
			expect: true,
		},
		{
			description: "match txt files at any level",
			patterns:    []string{"**/*.txt"},
			invert:      false,
			pull: pullrequest.PullRequest{
				Files: []string{
					"test/file2.txt",
					"test/testing/file2.txt",
					"test/testing/tested/file2.txt",
				},
			},
			expect: true,
		},
		{
			description: "match any file in test dir",
			patterns:    []string{"test/*"},
			invert:      false,
			pull: pullrequest.PullRequest{
				Files: []string{
					"file1.txt",
					"test/file2.txt",
				},
			},
			expect: true,
		},
		{
			description: "match any file recursively in test dir",
			patterns:    []string{"test/**"},
			invert:      false,
			pull: pullrequest.PullRequest{
				Files: []string{
					"file1.txt",
					"test/testing/file2.txt",
					"test/testing/tested/file2.txt",
				},
			},
			expect: true,
		},
		{
			description: "match multiple files",
			patterns:    []string{"ci/dockerfiles/**/*", "ci/dockerfiles/*", "ci/tasks/build-image.yml", "ci/pipelines/images.yml"},
			invert:      false,
			pull: pullrequest.PullRequest{
				Files: []string{
					"ci/Makefile",
					"ci/pipelines/images.yml",
					"terraform/Makefile",
				},
			},
			expect: true,
		},
		{
			description: "no match /**/*",
			patterns:    []string{"/ci/**/*"},
			invert:      false,
			pull: pullrequest.PullRequest{
				Files: []string{
					"src/teams/myteam/README.md",
					"src/teams/myteam/ci/pipeline.yml",
				},
			},
			expect: false,
		},
		{
			description: "match /**/*",
			patterns:    []string{"/ci/**/*"},
			invert:      false,
			pull: pullrequest.PullRequest{
				Files: []string{
					"ci/README.md",
					"ci/pipeline.yml",
				},
			},
			expect: true,
		},
		{
			description: "match **/*",
			patterns:    []string{"src/teams/myteam/**/*"},
			invert:      false,
			pull: pullrequest.PullRequest{
				Files: []string{
					"src/teams/myteam/README.md",
					"src/teams/myteam/ci/pipeline.yml",
				},
			},
			expect: true,
		},
		{
			description: "match multiple files 2",
			patterns:    []string{"src/teams/myteam/**/*", "src/teams/myteam/*"},
			invert:      false,
			pull: pullrequest.PullRequest{
				Files: []string{
					"src/teams/myteam/README.md",
					"src/teams/myteam/ci/pipeline.yml",
				},
			},
			expect: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := pullrequest.Files(tc.patterns, tc.invert)(tc.pull)
			assert.Equal(t, tc.expect, out)
		})
	}
}
