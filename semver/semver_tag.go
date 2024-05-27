package semver

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/KaribuLab/kli/git"
	"github.com/spf13/cobra"
)

var majorRegex = regexp.MustCompile(`(?i)(!|BREAKING CHANGE)`)
var minorRegex = regexp.MustCompile(`(?i)feat`)
var patchRegex = regexp.MustCompile(`(?i)fix`)

func generateTag(pattern string, major, minor, patch int) string {
	tag := strings.ReplaceAll(pattern, "{major}", fmt.Sprint(major))
	tag = strings.ReplaceAll(tag, "{minor}", fmt.Sprint(minor))
	tag = strings.ReplaceAll(tag, "{patch}", fmt.Sprint(patch))
	return tag
}

func tagExists(tag string, tags []git.GitTag) bool {
	for _, t := range tags {
		if t.Tag == tag {
			return true
		}
	}
	return false
}

func createTagIfNeeded(pattern string, major int, minor int, patch int, createTags bool, commitHash string, tags []git.GitTag, verbose bool, gitCmd git.Cmd, cmd *cobra.Command) string {
	if createTags {
		tag := generateTag(pattern, major, minor, patch)
		if !tagExists(tag, tags) {
			cmd.Println(tag)
			gitCmd.Tag(verbose, tag, commitHash)
			gitCmd.PushTags(verbose, tag)
			return tag
		}
	}
	return ""
}

func NewSemverCommand(gitCmd git.Cmd) *cobra.Command {
	semverCmd := &cobra.Command{
		Use:   "semver",
		Short: "semver is a semver tool",
		Long:  "semver is a semver tool that does things",
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := cmd.Flags().Lookup("pattern").Value.String()
			verbose := cmd.Flags().Lookup("verbose").Value.String() == "true"
			createTags := cmd.Flags().Lookup("tags").Value.String() == "true"
			var tag string
			var major int
			var minor int
			var patch int
			logs, err := gitCmd.GetLogs(verbose)
			if err != nil {
				cmd.PrintErrln(err)
				return nil
			}
			var tags []git.GitTag
			if createTags {
				branch, err := gitCmd.CurrentBranch(verbose)
				if err != nil {
					cmd.PrintErrln(err)
					return nil
				}
				if branch != "main" {
					cmd.PrintErrln("branch is not main")
					return nil
				}
				tags, err = gitCmd.GetTags(verbose)
				if err != nil {
					cmd.PrintErrln(err)
					return nil
				}
			}

			for _, log := range logs {
				switch {
				case majorRegex.MatchString(log.Message):
					major++
					minor = 0
					patch = 0
					tag = createTagIfNeeded(pattern, major, minor, patch, createTags, log.Commit, tags, verbose, gitCmd, cmd)

				case minorRegex.MatchString(log.Message):
					minor++
					patch = 0
					tag = createTagIfNeeded(pattern, major, minor, patch, createTags, log.Commit, tags, verbose, gitCmd, cmd)
				case patchRegex.MatchString(log.Message):
					patch++
					tag = createTagIfNeeded(pattern, major, minor, patch, createTags, log.Commit, tags, verbose, gitCmd, cmd)
				}
			}
			if tag != "" {
				return nil
			}
			tag = generateTag(pattern, major, minor, patch)
			cmd.Println(tag)
			return nil
		},
	}
	semverCmd.Flags().StringP("pattern", "p", "v{major}.{minor}.{patch}", "Pattern to use for the tag")
	semverCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	semverCmd.Flags().BoolP("tags", "t", false, "Create all tags if not present")
	return semverCmd
}
