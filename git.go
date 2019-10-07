package resource

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Git interface for testing purposes.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o fakes/fake_git.go . Git
type Git interface {
	Init(string) error
	Pull(string, string, int) error
	Clone(string, string, int) error
	RevParse(string) (string, error)
	Fetch(int, int) error
	Checkout(string, string) error
	Merge(string) error
	Rebase(string, string) error
	GitCryptUnlock(string) error
}

// NewGitClient ...
func NewGitClient(source *Source, dir string, output io.Writer) (*GitClient, error) {
	if source.SkipSSLVerification {
		os.Setenv("GIT_SSL_NO_VERIFY", "true")
	}
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

	log.Println("initialized repository:", branch)

	return g.Config()
}

// Pull ...
func (g *GitClient) Pull(uri, branch string, depth int) error {
	endpoint, err := g.Endpoint(uri)
	if err != nil {
		return err
	}

	if err := g.command("git", "remote", "add", "origin", endpoint+".git").Run(); err != nil {
		return fmt.Errorf("failed to set remote origin: %s", err)
	}

	args := []string{"pull", "origin", branch}
	args = appendDepth(args, depth)
	cmd := g.command("git", args...)

	// Discard output to have zero chance of logging the access token.
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard

	log.Println("pulling baseref:", args)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clone failed: %s", err)
	}

	return nil
}

// Config ...
func (g *GitClient) Config() error {
	if err := g.command("git", "config", "user.name", "concourse-ci").Run(); err != nil {
		return fmt.Errorf("failed to configure git user: %s", err)
	}
	if err := g.command("git", "config", "user.email", "concourse@local").Run(); err != nil {
		return fmt.Errorf("failed to configure git email: %s", err)
	}

	return nil
}

// Clone ...
func (g *GitClient) Clone(uri, branch string, depth int) error {
	endpoint, err := g.Endpoint(uri)
	if err != nil {
		return err
	}

	args := []string{"clone", endpoint + ".git", "-b", branch, "."}
	args = appendDepth(args, depth)
	cmd := g.command("git", args...)

	// Discard output to have zero chance of logging the access token.
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard

	log.Println("cloning baseref:", args)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clone failed: %s", err)
	}

	return g.Config()
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
func (g *GitClient) Fetch(prNumber, depth int) error {
	args := []string{
		"fetch",
		"origin",
		"-q",
		fmt.Sprintf("pull/%s/head", strconv.Itoa(prNumber)),
	}
	args = appendDepth(args, depth)
	cmd := g.command("git", args...)

	// Discard output to have zero chance of logging the access token.
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard

	log.Println("fetching headref:", args)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("fetch failed: %s", err)
	}
	return nil
}

// Checkout ...
func (g *GitClient) Checkout(branch, sha string) error {
	log.Println("checkout:", branch, sha)
	if err := g.command("git", "checkout", "-b", "pr-"+branch, sha).Run(); err != nil {
		return fmt.Errorf("checkout failed: %s", err)
	}

	return nil
}

// Merge ...
func (g *GitClient) Merge(sha string) error {
	log.Println("merging sha:", sha)
	if err := g.command("git", "merge", sha, "--no-stat").Run(); err != nil {
		return fmt.Errorf("merge failed: %s", err)
	}
	return nil
}

// Rebase ...
func (g *GitClient) Rebase(baseRef string, headSha string) error {
	log.Println("rebasing:", baseRef, headSha)
	if err := g.command("git", "rebase", baseRef, headSha).Run(); err != nil {
		return fmt.Errorf("rebase failed: %s", err)
	}
	return nil
}

// GitCryptUnlock unlocks the repository using git-crypt
func (g *GitClient) GitCryptUnlock(base64key string) error {
	keyDir, err := ioutil.TempDir("", "")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory")
	}
	defer os.RemoveAll(keyDir)
	decodedKey, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return fmt.Errorf("failed to decode git-crypt key")
	}
	keyPath := filepath.Join(keyDir, "git-crypt-key")
	if err := ioutil.WriteFile(keyPath, decodedKey, 600); err != nil {
		return fmt.Errorf("failed to write git-crypt key to file: %s", err)
	}
	if err := g.command("git-crypt", "unlock", keyPath).Run(); err != nil {
		return fmt.Errorf("git-crypt unlock failed: %s", err)
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

func appendDepth(args []string, depth int) []string {
	if depth > 0 {
		args = append(args, []string{"--depth", strconv.Itoa(depth)}...)
	}

	return args
}
