package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	version "github.com/mkaraki/mkarakiapps-releasing-tool"
	"github.com/mkaraki/mkarakiapps-releasing-tool/internal/app_version"
	"github.com/mkaraki/mkarakiapps-releasing-tool/internal/git_utils"
	"github.com/mkaraki/mkarakiapps-releasing-tool/internal/scm_api"
	"github.com/mkaraki/mkarakiapps-releasing-tool/internal/versioning_management"
)

func main() {
	fmt.Printf("auto-tag %s\n\n", version.APP_VERSION)
	fs := flag.NewFlagSet("auto-tag", flag.ExitOnError)

	versionFile := fs.String("version-file", "", "The version file")
	versionFileFormat := fs.String("version-file-format", "", "The version file format")
	localPatchPrefix := fs.String("local-patch-prefix", "", "The local patch prefix")

	scm := fs.String("scm", "", "The source control management system")
	githubToken := fs.String("github-token", "", "The GitHub token")
	owner := fs.String("owner", "", "The repository owner name")
	repo := fs.String("repo", "", "The repository name")

	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	if *versionFile == "" || *versionFileFormat == "" {
		panic("version-file and version-file-format are required")
	}

	if !app_version.IsSupportedFileType(*versionFileFormat) {
		panic("unsupported version file format: " + *versionFileFormat)
	}

	if !scm_api.IsSupportedScm(*scm) {
		panic("unsupported scm: " + *scm)
	}

	var scmApiConfig *scm_api.ScmApiConfig

	if *githubToken == "" {
		// Fallback GitHub CI token
		*githubToken = os.Getenv("GITHUB_TOKEN")
	}

	if *scm == "github" && (*githubToken == "" || *owner == "" || *repo == "") {
		panic("github-token, owner and repo are required for github scm")
	} else {
		scmApiConfig = &scm_api.ScmApiConfig{
			GitHubToken:     *githubToken,
			RepoOwnerString: *owner,
			RepoNameString:  *repo,
		}
	}

	_, err := os.Stat(*versionFile)
	if os.IsNotExist(err) {
		panic("version file does not exist")
	} else if err != nil {
		panic(err)
	}

	newVersionFileLines, err := app_version.ReadLinesFromFile(*versionFile)
	if err != nil {
		panic(err)
	}

	newVersion, err := app_version.ExtractVersionFromFileLines(*versionFileFormat, newVersionFileLines)
	if err != nil {
		panic(err)
	}

	isNewVersionClean := versioning_management.IsCleanReleaseVersionString(newVersion, *localPatchPrefix)
	if !isNewVersionClean {
		panic("new version is not a clean release version string")
	}

	if !strings.HasPrefix(newVersion, "v") {
		newVersion = "v" + newVersion
	}

	commitSha, err := git_utils.GetCurrentCommitSHA()
	if err != nil {
		panic(err)
	}

	if *scm == "github" {
		githubApi := &scm_api.GitHubApi{}
		err = githubApi.Tagging(commitSha, newVersion, "Release "+newVersion, scmApiConfig)
		if err != nil {
			panic(err)
		}
	} else if *scm == "local-dry" {
		fmt.Printf("Dry run: would have tagged commit %s with tag %s\n", commitSha, newVersion)
	}

	fmt.Printf("Done: tagged commit %s with tag %s\n", commitSha, newVersion)
}
