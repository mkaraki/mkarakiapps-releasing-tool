package scm_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GitHubApi struct{}

func (a *GitHubApi) Tagging(commitSha string, tagName string, tagMessage string, apiConfig *ScmApiConfig) error {
	token := apiConfig.GitHubToken
	if token == "" {
		return fmt.Errorf("GitHub API token is empty")
	}
	owner := apiConfig.RepoOwnerString
	if owner == "" {
		return fmt.Errorf("GitHub repository owner is empty")
	}
	repo := apiConfig.RepoNameString
	if repo == "" {
		return fmt.Errorf("GitHub repository name is empty")
	}

	if commitSha == "" {
		return fmt.Errorf("commit SHA is empty")
	}
	if tagName == "" {
		return fmt.Errorf("tag name is empty")
	}

	tagBody := map[string]string{
		"tag":     tagName,
		"message": tagMessage,
		"object":  commitSha,
		"type":    "commit",
	}

	tagJson, err := json.Marshal(tagBody)
	if err != nil {
		return fmt.Errorf("failed to marshal tag body: %v", err)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/tags", owner, repo)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(tagJson))
	if err != nil {
		return fmt.Errorf("failed to create request to GitHub API: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")

	resq, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to GitHub API: %v", err)
	}
	defer resq.Body.Close()

	if resq.StatusCode >= 300 {
		return fmt.Errorf("GitHub API returned error/redirect status code: %d", resq.StatusCode)
	}
	return nil
}
