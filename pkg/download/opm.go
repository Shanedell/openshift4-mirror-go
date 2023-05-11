package download

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Shanedell/openshift4-mirror-go/pkg/utils"
)

func OpmScript(osName string, version string) error {
	utils.Logger.Infoln("Starting opm download")

	filename := fmt.Sprintf("opm-%s.tar.gz", osName)
	downloadURL := strings.Join(
		[]string{utils.ClientsBaseURL, version, filename},
		"/",
	)

	if utils.BundleData.PreRelease {
		versionSplit := strings.Split(version, ".")
		downloadVersion := strings.Join(versionSplit[0:2], ".")
		downloadURL = strings.Join(
			[]string{
				utils.RhcosPreBaseURL,
				fmt.Sprintf("latest-%s", downloadVersion), filename,
			},
			"/",
		)
	}

	outputPath := filepath.Join(utils.BundleDirs.Clients, filename)

	if err := utils.DownloadFile(downloadURL, outputPath); err != nil {
		return err
	}

	utils.Logger.Infoln("Finished opm download")

	startScriptName := "opm"
	if runtime.GOOS != "linux" {
		startScriptName = fmt.Sprintf("%s-amd64-opm", runtime.GOOS)
	}

	err := utils.ExtractTar(
		utils.BundleDirs.Bin,
		outputPath,
		[]*utils.SaveFileToFrom{
			{From: startScriptName, To: "opm"},
		},
	)

	if err != nil {
		return err
	}

	// make sure opm file is executable
	err = utils.MakeFileExecutableFile(fmt.Sprintf("%s/opm", utils.BundleDirs.Bin))
	if err != nil {
		return err
	}

	return nil
}
