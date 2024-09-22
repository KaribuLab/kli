package project

import (
	"os/exec"
	"strings"
)

func runHook(workdir string, cmd string) error {
	tokens := strings.Split(cmd, " ")
	command := exec.Command(tokens[0], tokens[1:]...)
	command.Dir = workdir
	err := command.Run()
	if err != nil {
		return err
	}
	return nil
}
