package app

import (
	"os"
	"path/filepath"

	"github.com/shanedell/openshift4-mirror-go/pkg/download"
	"github.com/shanedell/openshift4-mirror-go/pkg/utils"
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

	if localOSTarball != "openshift-client-linux.tar.gz" {
		os.Remove(localOSTarball)
		os.Remove(filepath.Join(utils.BundleDirs.Bin, utils.OCLocalCmdPath))
	} else {
		if err := os.Rename(
			filepath.Join(utils.BundleDirs.Bin, utils.OCLocalCmdPath),
			filepath.Join(utils.BundleDirs.Bin, "oc"),
		); err != nil {
			return err
		}

		ocFileData, err := os.ReadFile(filepath.Join(utils.BundleDirs.Bin, "oc"))
		if err != nil {
			return err
		}

		if err := os.WriteFile(
			filepath.Join(utils.BundleDirs.Bin, "kubectl"),
			ocFileData,
			os.ModePerm,
		); err != nil {
			return err
		}
	}

	return nil
}
