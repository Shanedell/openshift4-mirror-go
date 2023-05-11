package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Shanedell/openshift4-mirror-go/pkg/utils"
)

func createIndexFile(imageToPrune string, osName string, jqArgs string, prunedCatalogIndexFile string) error {
	opmCmd := fmt.Sprintf("%s render %s | jq 'select ( %s )'", filepath.Join(utils.BundleDirs.Bin, "opm"), imageToPrune, jqArgs)

	cmd := utils.CreateBashCommand(opmCmd)

	f, err := os.Create(prunedCatalogIndexFile)
	if err != nil {
		return err
	}

	cmd.Stdout = f
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		cleanupOpmFiles(osName)
		return err
	}

	return nil
}

func createDockerFile(containerRuntimePath string, osName string, prunedCatalogFolder string) error {
	cmd := utils.CreateCommand(
		filepath.Join(utils.BundleDirs.Bin, "opm"),
		"alpha",
		"generate",
		"dockerfile",
		prunedCatalogFolder,
	)
	return runOpmCommand(cmd, osName, containerRuntimePath, true)
}

func buildImage(imageName string, containerRuntime string, containerRuntimePath string, osName string) error {
	cmd := utils.CreateCommand(
		containerRuntime,
		"build",
		"-t", imageName,
		"-f", "configs.Dockerfile",
		".",
	)
	cmd.Dir = "pruned-catalog"

	return runOpmCommand(cmd, osName, containerRuntimePath, true)
}

func pruneFile(pruneData *utils.PruneDataType, containerRuntime string, containerRuntimePath string, osName string) error {
	prunedCatalogFolder := "pruned-catalog/configs"
	prunedCatalogIndexFile := fmt.Sprintf("%s/index.json", prunedCatalogFolder)

	jqArgs := ""
	for i, operator := range pruneData.Operators {
		if len(pruneData.Operators) > 1 && i != len(pruneData.Operators)-1 {
			jqArgs += fmt.Sprintf(`.name == "%s" or .package == "%s" or `, operator, operator)
		} else {
			jqArgs += fmt.Sprintf(`.name == "%s" or .package == "%s"`, operator, operator)
		}
	}

	if _, err := os.Stat(prunedCatalogFolder); os.IsExist(err) {
		if err := os.RemoveAll(prunedCatalogFolder); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(prunedCatalogFolder, os.ModePerm); err != nil {
		return err
	}

	utils.Logger.Info("Creating index file...")
	err := createIndexFile(pruneData.ImageToPrune, osName, jqArgs, prunedCatalogIndexFile)
	if err != nil {
		return err
	}
	utils.Logger.Info("Finished creating index file...")

	utils.Logger.Info("Creating dockerfile...")
	err = createDockerFile(containerRuntimePath, osName, prunedCatalogFolder)
	if err != nil {
		return err
	}
	utils.Logger.Info("Finished creating dockerfile...")

	utils.Logger.Infof("Building image %s...\n", pruneData.TargetImage)
	err = buildImage(pruneData.TargetImage, containerRuntime, containerRuntimePath, osName)
	if err != nil {
		return err
	}
	utils.Logger.Info("Finished building image...")

	return nil
}
