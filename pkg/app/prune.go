package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Shanedell/openshift4-mirror-go/pkg/download"
	"github.com/Shanedell/openshift4-mirror-go/pkg/utils"
)

func cleanupOpmFiles(osName string) {
	os.Remove(filepath.Join(utils.BundleDirs.Clients, fmt.Sprintf("opm-%s.tar.gz", osName)))
	os.Remove(filepath.Join(utils.BundleDirs.Bin, "opm"))
}

func commonPruneCommandOpts(cmd *exec.Cmd, containerRuntimePath string, setEnv bool) {
	utils.SetCommandOutput(cmd)

	if setEnv {
		// update path to be able to find container runtime executable
		if containerRuntimePath != "" {
			// remove podman from path
			parts := strings.Split(containerRuntimePath, "/")
			parts = parts[:len(parts)-1]

			cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=$PATH:/usr/bin:/usr/local/bin:%s", strings.Join(parts, "/")))
		} else {
			cmd.Env = append(cmd.Env, "PATH=$PATH:/usr/bin:/usr/local/bin")
		}
	}
}

func runOpmCommand(cmd *exec.Cmd, osName string, containerRuntimePath string, setEnv bool) error {
	commonPruneCommandOpts(cmd, containerRuntimePath, setEnv)
	if err := cmd.Run(); err != nil {
		cleanupOpmFiles(osName)
		return err
	}
	return nil
}

func pruneSqlite(pruneData *utils.PruneDataType, containerRuntime string, containerRuntimePath string, osName string) error {
	cmd := utils.CreateCommand(
		filepath.Join(utils.BundleDirs.Bin, "opm"),
		"index",
		"prune",
		"-c", containerRuntime,
		"-f", pruneData.ImageToPrune,
		"-p", strings.Join(pruneData.Operators, ","),
		"-t", pruneData.TargetImage,
	)
	return runOpmCommand(cmd, osName, containerRuntimePath, true)
}

func PruneIndexImage(bundleDataIn *utils.BundleDataType, pruneData *utils.PruneDataType, containerRuntime string, containerRuntimePath string) error {
	if err := utils.Initialize(bundleDataIn); err != nil {
		return err
	}

	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "mac"
	}

	if err := download.OpmScript(osName, pruneData.OpmVersion); err != nil {
		return err
	}

	pruneFunction := pruneFile
	if pruneData.PruneType == "sqlite" {
		pruneFunction = pruneSqlite
	}

	err := pruneFunction(pruneData, containerRuntime, containerRuntimePath, osName)
	if err != nil {
		return err
	}

	cleanupOpmFiles(osName)
	return nil
}
