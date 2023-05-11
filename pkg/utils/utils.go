package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func GetVersionMinor(version string) string {
	versionParts := strings.Split(version, ".")
	return fmt.Sprintf("%s.%s", versionParts[0], versionParts[1])
}

func SetCommandOutput(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
}

func SaveFileToFromSliceContains(slice []*SaveFileToFrom, stringToCheckFor string) (bool, int) {
	for i, s := range slice {
		if s.From == stringToCheckFor {
			return true, i
		}
	}
	return false, -1
}

func DownloadFile(downloadURL string, outputPath string) error {
	if BundleData.SkipExisting && CheckExits(outputPath) {
		Logger.Infof(
			"Found existing file %s, skipping download of %s", outputPath, downloadURL,
		)
		return nil
	}

	Logger.Infof("Downloading %s\n", downloadURL)

	resp, err := http.Get(downloadURL) // nolint:gosec
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	Logger.Infof("Finished downloading %s", downloadURL)

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func SavePullSecret() error {
	pullSecretPath := filepath.Join(BundleDir, "pull-secret.json")
	Logger.Infof("Saving pull secret to %s\n", pullSecretPath)

	if CheckExits(pullSecretPath) {
		os.Remove(pullSecretPath)
	}

	out, err := os.Create(pullSecretPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.WriteString(BundleData.PullSecret)
	if err != nil {
		return err
	}

	out.Close()

	return nil
}

func GetLocalTarballName() (string, error) {
	var localOSTarball string

	osName, osArch := runtime.GOOS, runtime.GOARCH
	if osName == "darwin" {
		osName = "mac"
	}

	if osName == "windows" && osArch == "arm64" {
		return "", errors.New("error arm windows not supported")
	}

	if osName != "windows" && osName != "linux" && osName != "mac" {
		return "", fmt.Errorf("unsupported OS: %s", osName)
	}

	localOSTarball = fmt.Sprintf("openshift-client-%s", osName)

	if osArch == "arm64" {
		localOSTarball = fmt.Sprintf("%s-%s", localOSTarball, osArch)
	}

	localOSTarball = fmt.Sprintf("%s.tar.gz", localOSTarball)

	return localOSTarball, nil
}

func MakeFileExecutableFile(filepath string) error {
	return CreateCommand("chmod", "+x", filepath).Run()
}
