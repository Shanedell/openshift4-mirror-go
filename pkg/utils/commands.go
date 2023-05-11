package utils

import "os/exec"

func CreateCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

func CreateBashCommand(command string) *exec.Cmd {
	return exec.Command("bash", "-c", command)
}
