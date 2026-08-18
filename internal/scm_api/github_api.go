package scm_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		if resq.StatusCode >= 300 {
			body, _ := io.ReadAll(resq.Body)
			return fmt.Errorf(
				"GitHub API returned %d: %s",
				resq.StatusCode,
				string(body),
			)
		}
	}

	type TagResponse struct {
		SHA string `json:"sha"`
	}

	var tagResp TagResponse
	if err := json.NewDecoder(resq.Body).Decode(&tagResp); err != nil {
		return fmt.Errorf("failed to decode tag response: %v", err)
	}

	refBody := map[string]string{
		"ref": "refs/tags/" + tagName,
		"sha": tagResp.SHA,
	}

	refJSON, err := json.Marshal(refBody)
	if err != nil {
		return fmt.Errorf("failed to marshal ref body: %v", err)
	}

	refURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/git/refs",
		owner,
		repo,
	)

	refReq, err := http.NewRequest(
		"POST",
		refURL,
		bytes.NewBuffer(refJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to create ref request: %v", err)
	}

	refReq.Header.Set("Authorization", "Bearer "+token)
	refReq.Header.Set("Content-Type", "application/json")
	refReq.Header.Set("Accept", "application/vnd.github+json")

	refRes, err := http.DefaultClient.Do(refReq)
	if err != nil {
		return fmt.Errorf("failed to create Git ref: %v", err)
	}
	defer refRes.Body.Close()

	if refRes.StatusCode >= 300 {
		body, _ := io.ReadAll(refRes.Body)
		return fmt.Errorf(
			"GitHub API returned error creating ref: %d: %s",
			refRes.StatusCode,
			string(body),
		)
	}

	return nil
}
