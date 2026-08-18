package scm_api

import "slices"

type ScmApiConfig struct {
	RepoOwnerString string
	RepoNameString  string
	GitHubToken     string
}

type ScmApi interface {
	Tagging(commitSha string, tagName string, tagMessage string, apiConfig *ScmApiConfig) error
}

func IsSupportedScm(scm string) bool {
	supportedScms := []string{
		"github",
		"local-dry",
	}
	return slices.Contains(supportedScms, scm)
}
