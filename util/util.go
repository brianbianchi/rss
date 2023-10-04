package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetRootPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return fmt.Sprint(exPath, "/../")
}
