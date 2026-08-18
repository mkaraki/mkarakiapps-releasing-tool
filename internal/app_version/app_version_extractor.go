package app_version

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
)

func ReadLinesFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	buf := make([]byte, 1024)
	var line []byte
	for {
		n, err := file.Read(buf)
		if n > 0 {
			line = append(line, buf[:n]...)
			for {
				index := indexOfNewline(line)
				if index == -1 {
					break
				}
				lines = append(lines, string(line[:index]))
				line = line[index+1:]
			}
		}
		if err != nil {
			if err.Error() == "EOF" {
				if len(line) > 0 {
					lines = append(lines, string(line))
				}
				break
			}
			return nil, err
		}
	}

	return lines, nil
}

func indexOfNewline(line []byte) int {
	for i, b := range line {
		if b == '\n' {
			return i
		}
	}
	return -1
}

func ReadLinesFromGitFile(branch string, filePath string) ([]string, error) {
	cmd := exec.Command("git", "show", branch+":"+filePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := []string{}
	currentLine := []byte{}
	for _, b := range output {
		if b == '\n' {
			lines = append(lines, string(currentLine))
			currentLine = []byte{}
		} else {
			currentLine = append(currentLine, b)
		}
	}
	if len(currentLine) > 0 {
		lines = append(lines, string(currentLine))
	}
	return lines, nil
}

func trimAnyQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '\'' || s[0] == '"' {
			s = s[1 : len(s)-1]
		}
		if s[len(s)-1] == '\'' || s[len(s)-1] == '"' {
			s = s[1 : len(s)-1]
		}
	}
	return s
}

func ExtractVersionFromPhpFileLines(fileLines []string) (string, error) {
	// To share this regex to other similar languages (e.g. GoLang), make last `?` to be optional.
	constRe, err := regexp.Compile(`^const APP_VERSION\s*=\s*('([^']+)'|"([^"]+)")\s*;?`)
	if err != nil {
		return "", err
	}

	for _, line := range fileLines {
		matches := constRe.FindStringSubmatch(line)
		if matches != nil {
			ver := trimAnyQuotes(matches[1])
			return ver, nil
		}
	}

	return "", nil
}

func IsSupportedFileType(fileType string) bool {
	supportedFileTypes := []string{
		"php",
	}
	return slices.Contains(supportedFileTypes, fileType)
}

func ExtractVersionFromFileLines(fileType string, fileLines []string) (string, error) {
	switch fileType {
	case "php":
		return ExtractVersionFromPhpFileLines(fileLines)
	default:
		return "", fmt.Errorf("unsupported file type: %s", fileType)
	}
}
