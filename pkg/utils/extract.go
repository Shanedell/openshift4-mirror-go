package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func createFiles(outputFolder string, tarballFile string, filesToExtract []*SaveFileToFrom, contains bool, baseFilename string, index int, tarReader *tar.Reader, headerName string) error {
	if contains {
		Logger.Infof("Extracting %s from %s\n", baseFilename, tarballFile)

		outputFile := filepath.Join(outputFolder, filesToExtract[index].To)
		filesToCreate := []string{outputFile}
		if baseFilename == "oc" {
			// Create kubectl from oc if that is a fileToExtract
			containsKubectl, _ := SaveFileToFromSliceContains(filesToExtract, "kubectl")
			if containsKubectl {
				filesToCreate = append(
					filesToCreate,
					// kubectl is a link to oc so normally it doesn't extracted
					strings.ReplaceAll(headerName, "oc", "kubectl"),
				)
			}
		}

		for _, fileToCreate := range filesToCreate {
			outFile, err := os.Create(fileToCreate)
			if err != nil {
				return fmt.Errorf("failed to create file: %s", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to copy contents: %s", err)
			}
			outFile.Close()
		}
	}

	return nil
}

// Sanitize archive file pathing from "G305: Zip Slip vulnerability"
func SanitizeArchivePath(d, t string) (v string, err error) {
	v = filepath.Join(d, t)
	if strings.HasPrefix(v, filepath.Clean(d)) {
		return v, nil
	}

	return "", fmt.Errorf("content filepath is tainted: %s", t)
}

func ExtractTar(outputFolder string, tarballFile string, filesToExtract []*SaveFileToFrom) error {
	gzipStream, err := os.Open(tarballFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip stream: %s", err)
	}
	defer gzipStream.Close()

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %s", err)
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to get header: %s", err)
		}

		headerName, err := SanitizeArchivePath(outputFolder, header.Name)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(headerName, 0755); err != nil {
				return fmt.Errorf("failed to mkdir: %s", err)
			}
		case tar.TypeReg:
			baseFilename := path.Base(headerName)
			contains, index := SaveFileToFromSliceContains(filesToExtract, baseFilename)
			if err := createFiles(
				outputFolder, tarballFile, filesToExtract,
				contains, baseFilename, index, tarReader, headerName,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
