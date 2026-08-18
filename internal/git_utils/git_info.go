package git_utils

import (
	"os"
	"os/exec"
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
