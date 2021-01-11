package resource

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
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
	Pull(string, string, int, bool, bool) error
	RevParse(string) (string, error)
	Fetch(string, int, int, bool) error
	Checkout(string, string, bool) error
	Merge(string, bool) error
	Rebase(string, string, bool) error
	GitCryptUnlock(string) error
}

// NewGitClient ...
func NewGitClient(source *Source, dir string, output io.Writer) (*GitClient, error) {
	if source.SkipSSLVerification {
		os.Setenv("GIT_SSL_NO_VERIFY", "true")
	}
	if source.DisableGitLFS {
		os.Setenv("GIT_LFS_SKIP_SMUDGE", "true")
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
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		"X_OAUTH_BASIC_TOKEN="+g.AccessToken,
		"GIT_ASKPASS=/usr/local/bin/askpass.sh")
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
	if err := g.command("git", "config", "url.https://x-oauth-basic@github.com/.insteadOf", "git@github.com:").Run(); err != nil {
		return fmt.Errorf("failed to configure github url: %s", err)
	}
	if err := g.command("git", "config", "url.https://.insteadOf", "git://").Run(); err != nil {
		return fmt.Errorf("failed to configure github url: %s", err)
	}
	return nil
}

// Pull ...
func (g *GitClient) Pull(uri, branch string, depth int, submodules bool, fetchTags bool) error {
	endpoint, err := g.Endpoint(uri)
	if err != nil {
		return err
	}

	if err := g.command("git", "remote", "add", "origin", endpoint).Run(); err != nil {
		return fmt.Errorf("setting 'origin' remote to '%s' failed: %s", endpoint, err)
	}

	args := []string{"pull", "origin", branch}
	if depth > 0 {
		args = append(args, "--depth", strconv.Itoa(depth))
	}
	if fetchTags {
		args = append(args, "--tags")
	}
	if submodules {
		args = append(args, "--recurse-submodules")
	}
	cmd := g.command("git", args...)

	// Discard output to have zero chance of logging the access token.
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pull failed: %s", cmd)
	}
	if submodules {
		submodulesGet := g.command("git", "submodule", "update", "--init", "--recursive")
		if err := submodulesGet.Run(); err != nil {
			return fmt.Errorf("submodule update failed: %s", err)
		}
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
func (g *GitClient) Fetch(uri string, prNumber int, depth int, submodules bool) error {
	endpoint, err := g.Endpoint(uri)
	if err != nil {
		return err
	}

	args := []string{"fetch", endpoint, fmt.Sprintf("pull/%s/head", strconv.Itoa(prNumber))}
	if depth > 0 {
		args = append(args, "--depth", strconv.Itoa(depth))
	}
	if submodules {
		args = append(args, "--recurse-submodules")
	}
	cmd := g.command("git", args...)

	// Discard output to have zero chance of logging the access token.
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("fetch failed: %s", err)
	}
	return nil
}

// CheckOut
func (g *GitClient) Checkout(branch, sha string, submodules bool) error {
	if err := g.command("git", "checkout", "-b", branch, sha).Run(); err != nil {
		return fmt.Errorf("checkout failed: %s", err)
	}

	if submodules {
		if err := g.command("git", "submodule", "update", "--init", "--recursive", "--checkout").Run(); err != nil {
			return fmt.Errorf("submodule update failed: %s", err)
		}
	}

	return nil
}

// Merge ...
func (g *GitClient) Merge(sha string, submodules bool) error {
	if err := g.command("git", "merge", sha, "--no-stat").Run(); err != nil {
		return fmt.Errorf("merge failed: %s", err)
	}

	if submodules {
		if err := g.command("git", "submodule", "update", "--init", "--recursive", "--merge").Run(); err != nil {
			return fmt.Errorf("submodule update failed: %s", err)
		}
	}

	return nil
}

// Rebase ...
func (g *GitClient) Rebase(baseRef string, headSha string, submodules bool) error {
	if err := g.command("git", "rebase", baseRef, headSha).Run(); err != nil {
		return fmt.Errorf("rebase failed: %s", err)
	}

	if submodules {
		if err := g.command("git", "submodule", "update", "--init", "--recursive", "--rebase").Run(); err != nil {
			return fmt.Errorf("submodule update failed: %s", err)
		}
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
	if err := ioutil.WriteFile(keyPath, decodedKey, os.FileMode(0600)); err != nil {
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
