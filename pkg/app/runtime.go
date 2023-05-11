package app

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Shanedell/openshift4-mirror-go/pkg/utils"
)

func GetContainerRuntime() string {
	for _, runtime := range []string{"podman", "docker"} {
		cmd := utils.CreateCommand(runtime, "--version")

		if err := cmd.Run(); err == nil {
			return runtime
		}
	}

	panic("No container runtime found")
}

func BuildContainer(containerData *utils.ContainerDataType) error {
	utils.Logger.Infoln("Building the container image")

	arch := "amd64"
	if runtime.GOARCH == "arm64" {
		arch = runtime.GOARCH
	}

	cmd := utils.CreateCommand(
		containerData.Runtime,
		"build",
		"--tag", containerData.Image,
		"--build-arg", fmt.Sprintf("arch=%s", arch),
		"-f", "debug.Dockerfile",
		".",
	)
	utils.SetCommandOutput(cmd)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build image: %s", err)
	}

	utils.Logger.Infoln("Finished building the container images")
	return nil
}

func BuildContainerIfMissing(containerData *utils.ContainerDataType) error {
	var stdout bytes.Buffer

	cmd := utils.CreateCommand(
		containerData.Runtime,
		"images",
		containerData.Image,
		"--format", `"{{json .}}`,
	)
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to list image: %s", err)
	}

	if stdout.String() == "" {
		utils.Logger.Warningf("The container image does not exist: %s\n", containerData.Image)
		return BuildContainer(containerData)
	}

	return nil
}

func Shell(containerData *utils.ContainerDataType) error {
	if err := BuildContainerIfMissing(containerData); err != nil {
		return err
	}

	utils.Logger.Infoln("Starting shell in container")

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := []string{
		containerData.Runtime,
		"run",
		"--interactive",
		"--tty",
		"--rm",
		"--hostname", "openshift4-mirror",
		"--security-opt", "label=disable",
		"--volume", fmt.Sprintf("%s:/app", cwd),
	}

	for _, envVar := range os.Environ() {
		envVarSplit := strings.Split(envVar, "=")
		key, value := envVarSplit[0], envVarSplit[1]

		if strings.HasPrefix(key, "OPENSHIFT_MIRROR_") {
			cmd = append(cmd, "--env", fmt.Sprintf("%s=%s", key, value))
		}
	}

	cmd = append(cmd, containerData.Image)

	cmdToRun := utils.CreateCommand(cmd[0], cmd[1:]...)
	utils.SetCommandOutput(cmdToRun)
	err = cmdToRun.Run()

	return err
}
