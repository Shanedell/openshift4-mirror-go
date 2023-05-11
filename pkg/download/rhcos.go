package download

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Shanedell/openshift4-mirror-go/pkg/utils"
)

func rhcosFilename() (string, error) {
	filename := fmt.Sprintf("rhcos-%s.x86_64", utils.BundleData.Platform)

	switch utils.BundleData.Platform {
	case "aws":
		return fmt.Sprintf("%s.vmdk.gz", filename), nil
	case "azure":
		return fmt.Sprintf("%s.vhd.gz", filename), nil
	case "gcp":
		return fmt.Sprintf("%s.tar.gz", filename), nil
	case "metal":
		return fmt.Sprintf("%s.raw.gz", filename), nil
	case "openstack":
		return fmt.Sprintf("%s.qcow2.gz", filename), nil
	case "vmware":
		return fmt.Sprintf("%s.ova", filename), nil
	default:
		return "", fmt.Errorf(
			"invalid platform. Allowed platforms: [aws, azure, gcp, metal, openstack, vmware]",
		)
	}
}

func Rhcos() error {
	utils.Logger.Infoln("Starting RHCOS download")

	filename, err := rhcosFilename()
	if err != nil {
		return err
	}

	downloadURL := strings.Join(
		[]string{utils.RhcosBaseURL, utils.BundleData.OpenshiftVersion, "latest", filename},
		"/",
	)

	if utils.BundleData.PreRelease {
		versionSplit := strings.Split(utils.BundleData.OpenshiftVersion, ".")
		downloadVersion := strings.Join(versionSplit[0:2], ".")
		downloadURL = strings.Join(
			[]string{
				utils.RhcosPreBaseURL,
				fmt.Sprintf("latest-%s", downloadVersion), filename,
			},
			"/",
		)
	}

	outputPath := filepath.Join(utils.BundleDirs.Rhcos, filename)

	if err := utils.DownloadFile(downloadURL, outputPath); err != nil {
		return err
	}

	utils.Logger.Infoln("Finished RHCOS download")
	return nil
}
