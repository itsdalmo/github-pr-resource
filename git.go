package resource

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
)

// Git interface for testing purposes.
//go:generate mockgen -destination=mocks/mock_git.go -package=mocks github.com/telia-oss/github-pr-resource Git
type Git interface {
	Init(string) error
	Pull(string, string) error
	RevParse(string) (string, error)
	Fetch(string, int) error
	Merge(string) error
}

// NewGitClient ...
func NewGitClient(source *Source, dir string, output io.Writer) (*GitClient, error) {
	return &GitClient{
		AccessToken: source.AccessToken,
		Directory:   dir,
		Output:      output,
	}, nil
}

// GitClient ...
type GitClient struct {
	AccessToken string
	Directory   string
	Output      io.Writer
}

func (g *GitClient) command(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.Dir = g.Directory
	cmd.Stdout = g.Output
	cmd.Stderr = g.Output
	return cmd
}

// Init ...
func (g *GitClient) Init(branch string) error {
	if err := g.command("git", "init").Run(); err != nil {
		return fmt.Errorf("init failed: %s", err)
	}
	if err := g.command("git", "checkout", "-b", branch).Run(); err != nil {
		return fmt.Errorf("checkout to '%s' failed: %s", branch, err)
	}
	if err := g.command("git", "config", "user.name", "concourse-ci").Run(); err != nil {
		return fmt.Errorf("failed to configure git user: %s", err)
	}
	if err := g.command("git", "config", "user.email", "concourse@local").Run(); err != nil {
		return fmt.Errorf("failed to configure git email: %s", err)
	}
	return nil
}

// Pull ...
func (g *GitClient) Pull(uri, branch string) error {
	endpoint, err := g.Endpoint(uri)
	if err != nil {
		return err
	}
	cmd := g.command("git", "pull", endpoint+".git", branch)

	// Discard output to have zero chance of logging the access token.
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clone failed: %s", err)
	}
	return nil
}

// RevParse retrieves the SHA of the given branch.
func (g *GitClient) RevParse(branch string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--verify", branch)
	cmd.Dir = g.Directory
	sha, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("rev-parse '%s' failed: %s: %s", branch, err, string(sha))
	}
	return strings.TrimSpace(string(sha)), nil
}

// Fetch ...
func (g *GitClient) Fetch(uri string, prNumber int) error {
	endpoint, err := g.Endpoint(uri)
	if err != nil {
		return err
	}
	cmd := g.command("git", "fetch", endpoint, fmt.Sprintf("pull/%s/head", strconv.Itoa(prNumber)))

	// Discard output to have zero chance of logging the access token.
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("fetch failed: %s", err)
	}
	return nil
}

// Merge ...
func (g *GitClient) Merge(sha string) error {
	if err := g.command("git", "merge", sha, "--no-stat").Run(); err != nil {
		return fmt.Errorf("merge failed: %s", err)
	}
	return nil
}

// Endpoint takes an uri and produces an endpoint with the login information baked in.
func (g *GitClient) Endpoint(uri string) (string, error) {
	endpoint, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("failed to parse commit url: %s", err)
	}
	endpoint.User = url.UserPassword("x-oauth-basic", g.AccessToken)
	return endpoint.String(), nil
}
