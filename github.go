package resource

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Github for testing purposes.
//go:generate counterfeiter -o fakes/fake_github.go . Github
type Github interface {
	ListOpenPullRequests() ([]*PullRequest, error)
	ListModifiedFiles(int) ([]string, error)
	PostComment(string, string) error
	GetPullRequest(string, string) (*PullRequest, error)
	UpdateCommitStatus(string, string, string) error
}

// GithubClient for handling requests to the Github V3 and V4 APIs.
type GithubClient struct {
	V3         *github.Client
	V4         *githubv4.Client
	Repository string
	Owner      string
}

// NewGithubClient ...
func NewGithubClient(s *Source) (*GithubClient, error) {
	owner, repository, err := parseRepository(s.Repository)
	if err != nil {
		return nil, err
	}

	// Skip SSL verification for self-signed certificates
	// source: https://github.com/google/go-github/pull/598#issuecomment-333039238
	var ctx context.Context
	if s.SkipSSLVerification {
		insecureClient := &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		}
		ctx = context.WithValue(context.TODO(), oauth2.HTTPClient, insecureClient)
	} else {
		ctx = context.TODO()
	}

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.AccessToken},
	))

	var v3 *github.Client
	if s.V3Endpoint != "" {
		endpoint, err := url.Parse(s.V3Endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to parse v3 endpoint: %s", err)
		}
		v3, err = github.NewEnterpriseClient(endpoint.String(), endpoint.String(), client)
		if err != nil {
			return nil, err
		}
	} else {
		v3 = github.NewClient(client)
	}

	var v4 *githubv4.Client
	if s.V4Endpoint != "" {
		endpoint, err := url.Parse(s.V4Endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to parse v4 endpoint: %s", err)
		}
		v4 = githubv4.NewEnterpriseClient(endpoint.String(), client)
		if err != nil {
			return nil, err
		}
	} else {
		v4 = githubv4.NewClient(client)
	}

	return &GithubClient{
		V3:         v3,
		V4:         v4,
		Owner:      owner,
		Repository: repository,
	}, nil
}

// ListOpenPullRequests gets the last commit on all open pull requests.
func (m *GithubClient) ListOpenPullRequests() ([]*PullRequest, error) {
	var query struct {
		Repository struct {
			PullRequests struct {
				Edges []struct {
					Node struct {
						PullRequestObject
						Commits struct {
							Edges []struct {
								Node struct {
									Commit CommitObject
								}
							}
						} `graphql:"commits(last:$commitsLast)"`
					}
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"pullRequests(first:$prFirst,states:$prStates,after:$prCursor)"`
		} `graphql:"repository(owner:$repositoryOwner,name:$repositoryName)"`
	}

	vars := map[string]interface{}{
		"repositoryOwner": githubv4.String(m.Owner),
		"repositoryName":  githubv4.String(m.Repository),
		"prFirst":         githubv4.Int(100),
		"prStates":        []githubv4.PullRequestState{githubv4.PullRequestStateOpen},
		"prCursor":        (*githubv4.String)(nil),
		"commitsLast":     githubv4.Int(1),
	}

	var response []*PullRequest
	for {
		if err := m.V4.Query(context.TODO(), &query, vars); err != nil {
			return nil, err
		}
		for _, p := range query.Repository.PullRequests.Edges {
			for _, c := range p.Node.Commits.Edges {
				response = append(response, &PullRequest{
					PullRequestObject: p.Node.PullRequestObject,
					Tip:               c.Node.Commit,
				})
			}
		}
		if !query.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}
		vars["prCursor"] = query.Repository.PullRequests.PageInfo.EndCursor
	}
	return response, nil
}

// ListModifiedFiles in a pull request (not supported by V4 API).
func (m *GithubClient) ListModifiedFiles(prNumber int) ([]string, error) {
	var files []string

	opt := &github.ListOptions{
		PerPage: 100,
	}
	for {
		result, response, err := m.V3.PullRequests.ListFiles(
			context.TODO(),
			m.Owner,
			m.Repository,
			prNumber,
			opt,
		)
		if err != nil {
			return nil, err
		}
		for _, f := range result {
			files = append(files, *f.Filename)
		}
		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}
	return files, nil
}

// PostComment to a pull request or issue.
func (m *GithubClient) PostComment(prNumber, comment string) error {
	pr, err := strconv.Atoi(prNumber)
	if err != nil {
		return fmt.Errorf("failed to convert pull request number to int: %s", err)
	}

	_, _, err = m.V3.Issues.CreateComment(
		context.TODO(),
		m.Owner,
		m.Repository,
		pr,
		&github.IssueComment{
			Body: github.String(comment),
		},
	)
	return err
}

// GetPullRequest ...
func (m *GithubClient) GetPullRequest(prNumber, commitRef string) (*PullRequest, error) {
	pr, err := strconv.Atoi(prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to convert pull request number to int: %s", err)
	}

	var query struct {
		Repository struct {
			PullRequest struct {
				PullRequestObject
				Commits struct {
					Edges []struct {
						Node struct {
							Commit CommitObject
						}
					}
				} `graphql:"commits(last:$commitsLast)"`
			} `graphql:"pullRequest(number:$prNumber)"`
		} `graphql:"repository(owner:$repositoryOwner,name:$repositoryName)"`
	}

	vars := map[string]interface{}{
		"repositoryOwner": githubv4.String(m.Owner),
		"repositoryName":  githubv4.String(m.Repository),
		"prNumber":        githubv4.Int(pr),
		"commitsLast":     githubv4.Int(100),
	}

	// TODO: Pagination - in case someone pushes > 100 commits before the build has time to start :p
	if err := m.V4.Query(context.TODO(), &query, vars); err != nil {
		return nil, err
	}
	for _, c := range query.Repository.PullRequest.Commits.Edges {
		if c.Node.Commit.OID == commitRef {
			// Return as soon as we find the correct ref.
			return &PullRequest{
				PullRequestObject: query.Repository.PullRequest.PullRequestObject,
				Tip:               c.Node.Commit,
			}, nil
		}
	}

	// Return an error if the commit was not found
	return nil, fmt.Errorf("commit with ref '%s' does not exist", commitRef)
}

// UpdateCommitStatus for a given commit (not supported by V4 API).
func (m *GithubClient) UpdateCommitStatus(commitRef, statusContext, status string) error {
	c := []string{"concourse-ci"}
	if statusContext == "" {
		c = append(c, "status")
	} else {
		c = append(c, statusContext)
	}
	statusContext = strings.Join(c, "/")

	// Format build page
	build := os.Getenv("ATC_EXTERNAL_URL")
	if build != "" {
		build = strings.Join([]string{build, "builds", os.Getenv("BUILD_ID")}, "/")
	}

	_, _, err := m.V3.Repositories.CreateStatus(
		context.TODO(),
		m.Owner,
		m.Repository,
		commitRef,
		&github.RepoStatus{
			State:       github.String(strings.ToLower(status)),
			TargetURL:   github.String(build),
			Description: github.String(fmt.Sprintf("Concourse CI build %s", status)),
			Context:     github.String(statusContext),
		},
	)
	return err
}

func parseRepository(s string) (string, string, error) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return "", "", errors.New("malformed repository")
	}
	return parts[0], parts[1], nil
}
