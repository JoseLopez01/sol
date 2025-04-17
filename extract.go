package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractFile(file io.ReadCloser, dest string) error {
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		target, err := getPath(dest, header.Name)
		if err != nil {
			return fmt.Errorf("failed to get target path: %w", err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, target); err != nil {
				return fmt.Errorf("failed to create symlink: %w", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory for file: %w", err)
			}

			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				return fmt.Errorf("failed to extract file: %w", err)
			}

			if err := outFile.Close(); err != nil {
				return fmt.Errorf("failed to close file: %w", err)
			}

			if err := os.Chmod(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to set file permissions: %w", err)
			}
		default:
			return fmt.Errorf("unsupported tar header type: %c", header.Typeflag)
		}
	}

	return nil
}

func getPath(dest, name string) (string, error) {
	parts := strings.Split(name, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid tar header name: %s", name)
	}

	target := filepath.Join(dest, filepath.Join(parts[1:]...))
	return target, nil
}
