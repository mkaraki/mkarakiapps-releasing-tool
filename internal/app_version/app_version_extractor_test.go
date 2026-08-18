package app_version

import "testing"

func TestExtractVersionFromPhpFileLines(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{"<?php", "const APP_VERSION = 'v1.2.3'"}, "v1.2.3"},
		{[]string{"<?php", `const APP_VERSION = "1.2.3";`}, "1.2.3"},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result, err := ExtractVersionFromPhpFileLines(test.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}
