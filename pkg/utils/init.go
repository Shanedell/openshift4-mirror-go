package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

func Initialize(bundleDataIn *BundleDataType) error {
	BundleData = bundleDataIn

	versionMinor := BundleData.OpenshiftVersion
	if !strings.Contains(BundleData.OpenshiftVersion, "latest-") {
		versionMinor = GetVersionMinor(BundleData.OpenshiftVersion)
	}

	if BundleData.BundleDir != "" {
		BundleDir = filepath.Join(Cwd, BundleData.BundleDir, BundleData.OpenshiftVersion)
	} else {
		BundleDir = filepath.Join(Cwd, "bundle", BundleData.OpenshiftVersion)
	}
	BundleDirs = &BundleDirsType{
		Bin:      filepath.Join(BundleDir, "bin"),
		Release:  filepath.Join(BundleDir, "release"),
		Rhcos:    filepath.Join(BundleDir, "rhcos"),
		Catalogs: filepath.Join(BundleDir, "catalogs"),
		Clients:  filepath.Join(BundleDir, "clients"),
	}

	CatalogIndexes = map[string]string{
		"redhat-operators": BundleData.RedhatOperatorIndexImage,
		"certified-operators": fmt.Sprintf(
			"registry.redhat.io/redhat/certified-operator-index:v%s", versionMinor,
		),
		"redhat-marketplace": fmt.Sprintf(
			"registry.redhat.io/redhat/redhat-marketplace-index:v%s", versionMinor,
		),
		"community-operators": "registry.redhat.io/redhat/community-operator-index:latest",
	}

	if err := CreateBundlesDir(); err != nil {
		return err
	}

	return SavePullSecret()
}
