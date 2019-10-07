package resource_test

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	"github.com/shurcooL/githubv4"
	resource "github.com/telia-oss/github-pr-resource"
	_ "github.com/telia-oss/github-pr-resource/log"
	"github.com/telia-oss/github-pr-resource/pullrequest"
)

func createTestPR(count int, baseName string, skipCI, isCrossRepo, created, nocommit bool, approvedReviews int, labels []string) pullrequest.PullRequest {
	n := strconv.Itoa(count)
	u := time.Now().AddDate(0, 0, count)

	c := u
	if !created {
		c = time.Now().AddDate(0, 0, count-1)
	}

	m := fmt.Sprintf("commit message%s", n)
	if skipCI {
		m = "[skip ci]" + m
	}

	commit := resource.CommitObject{
		ID:             fmt.Sprintf("commit%s", n),
		OID:            fmt.Sprintf("oid%s", n),
		AbbreviatedOID: fmt.Sprintf("oid%s", n),
		CommittedDate:  githubv4.DateTime{Time: u},
		Message:        m,
		Author: struct{ User struct{ Login string } }{
			User: struct{ Login string }{
				Login: fmt.Sprintf("login%s", n),
			},
		},
	}
	if nocommit {
		commit.ID = ""
	}

	pr := resource.PullRequestFactory(resource.PullRequestObject{
		ID:                fmt.Sprintf("pr%s", n),
		Number:            count,
		Title:             fmt.Sprintf("pr%s title", n),
		URL:               fmt.Sprintf("pr%s url", n),
		BaseRefName:       baseName,
		BaseRefOID:        "sha",
		HeadRefName:       fmt.Sprintf("pr%s", n),
		IsCrossRepository: isCrossRepo,
		CreatedAt:         githubv4.DateTime{Time: c},
		UpdatedAt:         githubv4.DateTime{Time: u},
		HeadRef: struct {
			ID     string
			Name   string
			Target struct {
				resource.CommitObject `graphql:"... on Commit"`
			}
		}{
			ID:   fmt.Sprintf("commit%s", n),
			Name: fmt.Sprintf("pr%s", n),
			Target: struct {
				resource.CommitObject `graphql:"... on Commit"`
			}{commit},
		},
		Repository: struct{ URL string }{
			URL: fmt.Sprintf("repo%s url", n),
		},
	})

	pr.ApprovedReviewCount = approvedReviews
	pr.Labels = labels

	return pr
}

func createTestDirectory(t *testing.T) string {
	dir, err := ioutil.TempDir("", "github-pr-resource")
	if err != nil {
		t.Fatalf("failed to create temporary directory")
	}
	return dir
}

func readTestFile(t *testing.T, path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read: %s: %s", path, err)
	}
	return string(b)
}
