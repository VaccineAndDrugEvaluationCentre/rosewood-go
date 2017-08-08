package carpenter

import (
	"os"
	"path/filepath"
)

//todo merge with other utils

func GetWorkingDir() string {
	dir, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(dir)
}
