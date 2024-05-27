package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Cmd interface {
	Run(verbose bool, args ...string) (string, error)
	GetLogs(verbose bool) ([]GitLog, error)
	GetTags(verbose bool) ([]GitTag, error)
	Tag(verbose bool, tag string, commit string) error
	CurrentBranch(verbose bool) (string, error)
	PushTags(verbose bool, tag string) error
}

// GitCmd is a struct that holds the path to the git executable
type GitCmd struct {
	Path string
}

// NewGitCmd creates a new GitCmd struct
func NewGitCmd() *GitCmd {
	return &GitCmd{Path: "git"}
}

// Run executes a git command
func (g *GitCmd) Run(verbose bool, args ...string) (string, error) {
	cmd := exec.Command(g.Path, args...)
	out, err := cmd.CombinedOutput()
	if verbose {
		fmt.Println(strings.Join(cmd.Args, " "))
		fmt.Println(string(out))
	}
	if err != nil {
		return "", fmt.Errorf("error running git command: %s", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// Log returns a list of GitLog structs
func (g *GitCmd) GetLogs(verbose bool) ([]GitLog, error) {
	out, err := g.Run(verbose, "log", "--reverse", "--pretty=format:%H|%an|%s")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(out, "\n")
	logs := make([]GitLog, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, "|")
		logs[i] = GitLog{
			Commit: parts[0],
		}
		if len(parts) > 1 {
			logs[i].Author = parts[1]
		}
		if len(parts) > 2 {
			logs[i].Message = parts[2]
		}
	}
	return logs, nil
}

// Tag returns a list of GitTag structs
func (g *GitCmd) GetTags(verbose bool) ([]GitTag, error) {
	out, err := g.Run(verbose, "tag", "-l", "--format='%(objectname)|%(refname:short)'")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(out, "\n")
	tags := make([]GitTag, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, "|")
		tags[i] = GitTag{
			Commit: parts[0],
		}
		if len(parts) > 1 {
			tags[i].Tag = parts[1]
		}
	}
	return tags, nil
}

// Tag creates a new tag
func (g *GitCmd) Tag(verbose bool, tag string, commit string) error {
	out, err := g.Run(true, "tag", tag, commit)
	if verbose {
		fmt.Println(out)
	}
	if err != nil {
		return fmt.Errorf("error creating tag: %s", err)
	}
	return nil
}

// CurrentBranch returns the current branch
func (g *GitCmd) CurrentBranch(verbose bool) (string, error) {
	out, err := g.Run(verbose, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// PushTags pushes tags to the remote
func (g *GitCmd) PushTags(verbose bool, tag string) error {
	_, err := g.Run(verbose, "push", "origin", tag)
	if err != nil {
		return fmt.Errorf("error pushing tags: %s", err)
	}
	return nil
}
