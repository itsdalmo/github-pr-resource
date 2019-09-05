package resource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/telia-oss/github-pr-resource/pullrequest"
)

// Get (business logic)
func Get(request GetRequest, github Github, git Git, outputDir string) (*GetResponse, error) {
	if request.Params.SkipDownload {
		return &GetResponse{Version: request.Version}, nil
	}

	pull, err := github.GetPullRequest(request.Version.PR, request.Version.Commit)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pull request: %s", err)
	}

	// Initialize and pull the base for the PR
	err = getRepository(request, git, pull)
	if err != nil {
		return nil, fmt.Errorf("failed to get / clone repository: %s", err)
	}

	// Get the last commit SHA in base for the metadata
	baseSHA, err := git.RevParse(pull.BaseRefName)
	if err != nil {
		return nil, err
	}

	// Fetch the PR and merge the specified commit into the base
	if err := git.Fetch(pull.Number, request.Params.GitDepth); err != nil {
		return nil, err
	}

	switch tool := request.Params.IntegrationTool; tool {
	case "rebase":
		if err := git.Rebase(pull.BaseRefName, request.Version.Commit); err != nil {
			return nil, err
		}
	case "merge", "":
		if err := git.Merge(request.Version.Commit); err != nil {
			return nil, err
		}
	case "checkout":
		if err := git.Checkout(pull.HeadRefName, request.Version.Commit); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid integration tool specified: %s", tool)
	}

	if request.Source.GitCryptKey != "" {
		if err := git.GitCryptUnlock(request.Source.GitCryptKey); err != nil {
			return nil, err
		}
	}

	// Create the metadata
	var metadata Metadata
	metadata.Add("pr", strconv.Itoa(pull.Number))
	metadata.Add("url", pull.URL)
	metadata.Add("head_name", pull.HeadRefName)
	metadata.Add("head_sha", pull.HeadRef.OID)
	metadata.Add("head_short_sha", pull.HeadRef.AbbreviatedOID)
	metadata.Add("base_name", pull.BaseRefName)
	metadata.Add("base_sha", baseSHA)
	metadata.Add("message", pull.HeadRef.Message)
	metadata.Add("author", pull.HeadRef.Author)
	metadata.Add("events", fmt.Sprintf("%v", pull.Events))

	// Write version and metadata for reuse in PUT
	path := filepath.Join(outputDir, ".git", "resource")
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %s", err)
	}
	b, err := json.Marshal(&request.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal version: %s", err)
	}
	if err := ioutil.WriteFile(filepath.Join(path, "version.json"), b, 0644); err != nil {
		return nil, fmt.Errorf("failed to write version: %s", err)
	}
	b, err = json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %s", err)
	}
	if err := ioutil.WriteFile(filepath.Join(path, "metadata.json"), b, 0644); err != nil {
		return nil, fmt.Errorf("failed to write metadata: %s", err)
	}

	for _, d := range metadata {
		filename := d.Name
		content := []byte(d.Value)
		if err := ioutil.WriteFile(filepath.Join(path, filename), content, 0644); err != nil {
			return nil, fmt.Errorf("failed to write metadata file %s: %s", filename, err)
		}

	}

	if request.Params.ListChangedFiles {
		cfol, err := github.GetChangedFiles(request.Version.PR)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch list of changed files: %s", err)
		}

		var fl []byte

		for _, v := range cfol {
			fl = append(fl, []byte(v+"\n")...)
		}

		// Create List with changed files
		if err := ioutil.WriteFile(filepath.Join(path, "changed_files"), fl, 0644); err != nil {
			return nil, fmt.Errorf("failed to write file list: %s", err)
		}
	}

	return &GetResponse{
		Version:  request.Version,
		Metadata: metadata,
	}, nil
}

func getRepository(request GetRequest, git Git, pull pullrequest.PullRequest) error {
	if request.Params.Clone {
		if err := git.Clone(pull.RepositoryURL, pull.BaseRefName, request.Params.GitDepth); err != nil {
			return err
		}

		return nil
	}

	if err := git.Init(pull.BaseRefName); err != nil {
		return err
	}
	if err := git.Pull(pull.RepositoryURL, pull.BaseRefName, request.Params.GitDepth); err != nil {
		return err
	}

	return nil
}

// GetParameters ...
type GetParameters struct {
	SkipDownload     bool   `json:"skip_download"`
	IntegrationTool  string `json:"integration_tool"`
	GitDepth         int    `json:"git_depth"`
	ListChangedFiles bool   `json:"list_changed_files"`
	Clone            bool   `json:"clone"`
}

// GetRequest ...
type GetRequest struct {
	Source  Source        `json:"source"`
	Version Version       `json:"version"`
	Params  GetParameters `json:"params"`
}

// GetResponse ...
type GetResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata,omitempty"`
}
