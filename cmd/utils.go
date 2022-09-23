package cmd

import (
	"os"
)

func fileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func folderExists(folderName string) bool {
	_, err := os.Stat(folderName)
	return !os.IsNotExist(err)
}

func isValidPort(p int) bool {
	if p < 0 || p > 65535 {
		return false
	}

	return true
}
