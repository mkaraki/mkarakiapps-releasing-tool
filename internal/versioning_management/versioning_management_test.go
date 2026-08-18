package versioning_management

import "testing"

func TestIsCleanReleaseVersionString(t *testing.T) {
	tests := []struct {
		versionString    string
		localPatchPrefix string
		expected         bool
	}{
		{"v1.2.3", "", true},
		{"1.2.3", "", true},
		{"v1.2.3-local1", "local", true},
		{"v1.2.3-local1", "", false},
		{"v1.2.3-remote1", "local", false},
	}

	for _, test := range tests {
		t.Run(test.versionString, func(t *testing.T) {
			result := IsCleanReleaseVersionString(test.versionString, test.localPatchPrefix)
			if result != test.expected {
				t.Errorf("IsCleanReleaseVersionString(%q, %q) = %v; want %v",
					test.versionString, test.localPatchPrefix, result, test.expected)
			}
		})
	}
}

func TestParseReleaseVersionString(t *testing.T) {
	tests := []struct {
		versionString    string
		localPatchPrefix string
		expected         ReleaseVersion
	}{
		{"v1.2.3", "", ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"", -1}, ""}},
		{"1.2.3", "", ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"", -1}, ""}},
		{"v1.2.3-local1", "local", ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"local", 1}, ""}},
		{"v1.2.3-local1", "", ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"", -1}, ""}},
		{"v1.2.3-remote1", "local", ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"local", -1}, ""}},
	}

	for _, test := range tests {
		t.Run(test.versionString, func(t *testing.T) {
			result := ParseReleaseVersionString(test.versionString, test.localPatchPrefix)
			if result != test.expected {
				t.Errorf("ParseReleaseVersionString(%q, %q) = %v; want %v",
					test.versionString, test.localPatchPrefix, result, test.expected)
			}
		})
	}
}

func TestIsAVersionLargerThanB(t *testing.T) {
	tests := []struct {
		a        ReleaseVersion
		b        ReleaseVersion
		expected bool
	}{
		{
			ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"", -1}, ""},
			ReleaseVersion{SemanticVersion{1, 2, 2}, LocalPatch{"", -1}, ""},
			true,
		},
		{
			ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"", -1}, ""},
			ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"", -1}, ""},
			false,
		},
		{
			ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"", -1}, ""},
			ReleaseVersion{SemanticVersion{1, 2, 4}, LocalPatch{"", -1}, ""},
			false,
		},
		{
			ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"local", 1}, ""},
			ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"local", 0}, ""},
			true,
		},
		{
			ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"", -1}, ""},
			ReleaseVersion{SemanticVersion{2, 0, 0}, LocalPatch{"", -1}, ""},
			false,
		},
		{
			ReleaseVersion{SemanticVersion{2, 0, 0}, LocalPatch{"", -1}, ""},
			ReleaseVersion{SemanticVersion{1, 2, 3}, LocalPatch{"", -1}, ""},
			true,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := IsAVersionLargerThanB(test.a, test.b)
			if result != test.expected {
				t.Errorf("IsAVersionLargerThanB(%v, %v) = %v; want %v",
					test.a, test.b, result, test.expected)
			}
		})
	}
}
