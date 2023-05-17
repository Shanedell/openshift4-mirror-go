package download

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/Shanedell/openshift4-mirror-go/pkg/utils"
)

func mirrorCatalogs(outputDir string, catalog string) error {
	cmd := utils.CreateCommand(
		filepath.Join(utils.BundleDirs.Bin, utils.OCLocalCmdPath),
		"adm",
		"catalog",
		"mirror",
		"--registry-config", filepath.Join(utils.BundleDir, "pull-secret.json"),
		"--index-filter-by-os", "linux/amd64",
		"--continue-on-error=true",
		"--to-manifests", outputDir,
		utils.CatalogIndexes[catalog],
		"file://local",
	)
	utils.SetCommandOutput(cmd)
	cmd.Dir = outputDir

	return cmd.Run()
}

func Catalogs() error {
	utils.Logger.Infoln("Starting catalogs download")

	for _, catalog := range utils.BundleData.Catalogs {
		outputDir := filepath.Join(utils.BundleDirs.Catalogs, catalog)

		if utils.BundleData.SkipExisting && utils.CheckExits(outputDir) {
			continue
		}

		if !utils.CheckExits(outputDir) {
			err := os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				return err
			}
		}

		utils.Logger.Infof("Mirroring catalogs for %s from %s", catalog, utils.CatalogIndexes[catalog])

		if err := mirrorCatalogs(outputDir, catalog); err != nil {
			return err
		}

		mappingFilePath := filepath.Join(outputDir, "mapping.txt")
		mappingLocalFile := filepath.Join(outputDir, "mapping.files.txt")

		data, err := os.ReadFile(mappingFilePath)
		if err != nil {
			return err
		}

		dummyRegistryRegexp := regexp.MustCompile(".*file://")

		dataUpdated := dummyRegistryRegexp.ReplaceAll(data, []byte("file://"))

		err = os.WriteFile(mappingLocalFile, dataUpdated, 0644) // nolint:gosec
		if err != nil {
			return err
		}
	}

	utils.Logger.Infoln("Finished catalogs download")

	return nil
}
