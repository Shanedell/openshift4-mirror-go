package utils

import (
	"errors"
	"os"
)

func CheckExits(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func CreateBundlesDir() error {
	dirs := []string{
		BundleDirs.Bin,
		BundleDirs.Catalogs,
		BundleDirs.Clients,
		BundleDirs.Release,
		BundleDirs.Rhcos,
	}

	if CheckExits(BundleDir) && !BundleData.SkipExisting {
		if err := os.RemoveAll(BundleDir); err != nil {
			return nil
		}
	}

	for _, dir := range dirs {
		if err := CreateCleanDir(dir); err != nil {
			return err
		}
	}

	return nil
}

func CreateCleanDir(path string) error {
	if !CheckExits(path) {
		return os.MkdirAll(path, os.ModePerm)
	}

	if !BundleData.SkipExisting {
		if err := os.Remove(path); err != nil {
			return err
		}
		return os.MkdirAll(path, os.ModePerm)
	}

	return nil
}
