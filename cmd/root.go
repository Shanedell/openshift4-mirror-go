package cmd

import (
	"github.com/spf13/cobra"
)

var (
	bundleDir            string
	preRelease           bool
	targetRegistry       string
	containerRuntime     string
	containerRuntimePath string
)

var rootHelp = "openshift4_mirror - CLI for mirroring OpenShift 4 content."

func addRootCommands(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(
		&bundleDir,
		"bundle-dir",
		"",
		"directory to save downloaded content",
	)

	cmd.PersistentFlags().BoolVar(
		&preRelease,
		"pre-release",
		false,
		"pre-release version of OpenShift",
	)

	cmd.PersistentFlags().StringVarP(
		&targetRegistry,
		"target-registry",
		"r",
		"example.registry.com",
		"target registry to tag the image with",
	)

	cmd.PersistentFlags().StringVarP(
		&containerRuntime,
		"containerRuntime",
		"c",
		"",
		"container runtime. supported options [docker, podman]. if not specified, code looks for both and uses whichever is found first.",
	)
	cmd.PersistentFlags().StringVar(
		&containerRuntimePath,
		"containerRuntimePath",
		"",
		"full to container runtime. needed if executable not in /usr/bin or /usr/local/bin",
	)
}

func NewRootCommand() *cobra.Command {
	var rootCommand = &cobra.Command{
		Use:   "openshift4_mirror",
		Short: rootHelp,
		Long:  rootHelp,
		Run:   func(_ *cobra.Command, _ []string) {},
	}

	addRootCommands(rootCommand)

	rootCommand.AddCommand(
		NewBundleCommand(),
		NewBuildCommand(),
		NewPruneCommand(),
		NewShellCommand(),
	)

	return rootCommand
}

func main() {
	cmd := NewRootCommand()

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func Run() {
	main()
}
