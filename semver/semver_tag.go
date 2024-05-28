package semver

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/KaribuLab/kli/git"
	"github.com/spf13/cobra"
)

var majorRegex = regexp.MustCompile(`(?i)(!|BREAKING CHANGE)`)
var minorRegex = regexp.MustCompile(`(?i)feat`)
var patchRegex = regexp.MustCompile(`(?i)fix`)

func generateTag(verbose bool, pattern string, major, minor, patch int) string {
	tag := strings.ReplaceAll(pattern, "{major}", fmt.Sprint(major))
	tag = strings.ReplaceAll(tag, "{minor}", fmt.Sprint(minor))
	tag = strings.ReplaceAll(tag, "{patch}", fmt.Sprint(patch))
	if verbose {
		fmt.Println("Tag", tag)
	}
	return tag
}

func tagExists(verbose bool, tag string, tags []git.GitTag) bool {
	for _, t := range tags {
		if verbose {
			fmt.Printf("Checking tag %s - %s\n", t.Tag, tag)
		}
		if t.Tag == tag {
			return true
		}
	}
	return false
}

func createTagIfNeeded(pattern string, major int, minor int, patch int, createTags bool, commitHash string, tags []git.GitTag, verbose bool, gitCmd git.Cmd) string {
	if createTags {
		tag := generateTag(verbose, pattern, major, minor, patch)
		if !tagExists(verbose, tag, tags) {
			fmt.Println(tag)
			gitCmd.Tag(verbose, tag, commitHash)
			gitCmd.PushTags(verbose, tag)
			time.Sleep(10 * time.Second)
			return tag
		}
	}
	return ""
}

func removeTagIfNeeded(tag string, removeTags bool, tags []git.GitTag, verbose bool, gitCmd git.Cmd) {
	if removeTags {
		exists := tagExists(verbose, tag, tags)
		if verbose {
			fmt.Println("Tag exists", exists)
		}
		if exists {
			fmt.Println("Removing tag", tag)
			gitCmd.RemoveTag(verbose, tag)
		}
	}
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
			dryRun := cmd.Flags().Lookup("dryrun").Value.String() == "true"
			removeTags := cmd.Flags().Lookup("remove").Value.String() == "true"
			var tag string
			var major int
			var minor int
			var patch int
			logs, err := gitCmd.GetLogs(verbose)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "error getting logs: %s\n", err)
				return nil
			}
			var tags []git.GitTag
			if createTags || removeTags {
				branch, err := gitCmd.CurrentBranch(verbose)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "error getting branch: %s\n", err)
					return nil
				}
				if branch != "main" {
					fmt.Fprintf(cmd.ErrOrStderr(), "branch is not main: %s\n", branch)
					return nil
				}
				tags, err = gitCmd.GetTags(verbose)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "error getting tags: %s\n", err)
					return nil
				}
				if verbose {
					fmt.Println("tags", tags)
				}
			}

			for _, log := range logs {
				if verbose {
					fmt.Println(log)
				}
				switch {
				case majorRegex.MatchString(log.Message):
					major++
					minor = 0
					patch = 0
					tag = generateTag(verbose, pattern, major, minor, patch)
					if dryRun {
						fmt.Println(tag)
					} else if removeTags {
						removeTagIfNeeded(tag, removeTags, tags, verbose, gitCmd)
					} else {
						tag = createTagIfNeeded(pattern, major, minor, patch, createTags, log.Commit, tags, verbose, gitCmd)
					}

				case minorRegex.MatchString(log.Message):
					minor++
					patch = 0
					tag = generateTag(verbose, pattern, major, minor, patch)
					if dryRun {
						fmt.Println(tag)
					} else if removeTags {
						removeTagIfNeeded(tag, removeTags, tags, verbose, gitCmd)
					} else {
						tag = createTagIfNeeded(pattern, major, minor, patch, createTags, log.Commit, tags, verbose, gitCmd)
					}
				case patchRegex.MatchString(log.Message):
					patch++
					tag = generateTag(verbose, pattern, major, minor, patch)
					if dryRun {
						fmt.Println(tag)
					} else if removeTags {
						removeTagIfNeeded(tag, removeTags, tags, verbose, gitCmd)
					} else {
						tag = createTagIfNeeded(pattern, major, minor, patch, createTags, log.Commit, tags, verbose, gitCmd)
					}
				}
			}
			if tag != "" {
				return nil
			}
			tag = generateTag(verbose, pattern, major, minor, patch)
			fmt.Println(tag)
			return nil
		},
	}
	semverCmd.Flags().StringP("pattern", "p", "v{major}.{minor}.{patch}", "Pattern to use for the tag")
	semverCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	semverCmd.Flags().BoolP("tags", "t", false, "Create all tags if not present")
	semverCmd.Flags().BoolP("dryrun", "d", false, "Dry run mode")
	semverCmd.Flags().BoolP("remove", "r", false, "Remove tags")
	return semverCmd
}
