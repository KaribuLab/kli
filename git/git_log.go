package git

import "fmt"

type GitLog struct {
	Commit  string
	Author  string
	Message string
}

func (g *GitLog) String() string {
	return fmt.Sprintf("%s: %s - %s", g.Commit, g.Author, g.Message)
}
