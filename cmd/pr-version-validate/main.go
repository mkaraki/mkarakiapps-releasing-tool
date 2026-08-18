package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	version "github.com/mkaraki/mkarakiapps-releasing-tool"
	"github.com/mkaraki/mkarakiapps-releasing-tool/internal/app_version"
	"github.com/mkaraki/mkarakiapps-releasing-tool/internal/versioning_management"
)

func main() {
	fmt.Printf("pr-version-validate %s\n\n", version.APP_VERSION)

	fs := flag.NewFlagSet("pr-version-validate", flag.ExitOnError)
	targetBranch := fs.String("target-branch", "master", "The target branch")
	versionFile := fs.String("version-file", "", "The version file")
	versionFileFormat := fs.String("version-file-format", "", "The version file format (plain, php)")
	localPatchPrefix := fs.String("local-patch-prefix", "", "The local patch prefix")

	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	if *versionFile == "" || *versionFileFormat == "" {
		panic("version-file and version-file-format are required")
	}

	if !app_version.IsSupportedFileType(*versionFileFormat) {
		panic("unsupported version file format: " + *versionFileFormat)
	}

	if !strings.HasPrefix(*targetBranch, "origin/") {
		*targetBranch = "origin/" + *targetBranch
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

	targetBranchVersionFileLines, err := app_version.ReadLinesFromGitFile(*targetBranch, *versionFile)
	if err != nil || len(targetBranchVersionFileLines) == 0 {
		println("Target branch version file does not exist, is empty, or could not be read, skipping version comparison.")
		return
	}

	targetBranchVersion, err := app_version.ExtractVersionFromFileLines(*versionFileFormat, targetBranchVersionFileLines)
	if err != nil {
		println("Unable to extract version from target branch version file, skipping version comparison.")
		return
	}

	newVersionObject := versioning_management.ParseReleaseVersionString(newVersion, *localPatchPrefix)
	targetBranchVersionObject := versioning_management.ParseReleaseVersionString(targetBranchVersion, *localPatchPrefix)

	if !versioning_management.IsAVersionLargerThanB(newVersionObject, targetBranchVersionObject) {
		panic("new version is not larger than the target branch version")
	}

	println("Version validation passed.")
}
