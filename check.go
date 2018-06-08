package resource

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Check (business logic)
func Check(request CheckRequest, manager Github) (CheckResponse, error) {
	var response CheckResponse

	pulls, err := manager.ListOpenPullRequests()
	if err != nil {
		return nil, fmt.Errorf("failed to get last commits: %s", err)
	}
	var disableSkipCI bool
	if request.Source.DisableCISkip != "" {
		disableSkipCI, err = strconv.ParseBool(request.Source.DisableCISkip)
		if err != nil {
			return nil, fmt.Errorf("failed to parse disable_ci_skip: %s", err)
		}
	}

	var newPullsToReturn []*PullRequest
	var alreadySeenPullsToHide []*PullRequest

Loop:
	for _, p := range pulls {
		// [ci skip]/[skip ci] in Pull request title
		if !disableSkipCI && ContainsSkipCI(p.Title) {
			continue
		}
		// [ci skip]/[skip ci] in Commit message
		if !disableSkipCI && ContainsSkipCI(p.Tip.Message) {
			continue
		}

		// Fetch files once if paths/ignore_paths are specified.
		var files []string

		if len(request.Source.Paths) > 0 || len(request.Source.IgnorePaths) > 0 {
			files, err = manager.ListModifiedFiles(p.Number)
			if err != nil {
				return nil, fmt.Errorf("failed to list modified files: %s", err)
			}
		}

		// Skip version if no files match the specified paths.
		if len(request.Source.Paths) > 0 {
			var wanted []string
			for _, pattern := range request.Source.Paths {
				w, err := FilterPath(files, pattern)
				if err != nil {
					return nil, fmt.Errorf("path match failed: %s", err)
				}
				wanted = append(wanted, w...)
			}
			if len(wanted) == 0 {
				continue Loop
			}
		}

		// Skip version if all files are ignored.
		if len(request.Source.IgnorePaths) > 0 {
			wanted := files
			for _, pattern := range request.Source.IgnorePaths {
				wanted, err = FilterIgnorePath(wanted, pattern)
				if err != nil {
					return nil, fmt.Errorf("ignore path match failed: %s", err)
				}
			}
			if len(wanted) == 0 {
				continue Loop
			}
		}

		// Determine above/below the fold
		if AboveTheFold(GetVersionStringFromPullRequest(p), request.Version.AlreadySeen) {
			newPullsToReturn = append(newPullsToReturn, p)
		} else {
			alreadySeenPullsToHide = append(alreadySeenPullsToHide, p)
		}
	}
	var combinedVersions Pulls = append(newPullsToReturn, alreadySeenPullsToHide...)
	sort.Sort(combinedVersions)
	var versionsJustSeen = GenerateVersion(combinedVersions)

	// Add "above-the-fold" with new alreadySeen version strings
	for _, p := range newPullsToReturn {
		response = append(response, NewVersion(p, versionsJustSeen))
	}
	// Sort the commits by date
	sort.Sort(response)

	// If there are no new but an old version = return the old
	if len(response) == 0 && request.Version.AlreadySeen != "" {
		response = append(response, request.Version)
	}
	// If there are new versions and no previous = return just the latest
	if len(response) != 0 && request.Version.AlreadySeen == "" {
		response = CheckResponse{response[len(response)-1]}
	}
	return response, nil
}

// GetVersionStringFromPullRequest returns string-serialized representation of latest commit in a PR
func GetVersionStringFromPullRequest(pull *PullRequest) string {
	return strconv.Itoa(pull.Number) + ":" + strconv.FormatInt(pull.Tip.CommittedDate.Time.Unix(), 10)
}

// ExtractVersionFromVersionString takes a string-formatted pair of PR#:CommittedDate and decodes them
func ExtractVersionFromVersionString(alreadySeenPair string) AlreadySeenVersion {
	var pairs = strings.Split(alreadySeenPair, ":")
	committedDateAsInt, err := strconv.ParseInt(pairs[1], 10, 64)
	if err != nil {
		panic(err)
	}
	var committedDate = time.Unix(committedDateAsInt, 0)
	return AlreadySeenVersion{PR: pairs[0], committedDate: committedDate}
}

// GenerateVersion returns a string-formatted array of PR#:CommittedDate
func GenerateVersion(pulls []*PullRequest) string {
	var pairs []string
	for _, p := range pulls {
		pairs = append(pairs, GetVersionStringFromPullRequest(p))
	}
	return strings.Join(pairs, ",")
}

// AboveTheFold returns a boolean indicating if a given pull request commit is newer than any commits already seen for that pull request
func AboveTheFold(pullRequestVersion string, alreadySeen string) bool {
	if !strings.Contains(alreadySeen, ":") {
		return true
	}
	var pairs = strings.Split(alreadySeen, ",")
	var isAboveTheFold = false
	var isFoundInPairs = false
	var pullRequest = ExtractVersionFromVersionString(pullRequestVersion)
	for _, pair := range pairs {
		var thisPairVersion = ExtractVersionFromVersionString(pair)
		if thisPairVersion.PR == pullRequest.PR {
			isFoundInPairs = true
			if pullRequest.committedDate.After(thisPairVersion.committedDate) {
				isAboveTheFold = true
			}
		}
	}
	if !isFoundInPairs {
		isAboveTheFold = true
	}
	return isAboveTheFold
}

// ContainsSkipCI returns true if a string contains [ci skip] or [skip ci].
func ContainsSkipCI(s string) bool {
	re := regexp.MustCompile("(?i)\\[(ci skip|skip ci)\\]")
	return re.MatchString(s)
}

// FilterIgnorePath ...
func FilterIgnorePath(files []string, pattern string) ([]string, error) {
	var out []string
	for _, file := range files {
		match, err := filepath.Match(pattern, file)
		if err != nil {
			return nil, err
		}
		if !match {
			out = append(out, file)
		}
	}
	return out, nil
}

// FilterPath ...
func FilterPath(files []string, pattern string) ([]string, error) {
	var out []string
	for _, file := range files {
		match, err := filepath.Match(pattern, file)
		if err != nil {
			return nil, err
		}
		if match {
			out = append(out, file)
		}
	}
	return out, nil
}

// CheckRequest ...
type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// CheckResponse ...
type CheckResponse []Version

func (r CheckResponse) Len() int {
	return len(r)
}

func (r CheckResponse) Less(i, j int) bool {
	return r[j].CommittedDate.After(r[i].CommittedDate)
}

func (r CheckResponse) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Pulls ...
type Pulls []*PullRequest

func (r Pulls) Len() int {
	return len(r)
}
func (r Pulls) Less(i, j int) bool {
	return r[j].Number > r[i].Number
}
func (r Pulls) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
