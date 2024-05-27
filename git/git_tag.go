package git

import "fmt"

type GitTag struct {
	Tag    string
	Commit string
}

func (g *GitTag) String() string {
	return fmt.Sprintf("%s - %s", g.Commit, g.Tag)
}
