package utils

import (
	"os"
	"path/filepath"
)

func getDataDir() (string, error) {
	execDir, err := getExecDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(execDir, ".ggm"), nil
}

func getExecDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(execPath), nil
}
