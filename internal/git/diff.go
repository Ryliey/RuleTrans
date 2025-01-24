package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type FileChange struct {
	Status string
	Path   string
}

func GetDiffFiles(before, after string) ([]FileChange, error) {
	cmd := exec.Command("git", "diff", "--name-status", before, after)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git command failed: %w\nOutput: %s", err, string(output))
	}

	return parseGitOutput(string(output)), nil
}

func parseGitOutput(output string) []FileChange {
	var changes []FileChange
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if parts := strings.Fields(line); len(parts) >= 2 {
			changes = append(changes, FileChange{
				Status: string(parts[0][0]),
				Path:   parts[1],
			})
		}
	}
	return changes
}
