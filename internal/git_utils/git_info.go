package git_utils

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func getEnv(key string, defaultValue string) string {
	value := strings.TrimSpace(strings.Trim(os.Getenv(key), ""))
	if value == "" {
		return defaultValue
	}
	return value
}

func GetCurrentBranchName() (string, error) {
	// Run git command to get current branch name
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	branchName := strings.TrimSpace(string(output))

	if branchName == "" || branchName == "detached" {
		// Try to use env variable from GitHub Actions
		branchName = strings.TrimSpace(getEnv("GITHUB_HEAD_REF", ""))
	}

	if branchName == "" {
		// Try to use env variable from GitLab CI
		branchName = strings.TrimSpace(getEnv("CI_COMMIT_REF_NAME", ""))
	}

	return branchName, nil
}

func GetCurrentCommitSHA() (string, error) {
	// Run git command to get current commit SHA
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current commit SHA: %v. Output: %s", err, output)
	}

	commitSHA := strings.TrimSpace(string(output))

	shaRe, err := regexp.Compile(`^[0-9a-f]{40}$`)
	if err != nil {
		return "", fmt.Errorf("failed to compile SHA regex: %v", err)
	}

	if !shaRe.MatchString(commitSHA) {
		commitSHA = ""
	}

	if commitSHA == "" {
		// Try to use env variable from GitHub Actions
		commitSHA = strings.TrimSpace(getEnv("GITHUB_SHA", ""))
	}

	if commitSHA == "" {
		// Try to use env variable from GitLab CI
		commitSHA = strings.TrimSpace(getEnv("CI_COMMIT_SHA", ""))
	}

	if commitSHA == "" {
		return "", fmt.Errorf("failed to get current commit SHA")
	}

	return commitSHA, nil
}
