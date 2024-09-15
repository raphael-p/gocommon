package gocommon

import (
	"fmt"
	"os"
	"path/filepath"
)

// Checks if we are running the program as its executable.
// If so, fetch the executable directory; otherwise, returns an
// empty string.
func GetExecDirectory(executableName string) string {
	if os.Args[0] == executableName {
		ex, err := os.Executable()
		if err != nil {
			panic(fmt.Sprintf("failed to locate executable: %s", err))
		}
		return filepath.Dir(ex)
	}

	return ""
}
