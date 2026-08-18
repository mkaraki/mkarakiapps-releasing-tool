package versioning_management

import (
	"regexp"
	"strconv"
)

type SemanticVersion struct {
	Major int
	Minor int
	Patch int
}

type LocalPatch struct {
	PatchPrefix string
	LocalPatch  int
}

type ReleaseVersion struct {
	SemanticVersion
	LocalPatch
	DevelopmentSuffix string
}

func IsCleanReleaseVersionString(versionString string, localPatchPrefix string) bool {
	// Use regex to match a SemVer string
	re, err := regexp.Compile("^v?(\\d+)\\.(\\d+)\\.(\\d+)$")
	if err != nil {
		panic(err)
	}
	clean := re.MatchString(versionString)

	if clean {
		return true
	}

	if localPatchPrefix == "" {
		// If not clean and no local patch prefix. That is not clean SemVer.
		return false
	}

	re, err = regexp.Compile("^v?(\\d+)\\.(\\d+)\\.(\\d+)-" + regexp.QuoteMeta(localPatchPrefix) + "(\\d+)$")
	if err != nil {
		panic(err)
	}
	clean = re.MatchString(versionString)
	return clean
}

func ParseReleaseVersionString(versionString string, localPatchPrefix string) ReleaseVersion {
	var re *regexp.Regexp
	var err error

	if localPatchPrefix == "" {
		re, err = regexp.Compile("^v?(\\d+)\\.(\\d+)\\.(\\d+)")
	} else {
		re, err = regexp.Compile("^v?(\\d+)\\.(\\d+)\\.(\\d+)(?:-" + localPatchPrefix + "(\\d+))?")
	}
	if err != nil {
		panic(err)
	}

	matches := re.FindStringSubmatch(versionString)
	if matches == nil {
		panic("Invalid version string")
	}

	major := atoi(matches[1])
	minor := atoi(matches[2])
	patch := atoi(matches[3])

	var localPatch int = -1
	if localPatchPrefix != "" && len(matches) > 4 && matches[4] != "" {
		localPatch = atoi(matches[4])
	}

	return ReleaseVersion{
		SemanticVersion: SemanticVersion{
			Major: major,
			Minor: minor,
			Patch: patch,
		},
		LocalPatch: LocalPatch{
			PatchPrefix: localPatchPrefix,
			LocalPatch:  localPatch,
		},
	}
}

func IsAVersionLargerThanB(a, b ReleaseVersion) bool {
	if a.SemanticVersion.Major != b.SemanticVersion.Major {
		return a.SemanticVersion.Major > b.SemanticVersion.Major
	}
	if a.SemanticVersion.Minor != b.SemanticVersion.Minor {
		return a.SemanticVersion.Minor > b.SemanticVersion.Minor
	}
	if a.SemanticVersion.Patch != b.SemanticVersion.Patch {
		return a.SemanticVersion.Patch > b.SemanticVersion.Patch
	}
	if a.LocalPatch.LocalPatch != b.LocalPatch.LocalPatch {
		return a.LocalPatch.LocalPatch > b.LocalPatch.LocalPatch
	}
	return false
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}
