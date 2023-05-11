package download

import (
	"fmt"
	"path/filepath"

	"github.com/Shanedell/openshift4-mirror-go/pkg/utils"
)

func Release() error {
	utils.Logger.Infoln("Starting release download")
	outputPath := filepath.Join(utils.BundleDirs.Release, "v2")

	if utils.BundleData.SkipExisting && utils.CheckExits(outputPath) {
		utils.Logger.Infoln("Found existing release content, skipping download")
	} else {
		cmd := utils.CreateCommand(
			filepath.Join(utils.BundleDirs.Bin, utils.OCLocalCmdPath),
			"adm",
			"release",
			"mirror",
			"--registry-config", filepath.Join(utils.BundleDir, "pull-secret.json"),
			"--to-dir", utils.BundleDirs.Release,
			utils.BundleData.OpenshiftVersion,
		)
		utils.SetCommandOutput(cmd)

		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			return err
		}
	}

	utils.Logger.Infoln("Finished release download")
	return nil
}
