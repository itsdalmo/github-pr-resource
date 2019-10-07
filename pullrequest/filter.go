package pullrequest

import (
	"log"
	"regexp"
	"time"

	glob "github.com/sabhiram/go-gitignore"
)

// TimelineItem constants
const (
	BaseRefChangedEvent     = "BaseRefChangedEvent"
	BaseRefForcePushedEvent = "BaseRefForcePushedEvent"
	HeadRefForcePushedEvent = "HeadRefForcePushedEvent"
	IssueComment            = "IssueComment"
	PullRequestCommit       = "PullRequestCommit"
	ReopenedEvent           = "ReopenedEvent"
)

// Filter is a function that filters a slice of PRs, returning the filtered slice.
type Filter func(PullRequest) bool

// SkipCI returns true if the PR title or HeadRef message contains [skip ci] & the feature is not disabled
func SkipCI(disabled bool) Filter {
	return func(p PullRequest) bool {
		if disabled {
			return false
		}

		re := regexp.MustCompile("(?i)\\[(ci skip|skip ci)\\]")
		if re.MatchString(p.Title) || re.MatchString(p.HeadRef.Message) {
			log.Println("skipCI: true")
			return true
		}

		return false
	}
}

// Fork returns true if the source DisableForks is true && the PR is from a fork
func Fork(disabled bool) Filter {
	return func(p PullRequest) bool {
		if disabled && p.IsCrossRepository {
			log.Println("fork: true")
			return true
		}

		return false
	}
}

// BaseBranch returns true if the source BaseBranch is set & it does not match the PR
func BaseBranch(b string) Filter {
	return func(p PullRequest) bool {
		if b == "" {
			return false
		}

		if b != p.BaseRefName {
			log.Println("baseBranch: true")
			return true
		}

		return false
	}
}

// ApprovedReviewCount returns true if pr review count is lt than configured count
func ApprovedReviewCount(v int) Filter {
	return func(p PullRequest) bool {
		if p.ApprovedReviewCount < v {
			log.Println("review_count: true - ", p.ApprovedReviewCount, v)
			return true
		}
		log.Println("review_count: false - ", p.ApprovedReviewCount, v)
		return false
	}
}

// Labels returns true if pr does not have a configured label
func Labels(v []string) Filter {
	return func(p PullRequest) bool {
		if len(v) == 0 {
			log.Println("labels: false")
			return false
		}

		for _, i := range v {
			for _, k := range p.Labels {
				if i == k {
					log.Println("labels: false")
					return false
				}
			}
		}

		log.Println("labels: true")
		return true
	}
}

// Created returns true if the PR was created with no new commits or since the last check
func Created(v time.Time) Filter {
	return func(p PullRequest) bool {
		if p.CreatedAt.Equal(p.UpdatedAt) {
			log.Println("created: true")
			return true
		}
		if p.CreatedAt.After(latest(v, p.HeadRef.AuthoredDate, p.HeadRef.CommittedDate, p.HeadRef.PushedDate)) {
			log.Println("created: true")
			return true
		}
		return false
	}
}

// BuildCI returns true if a comment containing [build ci] was added since the last check
func BuildCI() Filter {
	return func(p PullRequest) bool {
		for _, c := range p.Comments {
			re := regexp.MustCompile("(?i)\\[(ci build|build ci)\\]")
			if re.MatchString(c.Body) {
				log.Println("buildCI: true")
				return true
			}
		}
		return false
	}
}

// NewCommits returns true if the PR has new commits since the input version.UpdatedDate
func NewCommits(v time.Time) Filter {
	return func(p PullRequest) bool {
		if v.IsZero() {
			log.Println("new commits: true")
			return true
		}

		if latest(p.HeadRef.AuthoredDate, p.HeadRef.CommittedDate, p.HeadRef.PushedDate).After(v) {
			log.Println("new commits: true")
			return true
		}
		return false
	}
}

// BaseRefChanged returns true if the PR contains a BaseRefChangedEvent since the last check
func BaseRefChanged() Filter {
	return filterEvent(BaseRefChangedEvent)
}

// BaseRefForcePushed returns true if the PR contains a BaseRefForcePushedEvent since the last check
func BaseRefForcePushed() Filter {
	return filterEvent(BaseRefForcePushedEvent)
}

// HeadRefForcePushed returns true if the PR contains a HeadRefForcePushedEvent since the last check
func HeadRefForcePushed() Filter {
	return filterEvent(HeadRefForcePushedEvent)
}

// Reopened returns true if the PR contains a ReopenedEvent since the last check
func Reopened() Filter {
	return filterEvent(ReopenedEvent)
}

func filterEvent(eventType string) Filter {
	return func(p PullRequest) bool {
		for _, i := range p.Events {
			log.Println("filter:", eventType, "item type:", i.Type)
			if eventType != i.Type {
				continue
			}

			log.Println(eventType, ": true")
			return true
		}
		return false
	}
}

// Patterns returns true if there is a pattern configured
func Patterns(patterns []string) Filter {
	return func(p PullRequest) bool {
		if len(patterns) > 0 {
			return true
		}

		return false
	}
}

// Files matches the PRs changed files against a set of glob patterns:
// (invert == false) returns true if patterns empty OR any of the changed files match the patterns
// (invert == true) returns true if all of the changed files match the patterns to ignore
func Files(patterns []string, invert bool) Filter {
	return func(p PullRequest) bool {
		matched := make([]int8, 0)
		gc, err := glob.CompileIgnoreLines(patterns...)
		if err != nil {
			return false
		}

		for _, f := range p.Files {
			log.Println("comparing patterns to changed file:", f)
			if gc.MatchesPath("/" + f) {
				if !invert {
					log.Println("paths: true -", "matched:", f)
					return true
				}

				matched = append(matched, 1)
				log.Println("matched:", f)
			}
		}

		log.Println("invert:", invert, "- matched:", len(matched), "- files:", len(p.Files))
		if invert && len(matched) == len(p.Files) {
			log.Println("ignore paths: true")
			return true
		}

		return false
	}
}

func latest(times ...time.Time) time.Time {
	var latest time.Time
	for _, t := range times {
		if t.After(latest) {
			latest = t
		}
	}
	return latest
}
