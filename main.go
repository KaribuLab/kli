package main

import (
	"fmt"
	"os"

	"github.com/KaribuLab/kli/git"
	"github.com/KaribuLab/kli/semver"
	"github.com/spf13/cobra"
)

func main() {
	gitCmd := git.NewGitCmd()
	rootCommand := &cobra.Command{
		Use:   "kli",
		Short: "kli util CLI tool",
		Long:  "kli util CLI tool for cool developers",
	}
	rootCommand.AddCommand(semver.NewSemverCommand(gitCmd))
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
