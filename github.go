package resource

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/shurcooL/githubv4"
	"github.com/telia-oss/github-pr-resource/pullrequest"
	"golang.org/x/oauth2"
)

// Github for testing purposes.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o fakes/fake_github.go . Github
type Github interface {
	ListOpenPullRequests(prSince time.Time) ([]pullrequest.PullRequest, error)
	PostComment(int, string) error
	GetPullRequest(int, string) (pullrequest.PullRequest, error)
	GetChangedFiles(int) ([]string, error)
	UpdateCommitStatus(string, string, string, string, string, string) error
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

	ctx := context.TODO()
	httpClient := http.Client{}

	// Skip SSL verification for self-signed certificates
	// source: https://github.com/google/go-github/pull/598#issuecomment-333039238
	if s.SkipSSLVerification {
		log.Println("disabling SSL verification")
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &httpClient)
	}

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.AccessToken},
	))

	if s.PreviewSchema {
		log.Println("attaching preview schema transport to client")
		client.Transport = &PreviewSchemaTransport{
			oauthTransport: client.Transport,
		}
	}

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

// ListOpenPullRequests gets the last commit on all open pull requests
func (m *GithubClient) ListOpenPullRequests(since time.Time) ([]pullrequest.PullRequest, error) {
	return m.searchOpenPullRequests(since, 100)
}

func (m *GithubClient) searchOpenPullRequests(since time.Time, number int) ([]pullrequest.PullRequest, error) {
	log.Println("building open pull requests query")

	var query struct {
		Search struct {
			Edges []struct {
				Node struct {
					PullRequestObject `graphql:"... on PullRequest"`
				}
			}
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"search(query:$q,type:ISSUE,last:$n,after:$c)"`
	}

	vars := map[string]interface{}{
		"c": (*githubv4.String)(nil),
		"s": githubv4.DateTime{Time: since},
		"n": githubv4.Int(number),
		"q": githubv4.String(fmt.Sprintf("is:pr is:open repo:%s/%s updated:>%s sort:updated", m.Owner, m.Repository, since.Format(time.RFC3339))),
	}

	var response []pullrequest.PullRequest
	for {
		if err := m.V4.Query(context.TODO(), &query, vars); err != nil {
			return nil, err
		}
		for _, p := range query.Search.Edges {
			response = append(response, PullRequestFactory(p.Node.PullRequestObject))
		}
		if number < 100 || !query.Search.PageInfo.HasNextPage {
			break
		}
		vars["c"] = query.Search.PageInfo.EndCursor
	}
	return response, nil
}

// PostComment to a pull request or issue.
func (m *GithubClient) PostComment(number int, comment string) error {
	_, _, err := m.V3.Issues.CreateComment(
		context.TODO(),
		m.Owner,
		m.Repository,
		number,
		&github.IssueComment{
			Body: github.String(comment),
		},
	)
	return err
}

// GetChangedFiles ...
func (m *GithubClient) GetChangedFiles(number int) ([]string, error) {
	log.Println("building pull request changed files query")

	var filequery struct {
		Repository struct {
			PullRequest struct {
				Files struct {
					Edges []struct {
						Node struct {
							ChangedFileObject
						}
					} `graphql:"edges"`
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					} `graphql:"pageInfo"`
				} `graphql:"files(first:100, after: $c)"`
			} `graphql:"pullRequest(number:$n)"`
		} `graphql:"repository(owner:$owner,name:$name)"`
	}

	files := []string{}
	cursor := ""

	for {
		vars := map[string]interface{}{
			"owner": githubv4.String(m.Owner),
			"name":  githubv4.String(m.Repository),
			"n":     githubv4.Int(number),
			"c":     githubv4.String(cursor),
		}

		if err := m.V4.Query(context.TODO(), &filequery, vars); err != nil {
			return nil, err
		}

		for _, f := range filequery.Repository.PullRequest.Files.Edges {
			files = append(files, f.Node.Path)
		}

		if !filequery.Repository.PullRequest.Files.PageInfo.HasNextPage {
			break
		}

		cursor = string(filequery.Repository.PullRequest.Files.PageInfo.EndCursor)
	}

	return files, nil
}

// GetPullRequest ...
func (m *GithubClient) GetPullRequest(number int, commitRef string) (pullrequest.PullRequest, error) {
	log.Println("building pull request query")

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
				} `graphql:"commits(last:$last)"`
			} `graphql:"pullRequest(number:$number)"`
		} `graphql:"repository(owner:$owner,name:$name)"`
	}

	vars := map[string]interface{}{
		"s":      githubv4.DateTime{Time: time.Now().AddDate(-1, 0, 0)},
		"owner":  githubv4.String(m.Owner),
		"name":   githubv4.String(m.Repository),
		"number": githubv4.Int(number),
		"last":   githubv4.Int(100),
	}

	// TODO: Pagination - in case someone pushes > 100 commits before the build has time to start :p
	if err := m.V4.Query(context.TODO(), &query, vars); err != nil {
		return pullrequest.PullRequest{}, err
	}

	for _, c := range query.Repository.PullRequest.Commits.Edges {
		if c.Node.Commit.OID == commitRef {
			// Return as soon as we find the correct ref.
			pull := PullRequestFactory(query.Repository.PullRequest.PullRequestObject)
			pull.HeadRef = commitFactory(c.Node.Commit)
			return pull, nil
		}
	}

	// Return an error if the commit was not found
	return pullrequest.PullRequest{}, fmt.Errorf("commit with ref '%s' does not exist", commitRef)
}

// UpdateCommitStatus for a given commit (not supported by V4 API).
func (m *GithubClient) UpdateCommitStatus(commitRef, baseContext, statusContext, status, targetURL, description string) error {
	if baseContext == "" {
		baseContext = "concourse-ci"
	}

	if statusContext == "" {
		statusContext = "status"
	}

	if targetURL == "" {
		targetURL = strings.Join([]string{os.Getenv("ATC_EXTERNAL_URL"), "builds", os.Getenv("BUILD_ID")}, "/")
	}

	if description == "" {
		description = fmt.Sprintf("Concourse CI build %s", status)
	}

	_, _, err := m.V3.Repositories.CreateStatus(
		context.TODO(),
		m.Owner,
		m.Repository,
		commitRef,
		&github.RepoStatus{
			State:       github.String(strings.ToLower(status)),
			TargetURL:   github.String(targetURL),
			Description: github.String(description),
			Context:     github.String(path.Join(baseContext, statusContext)),
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

// PullRequestFactory generates a PullRequest object from a PullRequestObject
func PullRequestFactory(p PullRequestObject) pullrequest.PullRequest {
	labels := make([]string, 0)

	for _, i := range p.Labels.Edges {
		labels = append(labels, i.Node.LabelObject.Name)
	}

	events := make([]pullrequest.Event, 0)
	comments := make([]pullrequest.Comment, 0)
	commits := make([]pullrequest.Commit, 0)

	for _, i := range p.TimelineItems.Edges {
		switch i.Node.Typename {
		case pullrequest.BaseRefChangedEvent:
			events = append(events, pullrequest.Event{
				Type:      pullrequest.BaseRefChangedEvent,
				CreatedAt: i.Node.BaseRefChangedEvent.CreatedAt.Time,
			})
		case pullrequest.BaseRefForcePushedEvent:
			events = append(events, pullrequest.Event{
				Type:      pullrequest.BaseRefForcePushedEvent,
				CreatedAt: i.Node.BaseRefForcePushedEvent.CreatedAt.Time,
			})
		case pullrequest.HeadRefForcePushedEvent:
			events = append(events, pullrequest.Event{
				Type:      pullrequest.HeadRefForcePushedEvent,
				CreatedAt: i.Node.HeadRefForcePushedEvent.CreatedAt.Time,
			})
		case pullrequest.ReopenedEvent:
			events = append(events, pullrequest.Event{
				Type:      pullrequest.ReopenedEvent,
				CreatedAt: i.Node.ReopenedEvent.CreatedAt.Time,
			})
		case pullrequest.IssueComment:
			comments = append(comments, pullrequest.Comment{
				CreatedAt: i.Node.IssueComment.CreatedAt.Time,
				Body:      i.Node.IssueComment.BodyText,
			})
		case pullrequest.PullRequestCommit:
			commits = append(commits, commitFactory(i.Node.PullRequestCommit.Commit))
		}
	}

	return pullrequest.PullRequest{
		ID:                  p.ID,
		Number:              p.Number,
		Title:               p.Title,
		URL:                 p.URL,
		RepositoryURL:       p.Repository.URL,
		BaseRefName:         p.BaseRefName,
		BaseRefOID:          p.BaseRefOID,
		HeadRefName:         p.HeadRefName,
		IsCrossRepository:   p.IsCrossRepository,
		CreatedAt:           p.CreatedAt.Time,
		UpdatedAt:           p.UpdatedAt.Time,
		HeadRef:             commitFactory(p.HeadRef.Target.CommitObject),
		Events:              events,
		Commits:             commits,
		Comments:            comments,
		Labels:              labels,
		ApprovedReviewCount: p.Reviews.TotalCount,
	}
}

func commitFactory(c CommitObject) pullrequest.Commit {
	return pullrequest.Commit{
		OID:            c.OID,
		AbbreviatedOID: c.AbbreviatedOID,
		AuthoredDate:   c.AuthoredDate.Time,
		CommittedDate:  c.CommittedDate.Time,
		PushedDate:     c.PushedDate.Time,
		Message:        c.Message,
		Author:         c.Author.User.Login,
	}
}

// PreviewSchemaTransport is used to access GraphQL schema's hidden behind an Accept header by GitHub
type PreviewSchemaTransport struct {
	oauthTransport http.RoundTripper
}

// RoundTrip appends the Accept header and then executes the parent RoundTrip Transport
func (t *PreviewSchemaTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	log.Println("setting accept header for timelineItems & files connections preview schemas")
	r.Header.Add("Accept", "application/vnd.github.starfire-preview+json, application/vnd.github.ocelot-preview+json")

	return t.oauthTransport.RoundTrip(r)
}
