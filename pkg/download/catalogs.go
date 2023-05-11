package download

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Shanedell/openshift4-mirror-go/pkg/utils"
)

func mirrorCatalogs(outputDir string, catalog string) error {
	cmd := utils.CreateCommand(
		filepath.Join(utils.BundleDirs.Bin, utils.OCLocalCmdPath),
		"adm",
		"catalog",
		"mirror",
		"--registry-config", filepath.Join(utils.BundleDir, "pull-secret.json"),
		"--manifests-only",
		"--to-manifests", outputDir,
		utils.CatalogIndexes[catalog],
		utils.BundleData.TargetRegistry,
	)
	utils.SetCommandOutput(cmd)

	return cmd.Run()
}

func mirrorImages(outputDir string, mappingLocalFilePath string) error {
	cmd := utils.CreateCommand(
		filepath.Join(utils.BundleDirs.Bin, utils.OCLocalCmdPath),
		"image",
		"mirror",
		"--registry-config", filepath.Join(utils.BundleDir, "pull-secret.json"),
		"--dir", outputDir,
		"--filter-by-os", "linux/amd64",
		"--continue-on-error=true",
		"--filename", mappingLocalFilePath,
	)
	utils.SetCommandOutput(cmd)

	return cmd.Run()
}

// updates the mapping data to remove registry names and append the index image name to the path
func updateMappingData(data []byte, dummyRegistryRegexps []*regexp.Regexp) [][]byte {
	indexImageName := ""
	for _, d := range strings.Split(string(data), "\n") {
		if strings.Contains(d, "operator-index") {
			indexImageName = strings.ReplaceAll(
				strings.Split(strings.Split(d, "=")[1], ":")[0],
				fmt.Sprintf("%s/", utils.BundleData.TargetRegistry),
				"",
			)
		} else if strings.Contains(d, "marketplace-index") {
			indexImageName = strings.ReplaceAll(
				strings.Split(strings.Split(d, "=")[1], ":")[0],
				fmt.Sprintf("%s/", utils.BundleData.TargetRegistry),
				"",
			)
		}
	}

	indexImagePathing := fmt.Sprintf("file://%s/", indexImageName)
	duplicateIndexImagePathRegexp := regexp.MustCompile(
		fmt.Sprintf("%s/%s", indexImageName, indexImageName),
	)

	dataUpdated := [][]byte{
		dummyRegistryRegexps[0].ReplaceAll(data, []byte(indexImagePathing)),
		dummyRegistryRegexps[1].ReplaceAll(data, []byte(indexImagePathing)),
	}

	if indexImageName != "" {
		dataUpdated = [][]byte{
			duplicateIndexImagePathRegexp.ReplaceAll(dataUpdated[0], []byte(indexImageName)),
			duplicateIndexImagePathRegexp.ReplaceAll(dataUpdated[1], []byte(indexImageName)),
		}
	}

	return dataUpdated
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

		utils.Logger.Infof("Mirroring catalog manifests for %s from %s", catalog, utils.CatalogIndexes[catalog])

		if err := mirrorCatalogs(outputDir, catalog); err != nil {
			return err
		}

		mappingFilePath := filepath.Join(outputDir, "mapping.txt")
		mappingLocalFiles := []string{
			filepath.Join(outputDir, "mapping.local.txt"),
			filepath.Join(outputDir, "mapping.files.txt"),
		}

		data, err := os.ReadFile(mappingFilePath)
		if err != nil {
			return err
		}

		dummyRegistryRegexps := []*regexp.Regexp{
			regexp.MustCompile(fmt.Sprintf("%s/", utils.BundleData.TargetRegistry)),
			regexp.MustCompile(fmt.Sprintf(".*%s/", utils.BundleData.TargetRegistry)),
		}

		dataUpdated := updateMappingData(data, dummyRegistryRegexps)

		for i, mappingLocalFile := range mappingLocalFiles {
			err = os.WriteFile(mappingLocalFile, dataUpdated[i], 0644) // nolint:gosec
			if err != nil {
				return err
			}
		}

		utils.Logger.Infof("Mirroring catalog images for %s\n", catalog)

		if err := mirrorImages(outputDir, mappingLocalFiles[0]); err != nil {
			return err
		}
	}

	utils.Logger.Infoln("Finished catalogs download")

	return nil
}
