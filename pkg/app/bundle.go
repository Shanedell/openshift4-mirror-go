package app

import (
	"os"
	"path/filepath"

	"github.com/Shanedell/openshift4-mirror-go/pkg/download"
	"github.com/Shanedell/openshift4-mirror-go/pkg/utils"
)

func Bundle(bundleDataIn *utils.BundleDataType) error {
	if utils.Err != nil {
		return utils.Err
	}

	if err := utils.Initialize(bundleDataIn); err != nil {
		return err
	}

	if err := download.Clients(); err != nil {
		return err
	}

	if !utils.BundleData.SkipRelease {
		if err := download.Release(); err != nil {
			return err
		}
	}

	if !utils.BundleData.SkipCatalogs {
		if err := download.Catalogs(); err != nil {
			return err
		}
	}

	if !utils.BundleData.SkipRhcos {
		if err := download.Rhcos(); err != nil {
			return err
		}
	}

	// cleanup files only needed local
	localOSTarball, err := utils.GetLocalTarballName()
	if err != nil {
		return err
	}
	os.Remove(localOSTarball)
	os.Remove(filepath.Join(utils.BundleDirs.Bin, utils.OCLocalCmdPath))

	return nil
}
