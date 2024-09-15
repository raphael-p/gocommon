package gocommon

import (
	"fmt"
	"os"
	"path/filepath"
)

// Find the directory the program's executable is stored in.
func GetExecDirectory(executableName, fallbackPath string) string {
	if os.Args[0] == executableName {
		// if executed from compiled binary, use its directory
		ex, err := os.Executable()
		if err != nil {
			panic(fmt.Sprintf("failed to locate executable: %s", err))
		}
		return filepath.Dir(ex)
	} else {
		// otherwise, use fallback directory (relative path)
		return fallbackPath
	}
}
