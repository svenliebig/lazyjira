package git

import (
	"errors"
	"os/exec"
	"strings"
)

var ErrNotGitRepo = errors.New("not a git repository")

// CommitsForIssue returns commit messages that mention the issue key.
func CommitsForIssue(issueKey string) ([]string, error) {
	cmd := exec.Command("git", "log", "--oneline", "--all")
	out, err := cmd.Output()
	if err != nil {
		return nil, ErrNotGitRepo
	}
	var commits []string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, issueKey) {
			commits = append(commits, line)
		}
	}
	return commits, nil
}
