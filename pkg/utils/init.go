package utils

import (
	"path/filepath"
)

func Initialize(bundleDataIn *BundleDataType) error {
	BundleData = bundleDataIn

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
		"redhat-operators":    BundleData.RedhatOperatorIndexImage,
		"redhat-marketplace":  BundleData.RedhatMarketplaceIndexImage,
		"certified-operators": BundleData.CertifiedOperatorIndexImage,
		"community-operators": BundleData.CommunityOperatorIndexImage,
	}

	if err := CreateBundlesDir(); err != nil {
		return err
	}

	return SavePullSecret()
}
