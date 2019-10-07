package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	resource "github.com/telia-oss/github-pr-resource"
	rlog "github.com/telia-oss/github-pr-resource/log"
)

func main() {
	var request resource.GetRequest

	input := rlog.WriteStdin()
	defer rlog.Close()

	if err := json.Unmarshal(input, &request); err != nil {
		log.Fatalf("failed to unmarshal request: %s", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("missing arguments")
	}
	outputDir := os.Args[1]
	if err := request.Source.Validate(); err != nil {
		log.Fatalf("invalid source configuration: %s", err)
	}
	git, err := resource.NewGitClient(&request.Source, outputDir, os.Stderr)
	if err != nil {
		log.Fatalf("failed to create git client: %s", err)
	}
	github, err := resource.NewGithubClient(&request.Source)
	if err != nil {
		log.Fatalf("failed to create github manager: %s", err)
	}
	response, err := resource.Get(request, github, git, outputDir)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("get failed: %s", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to marshal response: %s", err)
	}
}
