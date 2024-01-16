package cmd

import (
	"errors"

	"github.com/shanedell/openshift4-mirror-go/pkg/app"
	"github.com/shanedell/openshift4-mirror-go/pkg/utils"
	"github.com/spf13/cobra"
)

var shellHelp = "open a shell in the container environment"

func NewShellCommand() *cobra.Command {
	shellCommand := &cobra.Command{
		Use:   "shell",
		Short: shellHelp,
		Long:  shellHelp,
		RunE:  shellMain,
	}

	return shellCommand
}

func shellMain(_ *cobra.Command, _ []string) error {
	if containerRuntime == "" {
		containerRuntime = app.GetContainerRuntime()
	}

	if openshiftVersion == "" {
		return errors.New("please provide OpenShift Version using -v/--openshift-version")
	}

	containerData := &utils.ContainerDataType{
		OpenshiftVersion: openshiftVersion,
		Runtime:          containerRuntime,
		Image:            "localhost/openshift4-mirror:latest",
	}

	return app.Shell(containerData)
}
