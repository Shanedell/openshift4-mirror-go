package cmd

import (
	"errors"

	"github.com/shanedell/openshift4-mirror-go/pkg/app"
	"github.com/shanedell/openshift4-mirror-go/pkg/utils"
	"github.com/spf13/cobra"
)

var buildHelp = "build the container image"

func NewBuildCommand() *cobra.Command {
	buildCommand := &cobra.Command{
		Use:   "build",
		Short: buildHelp,
		Long:  buildHelp,
		RunE:  buildMain,
	}

	return buildCommand
}

func buildMain(_ *cobra.Command, _ []string) error {
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

	return app.BuildContainer(containerData)
}
