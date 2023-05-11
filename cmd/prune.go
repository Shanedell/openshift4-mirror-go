package cmd

import (
	"errors"

	"github.com/Shanedell/openshift4-mirror-go/pkg/app"
	"github.com/Shanedell/openshift4-mirror-go/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	operators    []string
	imageToPrune string
	pruneType    string
	opmVersion   string
	targetImage  string
)

var pruneHelp = "prune the Red Hat Operator index image"

func NewPruneCommand() *cobra.Command {
	pruneCommand := &cobra.Command{
		Use:   "prune",
		Short: pruneHelp,
		Long:  pruneHelp,
		RunE:  pruneMain,
	}

	pruneCommand.PersistentFlags().StringVar(
		&pruneType,
		"type",
		"sqlite",
		"index image prunnings type. supported options: [sqlite, file]",
	)

	err := pruneCommand.MarkPersistentFlagRequired("type")
	if err != nil {
		panic(err)
	}

	pruneCommand.PersistentFlags().StringSliceVarP(
		&operators,
		"operators",
		"o",
		nil,
		"the operator(s) desired. Rest are pruned out",
	)

	pruneCommand.PersistentFlags().StringVar(
		&imageToPrune,
		"prune-image",
		"registry.redhat.io/redhat/redhat-operator-index:v4.10",
		"image to prune",
	)

	pruneCommand.PersistentFlags().StringVar(
		&opmVersion,
		"opm-version",
		"latest-4.9",
		"version of opm to download/use",
	)

	pruneCommand.PersistentFlags().StringVarP(
		&targetImage,
		"target-image",
		"t",
		"example.com/redhat-operators-index:latest",
		"complete image name to tag final image as.",
	)

	return pruneCommand
}

func pruneMain(_ *cobra.Command, _ []string) error {
	if containerRuntime == "" {
		containerRuntime = app.GetContainerRuntime()
	}

	if pruneType != "file" && pruneType != "sqlite" {
		panic(errors.New("unsupported prune type. Supported options: [file, sqlite]"))
	}

	bundleData := &utils.BundleDataType{
		OpenshiftVersion: opmVersion,
		PreRelease:       preRelease,
		BundleDir:        bundleDir,
	}

	pruneData := &utils.PruneDataType{
		Operators:    operators,
		ImageToPrune: imageToPrune,
		PruneType:    pruneType,
		OpmVersion:   opmVersion,
		TargetImage:  targetImage,
	}

	return app.PruneIndexImage(
		bundleData,
		pruneData,
		containerRuntime,
		containerRuntimePath,
	)
}
