package download

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Shanedell/openshift4-mirror-go/pkg/utils"
)

func Clients() error {
	utils.Logger.Infoln("Starting client download")

	localOSTarball, err := utils.GetLocalTarballName()
	if err != nil {
		return err
	}

	for k, v := range map[string][]*utils.SaveFileToFrom{
		"openshift-install-linux.tar.gz": {
			{
				From: "openshift-install", To: "openshift-install",
			},
		},
		"openshift-client-linux.tar.gz": {
			{
				From: "oc", To: "oc",
			},
			{
				From: "kubectl", To: "kubectl",
			},
		},
		localOSTarball: {
			{
				From: "oc", To: utils.OCLocalCmdPath,
			},
		},
		"sha256sum.txt": nil,
	} {
		if err := downloadClient(k, v); err != nil {
			return nil
		}
	}

	// oc command is executable
	err = utils.MakeFileExecutableFile(fmt.Sprintf("%s/%s", utils.BundleDirs.Bin, utils.OCLocalCmdPath))
	if err != nil {
		return err
	}

	utils.Logger.Infoln("Finished client download")
	return nil
}

func downloadClient(filename string, filesToExtract []*utils.SaveFileToFrom) error {
	downloadURL := strings.Join(
		[]string{utils.ClientsBaseURL, utils.BundleData.OpenshiftVersion, filename},
		"/",
	)
	outputPath := filepath.Join(utils.BundleDirs.Clients, filename)

	if err := utils.DownloadFile(downloadURL, outputPath); err != nil {
		return err
	}

	if filesToExtract != nil {
		// Extract tarball contents
		filesAllExist := true

		for _, fileToCheck := range filesToExtract {
			if !utils.CheckExits(filepath.Join(utils.BundleDirs.Bin, fileToCheck.To)) {
				filesAllExist = false
				break
			}
		}

		if filesAllExist && utils.BundleData.SkipExisting {
			return nil
		}

		if err := utils.ExtractTar(utils.BundleDirs.Bin, outputPath, filesToExtract); err != nil {
			return err
		}
	}

	return nil
}
