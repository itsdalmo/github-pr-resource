package resource

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/telia-oss/github-pr-resource/pullrequest"
)

// Source represents the configuration for the resource.
type Source struct {
	// Repository to check, get, put
	Repository string `json:"repository"`
	// AccessToken for GitHub API with permissions to Repository
	AccessToken string `json:"access_token"`
	// V3Endpoint for GitHub Rest API (leave blank for cloud)
	V3Endpoint string `json:"v3_endpoint"`
	// V4Endpoint for GitHub GraphQL API (leave blank for cloud)
	V4Endpoint string `json:"v4_endpoint"`
	// Paths of Repository to return versions for
	Paths []string `json:"paths,omitempty"`
	// IgnorePaths of Repository to skip returning versions for
	IgnorePaths []string `json:"ignore_paths,omitempty"`
	// DisableCISkip disables ability to skip CI via PR title / message
	DisableCISkip bool `json:"disable_ci_skip,omitempty"`
	// SkipSSLVerification when executing GitHub API requests
	SkipSSLVerification bool `json:"skip_ssl_verification,omitempty"`
	// DisableForks disables versions from forks of Repository
	DisableForks bool `json:"disable_forks,omitempty"`
	// GitCryptKey enables GitCrypt unlocking
	GitCryptKey string `json:"git_crypt_key,omitempty"`
	// BaseBranch returns versions only for matching base branch
	BaseBranch string `json:"base_branch,omitempty"`
	// PreviewSchema enables GraphQL preview schemas, see: https://developer.github.com/v4/previews/
	PreviewSchema bool `json:"preview_schema,omitempty"`
	// RequiredReviewApprovals returns versions when PR has >= approvals
	RequiredReviewApprovals int `json:"required_review_approvals,omitempty"`
	// Labels returns versions for PRs matching labels
	Labels []string `json:"labels,omitempty"`
}

// Validate the source configuration.
func (s *Source) Validate() error {
	if s.AccessToken == "" || s.Repository == "" {
		return errors.New("access_token & repository are required")
	}

	if len(s.V3Endpoint)+len(s.V4Endpoint) > 0 && (s.V3Endpoint == "" || s.V4Endpoint == "") {
		return errors.New("both v3_endpoint & v4_endpoint endpoints are required for GitHub Enterprise")
	}

	return nil
}

// Metadata output from get/put steps.
type Metadata []*MetadataField

// Add a MetadataField to the Metadata.
func (m *Metadata) Add(name, value string) {
	*m = append(*m, &MetadataField{Name: name, Value: value})
}

// MetadataField ...
type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Version communicated with Concourse.
type Version struct {
	PR          int       `json:"pr"`
	Commit      string    `json:"commit"`
	UpdatedDate time.Time `json:"updated"`
}

// MarshalJSON custom marshaller to convert PR number
func (v *Version) MarshalJSON() ([]byte, error) {
	type Alias Version
	return json.Marshal(&struct {
		PR string `json:"pr"`
		*Alias
	}{
		PR:    strconv.Itoa(v.PR),
		Alias: (*Alias)(v),
	})
}

// UnmarshalJSON custom unmarshaller to convert PR number
func (v *Version) UnmarshalJSON(data []byte) error {
	type Alias Version
	aux := struct {
		PR string `json:"pr"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if aux.PR != "" {
		v.PR, err = strconv.Atoi(aux.PR)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewVersion constructs a new Version
func NewVersion(p pullrequest.PullRequest) Version {
	return Version{
		PR:          p.Number,
		Commit:      p.HeadRef.OID,
		UpdatedDate: p.UpdatedAt,
	}
}

// PullRequestObject represents the GraphQL commit node.
// https://developer.github.com/v4/object/pullrequest/
type PullRequestObject struct {
	ID                string
	Number            int
	Title             string
	URL               string
	BaseRefName       string
	BaseRefOID        string
	HeadRefName       string
	IsCrossRepository bool
	CreatedAt         githubv4.DateTime
	UpdatedAt         githubv4.DateTime
	HeadRef           struct {
		ID     string
		Name   string
		Target struct {
			CommitObject `graphql:"... on Commit"`
		}
	}
	Repository struct {
		URL string
	}
	Labels struct {
		Edges []struct {
			Node struct {
				LabelObject
			}
		}
	} `graphql:"labels(first:100)"`
	Reviews struct {
		TotalCount int
	} `graphql:"reviews(states:APPROVED)"`
	TimelineItems struct {
		Edges []struct {
			Node struct {
				Typename            string `graphql:"__typename"`
				BaseRefChangedEvent struct {
					ID        string
					CreatedAt githubv4.DateTime
				} `graphql:"... on BaseRefChangedEvent"`
				BaseRefForcePushedEvent struct {
					ID        string
					CreatedAt githubv4.DateTime
				} `graphql:"... on BaseRefForcePushedEvent"`
				HeadRefForcePushedEvent struct {
					ID        string
					CreatedAt githubv4.DateTime
				} `graphql:"... on HeadRefForcePushedEvent"`
				IssueComment struct {
					ID        string
					CreatedAt githubv4.DateTime
					BodyText  string
				} `graphql:"... on IssueComment"`
				ReopenedEvent struct {
					ID        string
					CreatedAt githubv4.DateTime
				} `graphql:"... on ReopenedEvent"`
				PullRequestCommit struct {
					ID     string
					Commit CommitObject
				} `graphql:"... on PullRequestCommit"`
			}
		}
	} `graphql:"timelineItems(last:100,since:$s)"`
}

// CommitObject represents the GraphQL commit node.
// https://developer.github.com/v4/object/commit/
type CommitObject struct {
	ID             string
	OID            string
	AbbreviatedOID string
	AuthoredDate   githubv4.DateTime
	CommittedDate  githubv4.DateTime
	PushedDate     githubv4.DateTime
	Message        string
	Author         struct {
		User struct {
			Login string
		}
	}
}

// ChangedFileObject represents the GraphQL FilesChanged node.
// https://developer.github.com/v4/object/pullrequestchangedfile/
type ChangedFileObject struct {
	Path string
}

// LabelObject represents the GraphQL label node.
// https://developer.github.com/v4/object/label
type LabelObject struct {
	Name string
}
