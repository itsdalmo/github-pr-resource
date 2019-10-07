package pullrequest

import "time"

// PullRequest represents a pull request
type PullRequest struct {
	ID                  string
	Number              int
	Title               string
	URL                 string
	RepositoryURL       string
	BaseRefName         string
	BaseRefOID          string
	HeadRefName         string
	IsCrossRepository   bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
	HeadRef             Commit
	Events              []Event
	Comments            []Comment
	Commits             []Commit
	Files               []string
	Labels              []string
	ApprovedReviewCount int
}

// Commit represents a commit
type Commit struct {
	OID            string
	AbbreviatedOID string
	AuthoredDate   time.Time
	CommittedDate  time.Time
	PushedDate     time.Time
	Message        string
	Author         string
}

// Event represents an event that has been recorded on the PR
type Event struct {
	Type      string
	CreatedAt time.Time
}

// Comment represents a comment on a PR
type Comment struct {
	CreatedAt time.Time
	Body      string
}
